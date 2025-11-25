package docker

import (
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/hashicorp/go-version"
)

type Container struct {
	container.InspectResponse
	Stats  container.StatsResponse
	Update bool
}

func (m *Monitor) Containers() []Container {
	out := make([]Container, 0, len(m.containers))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, c := range m.containers {
		out = append(out, *c)
	}
	return out
}

func (m *Monitor) Container(id string) (Container, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, ok := m.containers[id]; ok {
		return *c, nil
	}
	return Container{}, ErrContainerNotFound
}

func (m *Monitor) StartContainer(id string) error {
	return m.cli.ContainerStart(context.Background(), id, container.StartOptions{})
}

func (m *Monitor) StopContainer(id string) error {
	return m.cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (m *Monitor) RestartContainer(id string) error {
	return m.cli.ContainerRestart(context.Background(), id, container.StopOptions{})
}

func (m *Monitor) DeleteContainer(id string) error {
	dockerOpts := container.RemoveOptions{}
	return m.cli.ContainerRemove(context.Background(), id, dockerOpts)
}

func (m *Monitor) CreateContainerFromConfig(config *container.InspectResponse) (container.CreateResponse, error) {

	info, err := m.cli.ServerVersion(context.Background())
	if err != nil {
		log.Fatal("Failed to get docker version")
		return container.CreateResponse{}, err
	}

	sanitizeContainerJONVersion(config, info.APIVersion)

	netConfig := &network.NetworkingConfig{
		EndpointsConfig: config.NetworkSettings.Networks,
	}

	r, err := m.cli.ContainerCreate(context.Background(), config.Config, config.HostConfig, netConfig, nil, config.Name)
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

func (m *Monitor) refreshContainer(id string) {

	c, err := m.cli.ContainerInspect(context.Background(), id)

	m.mu.Lock()
	defer m.mu.Unlock()

	if err == nil {
		summary := &Container{
			InspectResponse: c,
		}
		if oc, ok := m.containers[id]; ok {
			// container already exists
			summary.Update = oc.Update
		} else {
			go m.watchStats(summary)
			go m.checkUpdate(summary)
		}
		m.containers[id] = summary
	} else {
		delete(m.containers, id)
	}
}

func (c *Monitor) snapshotContainers() []*Container {
	ctnrs, err := c.cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	containers := make([]*Container, len(ctnrs))

	for i, container := range ctnrs {
		inspect, err := c.cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			panic(err)
		}
		containers[i] = &Container{
			InspectResponse: inspect,
		}
	}

	return containers
}

func (m *Monitor) GetContainerRawConfig(containerID string) (*container.InspectResponse, error) {
	c, err := m.cli.ContainerInspect(context.Background(), containerID)
	return &c, err
}

// // PullImage is a convenience wrapper that discards pull progress.
// func (dc *DockerClient) PullImage(ctx context.Context, imageRef string) error {
// 	return dc.PullImageWithProgress(ctx, imageRef, func(map[string]any) {})
// }

// func (dc *DockerClient) GetContainerRawConfig(containerID string) (*et.ContainerJSON, error) {
// 	c, err := dc.cli.ContainerInspect(context.Background(), containerID)
// 	return &c, err
// }

// func (dc *DockerClient) CreateContainerFromConfig(config *et.ContainerJSON) (container.CreateResponse, error) {

// 	info, err := dc.cli.ServerVersion(context.Background())
// 	if err != nil {
// 		log.Fatal("Failed to get docker version")
// 		return container.CreateResponse{}, err
// 	}

// 	SanitizeContainerJONVersion(config, info.APIVersion)

// 	netConfig := &network.NetworkingConfig{
// 		EndpointsConfig: config.NetworkSettings.Networks,
// 	}

// 	r, err := dc.cli.ContainerCreate(context.Background(), config.Config, config.HostConfig, netConfig, nil, config.Name)
// 	return r, err
// }

// func SanitizeContainerJONVersion(containerJson *et.ContainerJSON, apiVersionString string) {

// 	apiVersion, err := version.NewVersion(apiVersionString)
// 	if err != nil {
// 		return
// 	}

// 	if apiVersion.LessThan(version.Must(version.NewVersion("1.44"))) {
// 		for netName, netConf := range containerJson.NetworkSettings.Networks {
// 			netConf.MacAddress = ""
// 			containerJson.NetworkSettings.Networks[netName] = netConf
// 		}
// 	}

// 	if containerJson.HostConfig.NetworkMode == "host" || strings.HasPrefix(string(containerJson.HostConfig.NetworkMode), "container:") {
// 		containerJson.Config.Hostname = ""
// 	}

// }
