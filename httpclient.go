package components

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/asecurityteam/settings"
	transportd "github.com/asecurityteam/transportd/pkg"
	componentsd "github.com/asecurityteam/transportd/pkg/components"
)

const (
	// HTTPTypeDefault is used to select the default Go HTTP client.
	HTTPTypeDefault = "DEFAULT"
	// HTTPTypeSmart is used to select the transportd HTTP client.
	HTTPTypeSmart = "SMART"
)

// HTTPDefaultConfig contains all settings for the default Go HTTP client.
type HTTPDefaultConfig struct{}

// HTTPDefaultComponent is a component for creating the default Go HTTP client.
type HTTPDefaultComponent struct{}

// Settings returns the default configuration.
func (*HTTPDefaultComponent) Settings() *HTTPDefaultConfig {
	return &HTTPDefaultConfig{}
}

// New constructs a client from the given configuration
func (*HTTPDefaultComponent) New(ctx context.Context, conf *HTTPDefaultConfig) (http.RoundTripper, error) {
	return http.DefaultTransport, nil
}

// HTTPSmartConfig contains all settings for the transportd HTTP client.
type HTTPSmartConfig struct {
	OpenAPI string `description:"The full OpenAPI specification with transportd extensions."`
}

// Name of the configuration tree.
func (*HTTPSmartConfig) Name() string {
	return "smart"
}

// HTTPSmartComponent is a component for creating a transportd HTTP client.
type HTTPSmartComponent struct {
	Plugins []transportd.NewComponent
}

// Settings returns the default configuration.
func (*HTTPSmartComponent) Settings() *HTTPSmartConfig {
	return &HTTPSmartConfig{}
}

// New constructs a client from the given configuration.
func (c *HTTPSmartComponent) New(ctx context.Context, conf *HTTPSmartConfig) (http.RoundTripper, error) {
	return transportd.NewTransport(ctx, []byte(conf.OpenAPI), c.Plugins...)
}

// HTTPConfig wraps all HTTP related settings.
type HTTPConfig struct {
	Type    string `description:"The type of HTTP client. Choices are SMART and DEFAULT."`
	Default *HTTPDefaultConfig
	Smart   *HTTPSmartConfig
}

// Name of the config.
func (*HTTPConfig) Name() string {
	return "httpclient"
}

// HTTPComponent is the top level HTTP client component.
type HTTPComponent struct {
	Default *HTTPDefaultComponent
	Smart   *HTTPSmartComponent
}

// NewHTTPComponent populates an HTTPComponent with defaults.
func NewHTTPComponent() *HTTPComponent {
	return &HTTPComponent{
		Default: &HTTPDefaultComponent{},
		Smart: &HTTPSmartComponent{
			Plugins: componentsd.Defaults,
		},
	}
}

// Settings returns the default configuration.
func (c *HTTPComponent) Settings() *HTTPConfig {
	return &HTTPConfig{
		Type:    "DEFAULT",
		Default: c.Default.Settings(),
		Smart:   c.Smart.Settings(),
	}
}

// New constructs a client from the given configuration.
func (c *HTTPComponent) New(ctx context.Context, conf *HTTPConfig) (http.RoundTripper, error) {
	switch {
	case strings.EqualFold(conf.Type, HTTPTypeDefault):
		return c.Default.New(ctx, conf.Default)
	case strings.EqualFold(conf.Type, HTTPTypeSmart):
		return c.Smart.New(ctx, conf.Smart)
	default:
		return nil, fmt.Errorf("unknown HTTP client type %s", conf.Type)
	}
}

// LoadHTTP is a convenience method for binding the source to the component.
func LoadHTTP(ctx context.Context, source settings.Source, c *HTTPComponent) (http.RoundTripper, error) {
	dst := new(http.RoundTripper)
	err := settings.NewComponent(ctx, source, c, dst)
	if err != nil {
		return nil, err
	}
	return *dst, nil
}

// NewHTTP is the top-level entry point for creating a new HTTP client.
// The default set of plugins will be installed for the smart client. Use the
// LoadHTTP() method if a custom set of plugins are required.
func NewHTTP(ctx context.Context, source settings.Source) (http.RoundTripper, error) {
	return LoadHTTP(ctx, source, NewHTTPComponent())
}
