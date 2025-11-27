package cache

import (
	"sync"

	"github.com/docker/docker/api/types/events"
	"github.com/kdruelle/gmd/docker/client"
)

type Cache struct {
	cli *client.Client
	mu  sync.RWMutex

	images     map[string]*client.Image
	containers map[string]*client.Container

	ievents <-chan events.Message
	ierrors <-chan error

	events            chan Event
	containerDeletion chan string
}

func NewCache(cli *client.Client) *Cache {
	c := &Cache{
		cli:               cli,
		images:            make(map[string]*client.Image),
		containers:        make(map[string]*client.Container),
		events:            make(chan Event, 20),
		containerDeletion: make(chan string, 20),
	}

	return c
}

func (c *Cache) LoadAndStart() error {

	c.ievents, c.ierrors = c.cli.StartEvents()

	imgs := c.snapshotImages()
	c.mu.Lock()
	for _, img := range imgs {
		c.images[img.ID] = img
	}
	c.mu.Unlock()

	c.events <- Event{EventType: ImagesLoadedEventType}

	conts := c.snapshotContainers()
	c.mu.Lock()
	for _, cont := range conts {
		c.containers[cont.ID] = cont
	}
	c.mu.Unlock()

	c.events <- Event{EventType: ContainersLoadedEventType}

	go c.listenEvents()
	go c.containerDeleteWorker()

	// for i := range m.containers {
	// 	go m.checkUpdate(m.containers[i])
	// 	go m.watchStats(m.containers[i])
	// }

	return nil
}
