package cache

import (
	"log"
	"slices"
	"strings"

	"github.com/kdruelle/gmd/docker/types"
)

// Images returns the list of images from the cache.
// The function locks the cache for reading and returns a copy of the underlying data, so it can be safely used without taking a write lock on the cache.
func (c *Cache) Images() []types.Image {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]types.Image, len(c.images))
	i := 0
	for _, img := range c.images {
		out[i] = *img
		i++
	}

	slices.SortFunc(out, func(a, b types.Image) int {
		tagA := a.Tag()
		tagB := b.Tag()
		return strings.Compare(tagA, tagB)
	})

	return out
}

// Image returns the image with the given ID from the cache.
// The function locks the cache for reading and returns a copy of the underlying data, so it can be safely used without taking a write lock on the cache.
// If the image is not found, ErrImageNotFound is returned.
// If the cache is empty, an empty slice is returned.
func (c *Cache) Image(id string) (types.Image, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c, ok := c.images[id]; ok {
		return *c, nil
	}
	return types.Image{}, ErrImageNotFound
}

// ImagesUnused returns the list of unused images from the cache.
// An image is considered unused if no container is using it.
// The function locks the cache for reading and returns a copy of the underlying data, so it can be safely used without taking a write lock on the cache.
// The returned slice is sorted by image name.
// If the cache is empty, an empty slice is returned.
func (c *Cache) ImagesUnused() []types.Image {
	c.mu.RLock()
	defer c.mu.RUnlock()

	used := make(map[string]bool, len(c.containers))
	for _, c := range c.containers {
		if c.Image != "" {
			used[c.Image] = true
		}
	}

	out := make([]types.Image, 0, len(c.images))
	for _, img := range c.images {

		if _, ok := used[img.ID]; !ok {
			out = append(out, *img)
		}

	}

	slices.SortFunc(out, func(a, b types.Image) int {
		tagA := a.Tag()
		tagB := b.Tag()
		return strings.Compare(tagA, tagB)
	})
	return out

}

func (c *Cache) refreshImage(id string) {
	imgs, err := c.cli.ImageList()
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.images, id)

	// faut retrouver l’image par ID
	for _, img := range imgs {
		if img.ID == id {
			c.images[id] = &types.Image{
				ID:          img.ID,
				RepoTags:    img.RepoTags,
				RepoDigests: img.RepoDigests,
				Size:        img.Size,
				ParentID:    img.ParentID,
			}
			return
		}
	}

}

// snapshotImages returns a snapshot of the images in the cache.
// The function first lists all images with the cli.ImageList() function,
// then creates a map of the images by ID. It then iterates over the list
// of images and adds each image to the map, ignoring any images that
// already exist in the map.
// Next, it iterates over the list of images again and updates the
// parents of each image by looking up the parent ID in the map.
// Finally, it flattens the map into a slice and returns the slice.
func (c *Cache) snapshotImages() []*types.Image {
	list, err := c.cli.ImageList()
	if err != nil {
		panic(err)
	}

	out := make(map[string]*types.Image)

	addImg := func(id string, tags []string, digs []string, size int64, parent string) {
		if _, ok := out[id]; ok {
			return
		}
		out[id] = &types.Image{
			ID:          id,
			RepoTags:    tags,
			RepoDigests: digs,
			Size:        size,
			ParentID:    parent,
		}
	}

	// 1. Add stadards images
	for _, img := range list {
		addImg(img.ID, img.RepoTags, img.RepoDigests, img.Size, img.ParentID)
	}

	// 2. Add parents via history
	for _, img := range list {
		history, err := c.cli.ImageHistory(img.ID)
		if err != nil {
			log.Println("history:", err)
			continue
		}

		for _, layer := range history {
			if layer.ID == "<missing>" || layer.ID == "" {
				continue
			}

			if _, ok := out[layer.ID]; !ok {
				// No tag info → this is an intermediate layer
				addImg(layer.ID, []string{}, []string{}, layer.Size, "")
			}
		}
	}

	// 3. Flatten into slice
	result := make([]*types.Image, 0, len(out))
	for _, img := range out {
		result = append(result, img)
	}

	return result
}
