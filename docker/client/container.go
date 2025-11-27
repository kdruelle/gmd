package client

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/hashicorp/go-version"
)

type Container struct {
	container.InspectResponse
}

func (c *Client) StartContainer(id string) error {
	return c.cli.ContainerStart(context.Background(), id, container.StartOptions{})
}

func (c *Client) StopContainer(id string) error {
	return c.cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (c *Client) RestartContainer(id string) error {
	return c.cli.ContainerRestart(context.Background(), id, container.StopOptions{})
}

func (c *Client) DeleteContainer(id string) error {
	dockerOpts := container.RemoveOptions{}
	return c.cli.ContainerRemove(context.Background(), id, dockerOpts)
}

func (c *Client) ContainerInspect(id string) (container.InspectResponse, error) {
	return c.cli.ContainerInspect(context.Background(), id)
}

func (c *Client) ContainerStats(id string) (container.StatsResponse, error) {
	var v container.StatsResponse
	stats, err := c.cli.ContainerStatsOneShot(context.Background(), id)

	if err != nil {
		return v, err
	}

	dec := json.NewDecoder(stats.Body)
	err = dec.Decode(&v)
	return v, err
}

func (c *Client) GetContainerRawConfig(containerID string) (*container.InspectResponse, error) {
	config, err := c.cli.ContainerInspect(context.Background(), containerID)
	return &config, err
}

func (c *Client) CreateContainerFromConfig(config *container.InspectResponse) (container.CreateResponse, error) {

	info, err := c.cli.ServerVersion(context.Background())
	if err != nil {
		log.Fatal("Failed to get docker version")
		return container.CreateResponse{}, err
	}

	sanitizeContainerJONVersion(config, info.APIVersion)

	netConfig := &network.NetworkingConfig{
		EndpointsConfig: config.NetworkSettings.Networks,
	}

	r, err := c.cli.ContainerCreate(context.Background(), config.Config, config.HostConfig, netConfig, nil, config.Name)
	return r, err
}

func (c *Client) ContainerList() ([]container.Summary, error) {
	return c.cli.ContainerList(context.Background(), container.ListOptions{All: true})
}

func sanitizeContainerJONVersion(containerJson *container.InspectResponse, apiVersionString string) {

	apiVersion, err := version.NewVersion(apiVersionString)
	if err != nil {
		log.Printf("Failed to get docker version")
		return
	}

	if apiVersion.LessThan(version.Must(version.NewVersion("1.44"))) {
		for netName, netConf := range containerJson.NetworkSettings.Networks {
			netConf.MacAddress = ""
			containerJson.NetworkSettings.Networks[netName] = netConf
		}
	}

	if containerJson.HostConfig.NetworkMode == "host" || strings.HasPrefix(string(containerJson.HostConfig.NetworkMode), "container:") {
		containerJson.Config.Hostname = ""

		containerJson.HostConfig.PortBindings = nil
		containerJson.Config.ExposedPorts = nil
		containerJson.HostConfig.PublishAllPorts = false
	}

}
