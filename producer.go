package components

import (
	"context"
	"fmt"
	"strings"

	"github.com/asecurityteam/settings"
)

const (
	// ProducerTypeBenthos is used to select a Benthos producer.
	ProducerTypeBenthos = "BENTHOS"
	// ProducerTypePOST is used to select an HTTP POST client.
	ProducerTypePOST = "POST"
)

// ProducerConfig wraps all producer related configuration.
type ProducerConfig struct {
	Type    string `default:"The type of producer. The choices are BENTHOS and POST."`
	Benthos *ProducerBenthosConfig
	POST    *ProducerPOSTConfig
}

// Name of the configuration.
func (*ProducerConfig) Name() string {
	return "producer"
}

// ProducerComponent is the top level producer component.
type ProducerComponent struct {
	Benthos *ProducerBenthosComponent
	POST    *ProducerPOSTComponent
}

// NewProducerComponent populates a ProducerComponent with defaults.
func NewProducerComponent() *ProducerComponent {
	return &ProducerComponent{
		Benthos: NewProducerBenthosComponent(),
		POST:    NewProducerPOSTComponent(),
	}
}

// Settings generates a default configuration.
func (c *ProducerComponent) Settings() *ProducerConfig {
	return &ProducerConfig{
		Type:    "BENTHOS",
		Benthos: c.Benthos.Settings(),
		POST:    c.POST.Settings(),
	}
}

// New constructs a Producer from the given configuration.
func (c *ProducerComponent) New(ctx context.Context, conf *ProducerConfig) (Producer, error) {
	switch {
	case strings.EqualFold(conf.Type, ProducerTypeBenthos):
		return c.Benthos.New(ctx, conf.Benthos)
	case strings.EqualFold(conf.Type, ProducerTypePOST):
		return c.POST.New(ctx, conf.POST)
	default:
		return nil, fmt.Errorf("unknown producer type %s", conf.Type)
	}
}

// LoadProducer is a convenience method for binding the source to the component.
func LoadProducer(ctx context.Context, source settings.Source, c *ProducerComponent) (Producer, error) {
	dst := new(Producer)
	err := settings.NewComponent(ctx, source, c, dst)
	if err != nil {
		return nil, err
	}
	return *dst, nil
}

// NewProducer is the top-level entry point for creating a producer client.
func NewProducer(ctx context.Context, source settings.Source) (Producer, error) {
	return LoadProducer(ctx, source, NewProducerComponent())
}
