package docker

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"slices"
	"strings"

	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	gcr "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func (m *Monitor) checkUpdate(c *Container) {

	m.mu.RLock()
	image := m.images[c.Image]
	m.mu.RUnlock()

	// log.Printf("check update for container %s:  image %s - %s - %+v", c.Name, c.Image, c.Config.Image)

	// if strings.Contains(imageRef, "@") {
	// 	return false, nil
	// }

	// if strings.HasPrefix(imageRef, "sha256") {
	// 	return true, nil
	// }

	localDigests, err := getLocalDigests(c.Config.Image)
	if err != nil {
		return
	}

	remoteDigest, err := getRemoteDigest(c.Config.Image)
	if err != nil {
		log.Printf("image : %s, localDigests: %v, err: %s", c.Image, localDigests, err)
		return
	}

	//log.Printf("image : %s, localDigests: %v, remoteDigest: %s", c.Image, localDigests, remoteDigest)

	f := func(s string) bool {
		return strings.HasPrefix(s, remoteDigest) || strings.HasSuffix(s, remoteDigest)
	}

	if slices.ContainsFunc(localDigests, f) {
		return
	}

	log.Printf("image to update : %s, container: %s, localDigests: %v, remoteDigest: %s", image.ID, c.ID, localDigests, remoteDigest)
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	c.Update = true
	e := Event{
		EventType: ContainerEventType,
		ActorID:   c.ID,
	}
	m.events <- e
}

func getLocalDigests(imageID string) ([]string, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Failed to create docker client")
	}
	ctx := context.Background()
	//imgInspect, _, err := cli.ImageInspectWithRaw(ctx, imageID)
	imgInspect, err := cli.ImageInspect(ctx, imageID)
	if err != nil {
		return nil, err
	}
	if len(imgInspect.RepoDigests) == 0 {
		return nil, fmt.Errorf("pas de RepoDigests pour %s", imageID)
	}

	return imgInspect.RepoDigests, nil
}

func getRemoteDigest(image string) (string, error) {

	log.Printf("getRemoteDigest for %s", image)

	ref, err := name.ParseReference(image)
	if err != nil {
		return "", err
	}

	// HEAD request for manifest digest
	desc, err := remote.Head(ref,
		remote.WithPlatform(gcr.Platform{Architecture: runtime.GOARCH, OS: runtime.GOOS}),
	)
	if err != nil {
		return "", err
	}

	return desc.Digest.String(), nil
}
