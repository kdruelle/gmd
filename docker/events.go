package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

type EventType string

const (
	ContainerEventType      EventType = "container"
	ContainerStatsEventType EventType = "stats"
	ImageEventType          EventType = "image"
)

type Action string

const (
	ActionCreate  Action = Action(events.ActionCreate)
	ActionStart   Action = Action(events.ActionStart)
	ActionRestart Action = Action(events.ActionRestart)
	ActionStop    Action = Action(events.ActionStop)
	ActionRemove  Action = Action(events.ActionRemove)
	ActionDie     Action = Action(events.ActionDie)
	ActionKill    Action = Action(events.ActionKill)
	ActionPause   Action = Action(events.ActionPause)
	ActionUnPause Action = Action(events.ActionUnPause)
	ActionRename  Action = Action(events.ActionRename)
	ActionDestroy Action = Action(events.ActionDestroy)
	ActionPush    Action = Action(events.ActionPush)
	ActionPull    Action = Action(events.ActionPull)
	ActionPrune   Action = Action(events.ActionPrune)
	ActionDelete  Action = Action(events.ActionDelete)
)

type Event struct {
	EventType EventType
	Action
	ActorID string
}

func (m *Monitor) Events() <-chan Event {
	return m.events
}

func eventsFilter() filters.Args {
	filter := filters.NewArgs()
	filter.Add("type", string(events.ContainerEventType))
	filter.Add("type", string(events.ImageEventType))

	for _, ev := range []events.Action{
		events.ActionCreate,
		events.ActionStart,
		events.ActionRestart,
		events.ActionStop,
		events.ActionRemove,
		events.ActionDie,
		events.ActionKill,
		events.ActionPause,
		events.ActionUnPause,
		events.ActionRename,
		events.ActionDestroy,
		events.ActionPush,
		events.ActionPull,
		events.ActionPrune,
		events.ActionDelete,
	} {
		filter.Add("event", string(ev))
	}

	return filter
}

func (m *Monitor) listenEvents() {

	for {
		select {
		case msg := <-m.eventsCh:
			if ev, err := m.handleEvent(msg); err == nil {
				m.events <- ev
			}
		case <-m.errsCh:
			//	m.errsCh <- err
			return
		}
	}
}
func (m *Monitor) handleEvent(e events.Message) (Event, error) {

	switch e.Type {
	case events.ContainerEventType:
		m.refreshContainer(e.Actor.ID)

		return Event{
			EventType: ContainerEventType,
			ActorID:   e.Actor.ID,
		}, nil

	case events.ImageEventType:
		m.refreshImage(e.Actor.ID)

		return Event{
			EventType: ImageEventType,
			ActorID:   e.Actor.ID,
		}, nil
	}

	return Event{}, fmt.Errorf("unhandled event")
}

func (c *Client) Events(ctx context.Context) (<-chan *Event, <-chan error) {

	filters := filters.NewArgs()
	filters.Add("type", string(events.ContainerEventType))
	filters.Add("type", string(events.ImageEventType))
	filters.Add("type", string(events.VolumeEventType))

	filters.Add("event", string(events.ActionCreate))
	filters.Add("event", string(events.ActionStart))
	filters.Add("event", string(events.ActionRestart))
	filters.Add("event", string(events.ActionStop))
	filters.Add("event", string(events.ActionRemove))
	filters.Add("event", string(events.ActionDie))
	filters.Add("event", string(events.ActionKill))
	filters.Add("event", string(events.ActionPause))
	filters.Add("event", string(events.ActionUnPause))
	filters.Add("event", string(events.ActionRename))
	filters.Add("event", string(events.ActionDestroy))

	filters.Add("event", string(events.ActionPush))
	filters.Add("event", string(events.ActionPull))
	filters.Add("event", string(events.ActionPrune))
	filters.Add("event", string(events.ActionDelete))

	// ActionCheckpoint   Action = "checkpoint"
	// ActionUpdate       Action = "update"
	// ActionOOM          Action = "oom"
	// ActionCommit       Action = "commit"
	// ActionEnable       Action = "enable"
	// ActionDisable      Action = "disable"
	// ActionConnect      Action = "connect"
	// ActionDisconnect   Action = "disconnect"
	// ActionReload       Action = "reload"
	// ActionMount        Action = "mount"
	// ActionUnmount      Action = "unmount"

	de, errors := c.cli.Events(ctx, events.ListOptions{
		Filters: filters,
	})

	c.events = make(chan *Event, 20)

	go func() {
		for e := range de {
			log.Printf("docker event received: %+v", e)
			c.handleEvent(e)
		}
		log.Printf("docker event watcher stopped")
	}()

	return c.events, errors

}

func (c *Client) handleEvent(msg events.Message) {

	event := &Event{
		EventType: EventType(msg.Type),
		ActorID:   string(msg.Actor.ID),
	}

	c.events <- event

}
