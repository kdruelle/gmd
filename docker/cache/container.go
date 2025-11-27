package cache

import (
	"log"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/kdruelle/gmd/docker/types"
)

// Containers returns all containers in the cache.
// The function locks the cache for reading and returns a slice of all containers.
// The returned slice is a copy of the underlying data, so it can be safely used without taking a write lock on the cache.
// The function does not return an error. If the cache is empty, an empty slice is returned.
func (c *Cache) Containers() []types.Container {
	out := make([]types.Container, 0, len(c.containers))

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, c := range c.containers {
		out = append(out, *c)
	}
	return out
}

// Container returns the container with the given ID from the cache.
// The function locks the cache for reading and returns a copy of the underlying data, so it can be safely used without taking a write lock on the cache.
// If the container is not found, ErrContainerNotFound is returned.
func (c *Cache) Container(id string) (types.Container, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c, ok := c.containers[id]; ok {
		return *c, nil
	}
	return types.Container{}, ErrContainerNotFound
}

// refreshContainer refreshes the cache with the given container event.
// It locks the cache for writing and updates the container with the given ID.
// If the event is an action of type "destroy", the container is removed from the cache.
// If an error occurs while refreshing the container, the container is removed from the cache.
// The function does not return an error.
func (c *Cache) refreshContainer(ev events.Message) {

	cont, err := c.cli.ContainerInspect(ev.Actor.ID)

	c.mu.Lock()

	if err == nil {

		summary := &types.Container{
			InspectResponse: cont,
		}
		c.containers[ev.Actor.ID] = summary

		if ev.Action == events.ActionDestroy {
			c.mu.Unlock()
			log.Printf("try to wait for container %s deletion", ev.Actor.ID)
			c.containerDeletion <- ev.Actor.ID
			return
		}

	} else {
		log.Printf("refresh container %s, delete container: %v", ev.Actor.ID, err)
		delete(c.containers, ev.Actor.ID)
	}
	c.mu.Unlock()
}

func (c *Cache) snapshotContainers() []*types.Container {
	ctnrs, err := c.cli.ContainerList()
	if err != nil {
		panic(err)
	}

	containers := make([]*types.Container, len(ctnrs))

	for i, container := range ctnrs {
		inspect, err := c.cli.ContainerInspect(container.ID)
		if err != nil {
			panic(err)
		}
		containers[i] = &types.Container{
			InspectResponse: inspect,
		}
	}

	return containers
}

func (c *Cache) containerDeleteWorker() {
	for id := range c.containerDeletion {

		for range 25 { // max 5 secondes
			cont, err := c.cli.ContainerInspect(id)
			if err != nil {
				log.Printf("delete container %s: %v", id, err)
				c.mu.Lock()
				delete(c.containers, id)
				c.mu.Unlock()

				c.events <- Event{EventType: ContainerEventType, ActorID: id}
				break
			}
			log.Printf("deleted container %s is still there: %+v", id, cont)
			time.Sleep(200 * time.Millisecond)
		}
	}
}
