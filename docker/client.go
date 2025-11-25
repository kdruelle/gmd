package docker

import "github.com/docker/docker/client"

type Client struct {
	cli    client.APIClient
	events chan *Event
}

func NewClient() *Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	//TODO: size should not be true, investigate later
	return &Client{
		cli: cli,
	}
}
