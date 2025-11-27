package images

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type ImagesLoadedMsg struct {
	Images []ImageItem
	Err    error
}

type DeleteImageMsg struct {
	ID  string
	Err error
}

func (m Model) FetchImagesCmd() tea.Cmd {
	return func() tea.Msg {
		images := m.cache.Images()
		imagesItems := make([]ImageItem, len(images))
		for i, img := range images {
			imagesItems[i] = ImageItem(img)
		}
		return ImagesLoadedMsg{Images: imagesItems, Err: nil}
	}
}

func (m Model) DeleteImagesCmd(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.cli.DeleteImage(context.Background(), id)
		if err != nil {
			return DeleteImageMsg{ID: id, Err: err}
		}
		return DeleteImageMsg{ID: id, Err: nil}
	}
}
