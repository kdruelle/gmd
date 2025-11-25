package containerupdate

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/kdruelle/gmd/docker"
	style "github.com/kdruelle/gmd/tui/styles"
)

type ControllerUpdateMsg struct {
}

type createResponse struct {
	Resp container.CreateResponse
	Err  error
}

type Controller struct {
	m          sync.RWMutex
	client     *docker.Monitor
	updateChan chan ControllerUpdateMsg

	order  []string
	layers map[string]string
	lines  []string
}

func New(client *docker.Monitor, updateChan chan ControllerUpdateMsg) *Controller {
	c := Controller{
		client:     client,
		updateChan: updateChan,
	}
	return &c
}

func (c *Controller) GetLines() []string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.lines
}

func (c *Controller) StartUpdate(container docker.Container) {
	c.order = []string{}
	c.layers = make(map[string]string)
	go c.updateContainer(container)
}

func (c *Controller) updateContainer(container docker.Container) {

	done := make(chan error)
	defer close(done)

	err := c.client.PullImageWithProgress(context.Background(), container.Config.Image, func(msg map[string]interface{}) {
		var ok bool
		var status, layerId string

		if status, ok = msg["status"].(string); !ok {
			return
		}

		if layerId, ok = msg["id"].(string); !ok {
			return
		}

		if layerId == "" {
			layerId = fmt.Sprintf("general-%d", len(c.layers)) // évite collision
		}

		line := status
		if progress, ok := msg["progress"].(string); ok {
			line += " " + progress
		}

		c.m.Lock()
		defer c.m.Unlock()
		if _, exists := c.layers[layerId]; !exists {
			c.order = append(c.order, layerId) // première fois qu’on voit ce layer
		}
		c.layers[layerId] = line

		c.lines = c.lines[:0]
		for _, id := range c.order {
			c.lines = append(c.lines, c.layers[id])
		}

		log.Printf("line: %s", line)

		c.updateChan <- ControllerUpdateMsg{}
	})

	if err != nil {
		log.Printf("Error pull for image %s : %v", container.Config.Image, err)
		c.m.Lock()
		c.lines = append(c.lines, fmt.Sprintf("Error pull image: %v", err))
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
		return
	}

	containerConfig, err := c.client.GetContainerRawConfig(container.ID)
	if err != nil {
		log.Printf("Error get config for container %s : %v", container.ID, err)
		c.m.Lock()
		c.lines = append(c.lines, fmt.Sprintf("Error get config: %v", err))
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
		return
	}

	add := true
	err = spinUntilDone(func() error {
		return c.client.StopContainer(container.ID)
	}, func(frame string) {
		c.m.Lock()
		if add {
			c.lines = append(c.lines, fmt.Sprintf("%s Stoping container: %s", frame, containerConfig.Name))
			add = false
		} else {
			c.lines[len(c.lines)-1] = fmt.Sprintf("%s Stoping container: %s", frame, containerConfig.Name)
		}
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
	})

	if err != nil {
		c.m.Lock()
		c.lines = append(c.lines, fmt.Sprintf("Error stop: %v", err))
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
		return
	}

	c.m.Lock()
	c.lines[len(c.lines)-1] = fmt.Sprintf("%s Stoping container: %s", style.ActiveItem.Render("✓"), containerConfig.Name)
	c.m.Unlock()
	c.updateChan <- ControllerUpdateMsg{}

	add = true
	err = spinUntilDone(func() error {
		return c.client.DeleteContainer(container.ID)
	}, func(frame string) {
		c.m.Lock()
		if add {
			c.lines = append(c.lines, fmt.Sprintf("%s Removing container: %s", frame, containerConfig.Name))
			add = false
		} else {
			c.lines[len(c.lines)-1] = fmt.Sprintf("%s Removing container: %s", frame, containerConfig.Name)
		}
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
	})

	if err != nil {
		c.m.Lock()
		c.lines = append(c.lines, fmt.Sprintf("Error remove: %v", err))
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
		return
	}

	c.m.Lock()
	c.lines[len(c.lines)-1] = fmt.Sprintf("%s Removing container: %s", style.ActiveItem.Render("✓"), containerConfig.Name)
	c.m.Unlock()
	c.updateChan <- ControllerUpdateMsg{}

	add = true
	cr := spinUntilDone(func() createResponse {
		r, e := c.client.CreateContainerFromConfig(containerConfig)
		return createResponse{Resp: r, Err: e}
	}, func(frame string) {
		c.m.Lock()
		if add {
			c.lines = append(c.lines, fmt.Sprintf("%s Creating container: %s", frame, containerConfig.Name))
			add = false
		} else {
			c.lines[len(c.lines)-1] = fmt.Sprintf("%s Creating container: %s", frame, containerConfig.Name)
		}
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
	})

	err = cr.Err

	if err != nil {
		c.m.Lock()
		c.lines = append(c.lines, fmt.Sprintf("Error create: %v", err))
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
		return
	}

	c.m.Lock()
	c.lines[len(c.lines)-1] = fmt.Sprintf("%s Creating container: %s", style.ActiveItem.Render("✓"), containerConfig.Name)
	c.m.Unlock()
	c.updateChan <- ControllerUpdateMsg{}

	add = true
	err = spinUntilDone(func() error {
		return c.client.StartContainer(cr.Resp.ID)
	}, func(frame string) {
		c.m.Lock()
		if add {
			c.lines = append(c.lines, fmt.Sprintf("%s Starting container: %s", frame, containerConfig.Name))
			add = false
		} else {
			c.lines[len(c.lines)-1] = fmt.Sprintf("%s Starting container: %s", frame, containerConfig.Name)
		}
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
	})

	if err != nil {
		c.m.Lock()
		c.lines = append(c.lines, fmt.Sprintf("Error start: %v", err))
		c.m.Unlock()
		c.updateChan <- ControllerUpdateMsg{}
		return
	}

	c.m.Lock()
	c.lines[len(c.lines)-1] = fmt.Sprintf("%s Starting container: %s", style.ActiveItem.Render("✓"), containerConfig.Name)
	c.m.Unlock()
	c.updateChan <- ControllerUpdateMsg{}

	// notification := listStatusMessageStyle.Render(fmt.Sprintf("Updated %s", containerConfig.Name))
	// notificationChan <- NewNotification(activeTab, notification)

	// updateTracker.Delete(item.Image)
	// progressTracker.Clear(item.GetId())

	// time.Sleep(1 * time.Second)

	// notifyUIRefresh(refreshChan, UiRefreshEvent{
	// 	Type: "all",
	// })

	// return
}

func spinUntilDone[T any](
	action func() T,
	updateLine func(frame string),
) T {

	done := make(chan T)

	// Lance l’action en arrière-plan
	go func() {
		done <- action()
	}()

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	index := 0

	for {
		select {
		case r := <-done:
			// Terminé
			return r

		case <-time.After(100 * time.Millisecond):
			// Frame suivante
			frame := spinner[index]
			index = (index + 1) % len(spinner)

			updateLine(style.Spinner.Render(frame))
		}
	}
}
