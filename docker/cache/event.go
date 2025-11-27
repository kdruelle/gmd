package cache

import (
	"fmt"
	"log"

	"github.com/docker/docker/api/types/events"
)

type EventType string

const (
	ImagesLoadedEventType     EventType = "images-loaded"
	ImageEventType            EventType = EventType(events.ImageEventType)
	ContainersLoadedEventType EventType = "containers-loaded"
	ContainerEventType        EventType = EventType(events.ContainerEventType)
)

type Event struct {
	EventType EventType
	ActorID   string
}

// Events returns a channel of Event objects.
// The channel is populated with events from the Docker daemon and
// is used to notify the UI of changes to the container and image lists.
func (c *Cache) Events() <-chan Event {
	return c.events
}

func (c *Cache) listenEvents() {
	for {
		select {
		case msg, ok := <-c.ievents:
			if !ok {
				return
			}
			log.Printf("lib docker - received event: %+v", msg)
			if ev, err := c.handleEvent(msg); err == nil {
				c.events <- ev
			}
		case <-c.ierrors:
			//	m.errsCh <- err
			return
		}
	}
}

func (c *Cache) handleEvent(e events.Message) (Event, error) {

	switch e.Type {
	case events.ContainerEventType:
		log.Printf("lib docker - received container event: %+v", e)
		c.refreshContainer(e)

		return Event{
			EventType: ContainerEventType,
			ActorID:   e.Actor.ID,
		}, nil

	case events.ImageEventType:
		log.Printf("lib docker - received image event: %+v", e)
		c.refreshImage(e.Actor.ID)

		return Event{
			EventType: ImageEventType,
			ActorID:   e.Actor.ID,
		}, nil
	}

	return Event{}, fmt.Errorf("unhandled event")
}
