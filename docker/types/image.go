package types

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
