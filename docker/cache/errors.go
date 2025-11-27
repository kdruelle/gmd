package cache

import "fmt"

var (
	ErrContainerNotFound = fmt.Errorf("container not found")
	ErrImageNotFound     = fmt.Errorf("image not found")
)
