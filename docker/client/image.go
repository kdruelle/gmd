package client

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types/image"
)

// DeleteImage deletes an image from the Docker daemon.
// It does not force the deletion of the image, and it does prune children.
// The function returns an error if the deletion fails.
func (c *Client) DeleteImage(ctx context.Context, imageID string) error {
	_, err := c.cli.ImageRemove(ctx, imageID, image.RemoveOptions{
		Force:         false,
		PruneChildren: true,
	})
	return err
}

// PullImageWithProgress pulls an image from the Docker Hub and prints
// the progress of the pull to the given function.
// The function returns an error if the pull fails.
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

// ImageList returns a list of images on the Docker daemon.
// The function returns an error if the list of images cannot be retrieved.
// The list of images includes all images on the daemon, including intermediate images.
// The list of images is sorted by image name.
func (c *Client) ImageList() ([]image.Summary, error) {
	return c.cli.ImageList(context.Background(), image.ListOptions{All: true})
}

// ImageHistory returns the history of an image on the Docker daemon.
// The function returns a slice of image.HistoryResponseItem, where each item
// represents a layer in the image's history. The slice is sorted by
// creation time.
// The function returns an error if the history of the image cannot be
// retrieved.
func (c *Client) ImageHistory(imageID string) ([]image.HistoryResponseItem, error) {
	return c.cli.ImageHistory(context.Background(), imageID)
}
