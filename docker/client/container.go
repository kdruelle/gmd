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

// ContainerList returns a list of all containers on the docker daemon.
// The function returns an error if the list of containers could not be retrieved.
// The returned list of containers is a slice of container.Summary objects.
// The container.Summary objects contain only the most basic information about the container, such as its ID, name, and status.
// The container.Summary objects are returned in a random order.
func (c *Client) ContainerList() ([]container.Summary, error) {
	return c.cli.ContainerList(context.Background(), container.ListOptions{All: true})
}

// StartContainer starts a container with the given ID.
// It returns an error if the container could not be started.
func (c *Client) StartContainer(id string) error {
	return c.cli.ContainerStart(context.Background(), id, container.StartOptions{})
}

// StopContainer stops a container with the given ID.
// It returns an error if the container could not be stopped.
func (c *Client) StopContainer(id string) error {
	return c.cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

// RestartContainer restarts a container with the given ID.
// It returns an error if the container could not be restarted.
func (c *Client) RestartContainer(id string) error {
	return c.cli.ContainerRestart(context.Background(), id, container.StopOptions{})
}

// DeleteContainer deletes a container with the given ID.
// It returns an error if the container could not be deleted.
func (c *Client) DeleteContainer(id string) error {
	dockerOpts := container.RemoveOptions{}
	return c.cli.ContainerRemove(context.Background(), id, dockerOpts)
}

// ContainerInspect returns the configuration of the container with the given ID.
// It returns an error if the container could not be inspected.
func (c *Client) ContainerInspect(id string) (container.InspectResponse, error) {
	return c.cli.ContainerInspect(context.Background(), id)
}

// ContainerStats returns the stats of a container with the given ID.
// It returns an error if the container could not be inspected.
// The returned stats are the result of a single shot stats query, and are not
// streamed. If the container is not running, the stats will be empty.
// If the container does not exist, an error will be returned.
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

// CreateContainerFromConfig creates a container based on the given container configuration.
// It returns an error if the container could not be created.
// The given container configuration is expected to be a container.InspectResponse object.
// The created container will have the same configuration as the given container.
// The function will sanitize the given container configuration to make it compatible with
// the docker daemon API version.
// The function will return a container.CreateResponse object containing information about the created container.
func (c *Client) CreateContainerFromConfig(config container.InspectResponse) (container.CreateResponse, error) {

	info, err := c.cli.ServerVersion(context.Background())
	if err != nil {
		log.Fatal("Failed to get docker version")
		return container.CreateResponse{}, err
	}

	sanitizeContainerJONVersion(&config, info.APIVersion)

	netConfig := &network.NetworkingConfig{
		EndpointsConfig: config.NetworkSettings.Networks,
	}

	r, err := c.cli.ContainerCreate(context.Background(), config.Config, config.HostConfig, netConfig, nil, config.Name)
	return r, err
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
