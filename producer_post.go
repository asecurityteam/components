package components

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type postProducer struct {
	Client   *http.Client
	Endpoint *url.URL
}

// Produce an event to the endpoint. Any 2xx series response is a success.
func (p *postProducer) Produce(ctx context.Context, event interface{}) (interface{}, error) {
	b, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	r, _ := http.NewRequest(http.MethodPost, p.Endpoint.String(), ioutil.NopCloser(bytes.NewReader(b)))
	res, err := p.Client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// Drain the body no matter what in order to allow for connection re-use
	// in HTTP/1.x systems.
	rb, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode > 200 {
		return nil, fmt.Errorf("failed to post. status(%d) reason(%s)", res.StatusCode, string(rb))
	}
	return event, nil
}

// ProducerPOSTConfig contains settings for the HTTP POST producer.
type ProducerPOSTConfig struct {
	Endpoint   string `description:"The URL to POST."`
	HTTPClient *HTTPConfig
}

// Name of the configuration section.
func (*ProducerPOSTConfig) Name() string {
	return "post"
}

// ProducerPOSTComponent is a component for creating an HTTP POST producer.
type ProducerPOSTComponent struct {
	HTTP *HTTPComponent
}

// NewProducerPOSTComponent populates a ProducerPOSTComponent with defaults.
func NewProducerPOSTComponent() *ProducerPOSTComponent {
	return &ProducerPOSTComponent{HTTP: NewHTTPComponent()}
}

// Settings returns the default configuration.
func (c *ProducerPOSTComponent) Settings() *ProducerPOSTConfig {
	return &ProducerPOSTConfig{
		HTTPClient: c.HTTP.Settings(),
	}
}

// New constructs a benthos configuration.
func (c *ProducerPOSTComponent) New(ctx context.Context, conf *ProducerPOSTConfig) (Producer, error) {
	if conf.Endpoint == "" {
		return nil, fmt.Errorf("missing POST producer endpoint")
	}
	u, err := url.Parse(conf.Endpoint)
	if err != nil {
		return nil, err
	}
	tr, err := c.HTTP.New(ctx, conf.HTTPClient)
	if err != nil {
		return nil, err
	}
	cl := &http.Client{Transport: tr}
	return &postProducer{Endpoint: u, Client: cl}, nil
}
