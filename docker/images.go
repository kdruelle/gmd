package docker

import (
	"context"
	"encoding/json"
	"log"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/image"
)

type Image struct {
	ID          string
	RepoTags    []string
	RepoDigests []string
	Size        int64
	ParentID    string
}

func (img Image) Tag() string {
	if len(img.RepoTags) > 0 {
		return img.RepoTags[0] // Portainer fait pareil : premier tag = dominant
	}
	if len(img.RepoDigests) > 0 {
		return img.RepoDigests[0]
	}
	return img.ID // fallback horrible mais nécessaire
}

func (m *Monitor) Images() []Image {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]Image, len(m.images))
	i := 0
	for _, img := range m.images {
		out[i] = *img
		i++
	}

	// Même tri que Portainer : les taggées d'abord
	slices.SortFunc(out, func(a, b Image) int {
		tagA := a.Tag()
		tagB := b.Tag()
		return strings.Compare(tagA, tagB)
	})

	return out
}

func (m *Monitor) ImagesUnused() []Image {
	m.mu.RLock()
	defer m.mu.RUnlock()

	used := make(map[string]bool, len(m.containers))
	log.Println("DEBUG === containers ===")
	for _, c := range m.containers {
		if c.Image != "" {
			log.Printf("container %s image %s", c.ID, c.Image)
			used[c.Image] = true
		}
	}

	out := make([]Image, 0, len(m.images))
	log.Println("DEBUG === images ===")
	for _, img := range m.images {
		log.Printf("img %s RepoDigests=%v", img.ID, img.RepoDigests)

		if _, ok := used[img.ID]; !ok {
			out = append(out, *img)
		}

	}
	// Même tri que Portainer : les taggées d'abord
	slices.SortFunc(out, func(a, b Image) int {
		tagA := a.Tag()
		tagB := b.Tag()
		return strings.Compare(tagA, tagB)
	})
	return out

}

func (m *Monitor) DeleteImage(ctx context.Context, imageID string) error {
	_, err := m.cli.ImageRemove(ctx, imageID, image.RemoveOptions{
		Force:         false,
		PruneChildren: true,
	})
	return err
}

func (dc *Monitor) PullImageWithProgress(ctx context.Context, imageRef string, progress func(map[string]interface{})) (err error) {
	reader, err := dc.cli.ImagePull(ctx, imageRef, image.PullOptions{})
	if err != nil {
		return err
	}
	defer func() {
		err = reader.Close()
	}()
	decoder := json.NewDecoder(reader)

	for decoder.More() {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
			return err
		}
		progress(msg)
	}

	return nil
}

func (m *Monitor) refreshImage(id string) {
	imgs, err := m.cli.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.images, id)

	// faut retrouver l’image par ID
	for _, img := range imgs {
		if img.ID == id {
			m.images[id] = &Image{
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

func (m *Monitor) snapshotImages() []*Image {
	ctx := context.Background()
	list, err := m.cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	out := make(map[string]*Image)

	addImg := func(id string, tags []string, digs []string, size int64, parent string) {
		if _, ok := out[id]; ok {
			return
		}
		out[id] = &Image{
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
		history, err := m.cli.ImageHistory(ctx, img.ID)
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
	result := make([]*Image, 0, len(out))
	for _, img := range out {
		result = append(result, img)
	}

	log.Printf("all images loaded")

	return result
}
