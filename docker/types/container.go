package types

import "github.com/docker/docker/api/types/container"

type Container struct {
	container.InspectResponse
}
