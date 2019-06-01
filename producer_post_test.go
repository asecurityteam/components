package components

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPostProducerCantMarshal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tr := NewMockRoundTripper(ctrl)
	u, _ := url.Parse("http://localhost")
	event := make(chan interface{})

	p := &postProducer{
		Client:   &http.Client{Transport: tr},
		Endpoint: u,
	}
	_, err := p.Produce(context.Background(), event)
	require.NotNil(t, err)
}

func TestPostProducerClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tr := NewMockRoundTripper(ctrl)
	u, _ := url.Parse("http://localhost")
	event := make(map[string]interface{})

	p := &postProducer{
		Client:   &http.Client{Transport: tr},
		Endpoint: u,
	}

	tr.EXPECT().RoundTrip(gomock.Any()).Return(nil, errors.New("error"))
	_, err := p.Produce(context.Background(), event)
	require.NotNil(t, err)
}

func TestPostProducerBadStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tr := NewMockRoundTripper(ctrl)
	u, _ := url.Parse("http://localhost")
	event := make(map[string]interface{})

	p := &postProducer{
		Client:   &http.Client{Transport: tr},
		Endpoint: u,
	}

	tr.EXPECT().RoundTrip(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       http.NoBody,
	}, nil)
	_, err := p.Produce(context.Background(), event)
	require.NotNil(t, err)
}

func TestPostProducerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tr := NewMockRoundTripper(ctrl)
	u, _ := url.Parse("http://localhost")
	event := make(map[string]interface{})

	p := &postProducer{
		Client:   &http.Client{Transport: tr},
		Endpoint: u,
	}

	tr.EXPECT().RoundTrip(gomock.Any()).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}, nil)
	_, err := p.Produce(context.Background(), event)
	require.Nil(t, err)
}

func TestProducerPOSTComponent_New(t *testing.T) {
	tests := []struct {
		name         string
		conf         *ProducerPOSTConfig
		wantErr      bool
		wantProducer bool
	}{
		{
			name: "missing url",
			conf: func() *ProducerPOSTConfig {
				return NewProducerPOSTComponent().Settings()
			}(),
			wantErr: true,
		},
		{
			name: "invalid url",
			conf: func() *ProducerPOSTConfig {
				conf := NewProducerPOSTComponent().Settings()
				conf.Endpoint = ":/localhost"
				return conf
			}(),
			wantErr: true,
		},
		{
			name: "default http",
			conf: func() *ProducerPOSTConfig {
				conf := NewProducerPOSTComponent().Settings()
				conf.Endpoint = "http://localhost"
				conf.HTTPClient.Type = HTTPTypeDefault
				return conf
			}(),
			wantProducer: true,
		},
		{
			name: "smart http",
			conf: func() *ProducerPOSTConfig {
				conf := NewProducerPOSTComponent().Settings()
				conf.Endpoint = "http://localhost"
				conf.HTTPClient.Type = HTTPTypeSmart
				conf.HTTPClient.Smart.OpenAPI = transportdConfig
				return conf
			}(),
			wantProducer: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewProducerPOSTComponent()
			got, err := c.New(context.Background(), tt.conf)
			if tt.wantErr {
				require.NotNil(t, err)
			}
			if tt.wantProducer {
				require.Nil(t, err)
				require.NotNil(t, got)
			}
		})
	}
}
