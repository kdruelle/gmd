package client

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// CheckUpdate checks if the given container needs to be updated.
// It returns true if an update is needed, false otherwise.
// It also returns an error if an error occurs during the check.
func (c *Client) CheckUpdate(containerID string) (bool, error) {

	container, err := c.ContainerInspect(containerID)
	if err != nil {
		return false, err
	}

	var image image.Summary
	images, err := c.ImageList()

	if err != nil {
		return false, err
	}

	for _, img := range images {
		if img.ID == container.Image {
			image = img
			break
		}
	}

	if image.ID == "" {
		return false, fmt.Errorf("image %s not found", container.Image)
	}
	// log.Printf("check update for container %s:  image %s - %s - %+v", c.Name, c.Image, c.Config.Image)

	// if strings.Contains(imageRef, "@") {
	// 	return false, nil
	// }

	// if strings.HasPrefix(imageRef, "sha256") {
	// 	return true, nil
	// }

	localDigests, err := c.getLocalDigests(container.Config.Image)
	if err != nil {
		return false, err
	}

	remoteDigest, err := getRemoteDigest(container.Config.Image)
	if err != nil {
		log.Printf("image : %s, localDigests: %v, err: %s", container.Image, localDigests, err)
		return false, err
	}

	log.Printf("image : %s, localDigests: %v, remoteDigest: %s", container.Image, localDigests, remoteDigest)

	f := func(s string) bool {
		return strings.HasPrefix(s, remoteDigest) || strings.HasSuffix(s, remoteDigest)
	}

	if slices.ContainsFunc(localDigests, f) {
		return false, nil
	}

	log.Printf("image to update : %s, container: %s, localDigests: %v, remoteDigest: %s", image.ID, container.ID, localDigests, remoteDigest)

	return true, nil
}

func (c *Client) getLocalDigests(imageID string) ([]string, error) {
	//imgInspect, _, err := cli.ImageInspectWithRaw(ctx, imageID)
	imgInspect, err := c.cli.ImageInspect(context.Background(), imageID)
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
		remote.WithPlatform(v1.Platform{Architecture: runtime.GOARCH, OS: runtime.GOOS}),
	)
	if err != nil {
		return "", err
	}

	return desc.Digest.String(), nil
}
