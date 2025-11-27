package client

import (
	"context"

	"github.com/docker/docker/client"
)

type Client struct {
	cli client.APIClient

	// events
	eventsContext context.Context
	eventsCancel  context.CancelFunc
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
}
