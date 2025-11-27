package client

import (
	"context"
	"encoding/json"

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

func (c *Client) DeleteImage(ctx context.Context, imageID string) error {
	_, err := c.cli.ImageRemove(ctx, imageID, image.RemoveOptions{
		Force:         false,
		PruneChildren: true,
	})
	return err
}

func (c *Client) PullImageWithProgress(ctx context.Context, imageRef string, progress func(map[string]interface{})) (err error) {
	reader, err := c.cli.ImagePull(ctx, imageRef, image.PullOptions{})
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

func (c *Client) ImageList() ([]image.Summary, error) {
	return c.cli.ImageList(context.Background(), image.ListOptions{All: true})
}

func (c *Client) ImageHistory(imageID string) ([]image.HistoryResponseItem, error) {
	return c.cli.ImageHistory(context.Background(), imageID)
}
