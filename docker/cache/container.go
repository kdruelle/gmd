package cache

import (
	"log"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/kdruelle/gmd/docker/client"
)

func (c *Cache) Containers() []client.Container {
	out := make([]client.Container, 0, len(c.containers))

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, c := range c.containers {
		out = append(out, *c)
	}
	return out
}

func (c *Cache) Container(id string) (client.Container, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c, ok := c.containers[id]; ok {
		return *c, nil
	}
	return client.Container{}, ErrContainerNotFound
}

func (c *Cache) refreshContainer(ev events.Message) {

	cont, err := c.cli.ContainerInspect(ev.Actor.ID)

	c.mu.Lock()
	defer c.mu.Unlock()

	if err == nil {

		summary := &client.Container{
			InspectResponse: cont,
		}
		c.containers[ev.Actor.ID] = summary

		if ev.Action == events.ActionDestroy {
			log.Printf("try to wait for container %s deletion", ev.Actor.ID)
			c.containerDeletion <- ev.Actor.ID
		}

	} else {
		log.Printf("refresh container %s, delete container: %v", ev.Actor.ID, err)
		delete(c.containers, ev.Actor.ID)
	}
}

func (c *Cache) snapshotContainers() []*client.Container {
	ctnrs, err := c.cli.ContainerList()
	if err != nil {
		panic(err)
	}

	containers := make([]*client.Container, len(ctnrs))

	for i, container := range ctnrs {
		inspect, err := c.cli.ContainerInspect(container.ID)
		if err != nil {
			panic(err)
		}
		containers[i] = &client.Container{
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
