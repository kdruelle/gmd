package cache

import (
	"log"
	"slices"
	"strings"

	"github.com/kdruelle/gmd/docker/client"
)

func (c *Cache) Images() []client.Image {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]client.Image, len(c.images))
	i := 0
	for _, img := range c.images {
		out[i] = *img
		i++
	}

	// Même tri que Portainer : les taggées d'abord
	slices.SortFunc(out, func(a, b client.Image) int {
		tagA := a.Tag()
		tagB := b.Tag()
		return strings.Compare(tagA, tagB)
	})

	return out
}

func (c *Cache) Image(id string) (client.Image, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c, ok := c.images[id]; ok {
		return *c, nil
	}
	return client.Image{}, ErrImageNotFound
}

func (c *Cache) ImagesUnused() []client.Image {
	c.mu.RLock()
	defer c.mu.RUnlock()

	used := make(map[string]bool, len(c.containers))
	log.Println("DEBUG === containers ===")
	for _, c := range c.containers {
		if c.Image != "" {
			log.Printf("container %s image %s", c.ID, c.Image)
			used[c.Image] = true
		}
	}

	out := make([]client.Image, 0, len(c.images))
	log.Println("DEBUG === images ===")
	for _, img := range c.images {
		log.Printf("img %s RepoDigests=%v", img.ID, img.RepoDigests)

		if _, ok := used[img.ID]; !ok {
			out = append(out, *img)
		}

	}
	// Même tri que Portainer : les taggées d'abord
	slices.SortFunc(out, func(a, b client.Image) int {
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
			c.images[id] = &client.Image{
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

func (c *Cache) snapshotImages() []*client.Image {
	list, err := c.cli.ImageList()
	if err != nil {
		panic(err)
	}

	out := make(map[string]*client.Image)

	addImg := func(id string, tags []string, digs []string, size int64, parent string) {
		if _, ok := out[id]; ok {
			return
		}
		out[id] = &client.Image{
			ID:          id,
			RepoTags:    tags,
			RepoDigests: digs,
			Size:        size,
			ParentID:    parent,
		}
	}

	// 1. Ajout des images "normales"
	for _, img := range list {
		addImg(img.ID, img.RepoTags, img.RepoDigests, img.Size, img.ParentID)
	}

	// 2. Ajout des parents via history
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
				// Pas d’infos de tag → c’est un layer intermédiaire
				addImg(layer.ID, []string{}, []string{}, layer.Size, "")
			}
		}
	}

	// 3. Flatten en slice
	result := make([]*client.Image, 0, len(out))
	for _, img := range out {
		result = append(result, img)
	}

	return result
}
