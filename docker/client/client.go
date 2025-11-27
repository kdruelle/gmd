package client

import (
	"context"

	"github.com/docker/docker/client"
)

// Client represents a client to the Docker daemon.
//
// It provides a way to interact with the daemon and receive events
// from the daemon.
type Client struct {
	cli           client.APIClient   // cli is the underlying client to the Docker daemon.
	eventsContext context.Context    // eventsContext is the context used for listening to events from the daemon.
	eventsCancel  context.CancelFunc // eventsCancel is the cancel function for the events context.
}

// NewClient returns a new Client object, which represents a client to the Docker daemon.
// If the creation of the client fails, it returns nil and an error.
func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
}
