package components

import (
	"context"
	"net/http"
)

// Producer is the interface used for sending all events to a destination.
type Producer interface {
	// Produce ships the event to the destination and returns the final
	// version of the data sent.
	Produce(ctx context.Context, event interface{}) (interface{}, error)
}

// RoundTripper is the interface that handles all HTTP operations. It is almost
// exclusively used with an http.Client wrapped around it. This is included here
// for documentation purposes only.
type RoundTripper = http.RoundTripper
