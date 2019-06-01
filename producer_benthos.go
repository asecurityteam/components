package components

import (
	"context"

	"github.com/Jeffail/benthos/lib/config"
	"github.com/Jeffail/benthos/lib/serverless"
	"github.com/Jeffail/benthos/lib/util/text"
	"gopkg.in/yaml.v2"
)

type benthosProducer struct {
	Handler *serverless.Handler
}

// Produce an event using the Benthos serverless handler. The handler is
// exactly a producer as we want it. We only adapt the function name here.
func (p *benthosProducer) Produce(ctx context.Context, event interface{}) (interface{}, error) {
	return p.Handler.Handle(ctx, event)
}

// ProducerBenthosConfig adapts the Benthos configuration system to this one.
type ProducerBenthosConfig struct {
	YAML string `description:"The YAML or JSON text of a Benthos configuration."`
}

// Name of the configuration section.
func (*ProducerBenthosConfig) Name() string {
	return "benthos"
}

// ProducerBenthosComponent is a component for creating a Benthos producer.
type ProducerBenthosComponent struct{}

// NewProducerBenthosComponent generates a ProducerBenthosComponent and
// populates it with defaults
func NewProducerBenthosComponent() *ProducerBenthosComponent {
	return &ProducerBenthosComponent{}
}

// Settings returns the default configuration.
func (*ProducerBenthosComponent) Settings() *ProducerBenthosConfig {
	return &ProducerBenthosConfig{}
}

// newBenthosConfig parses a Benthos YAML config file and replaces all
// environment variables. This is exactly what is done in the Benthos project
// when starting one of the main() functions.
func newBenthosConfig(b []byte) (config.Type, error) {
	conf := config.New()
	err := yaml.Unmarshal(text.ReplaceEnvVariables(b), &conf)
	return conf, err
}

// New constructs a benthos producer.
func (*ProducerBenthosComponent) New(ctx context.Context, conf *ProducerBenthosConfig) (Producer, error) {
	cfg, err := newBenthosConfig([]byte(conf.YAML))
	if err != nil {
		return nil, err
	}
	handler, err := serverless.NewHandler(cfg)
	if err != nil {
		return nil, err
	}
	return &benthosProducer{Handler: handler}, nil
}
