package images

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/kdruelle/gmd/docker/types"
)

type ImageItem types.Image

func (i ImageItem) Title() string { return types.Image(i).Tag() }
func (i ImageItem) Description() string {
	return fmt.Sprintf("%s - %s", i.ID, humanize.Bytes(uint64(i.Size)))
}
func (i ImageItem) FilterValue() string { return i.Title() }
