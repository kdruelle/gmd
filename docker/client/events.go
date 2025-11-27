package client

import (
	"context"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

func (c *Client) StartEvents() (<-chan events.Message, <-chan error) {

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

	c.eventsContext, c.eventsCancel = context.WithCancel(context.Background())

	ev, errors := c.cli.Events(c.eventsContext, events.ListOptions{
		Filters: filters,
	})

	return ev, errors

}

func (c *Client) StopEvents() {
	c.eventsCancel()
	<-c.eventsContext.Done()
}
