package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

type Monitor struct {
	cli *client.Client

	images     map[string]*Image
	containers map[string]*Container

	mu sync.RWMutex

	eventsCh <-chan events.Message
	errsCh   <-chan error

	events chan Event
}

func NewMonitor() (*Monitor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	m := &Monitor{
		cli:        cli,
		images:     make(map[string]*Image),
		containers: make(map[string]*Container),
		eventsCh:   make(chan events.Message),
		errsCh:     make(chan error),
		events:     make(chan Event, 200),
	}

	return m, nil
}

func (m *Monitor) Start() {
	ctx := context.Background()
	eventsCh, errsCh := m.cli.Events(ctx, events.ListOptions{Filters: eventsFilter()})
	m.eventsCh = eventsCh
	m.errsCh = errsCh

	drained := m.drainEvents()

	if err := m.loadInitialState(); err != nil {
		panic(err)
	}

	for _, ev := range drained {
		_, _ = m.handleEvent(ev)
	}

	go m.listenEvents()
	m.events <- Event{EventType: ImageEventType}
	m.events <- Event{EventType: ContainerEventType}

	for i := range m.containers {
		go m.checkUpdate(m.containers[i])
		go m.watchStats(m.containers[i])
	}
}

func (m *Monitor) drainEvents() []events.Message {
	drained := []events.Message{}
	for {
		select {
		case ev := <-m.eventsCh:
			drained = append(drained, ev)
		default:
			return drained
		}
	}
}

func (m *Monitor) loadInitialState() error {

	m.mu.Lock()
	defer m.mu.Unlock()

	imgs := m.snapshotImages()
	for _, img := range imgs {
		m.images[img.ID] = img
	}

	conts := m.snapshotContainers()
	for _, c := range conts {
		m.containers[c.ID] = c
	}

	return nil
}
