package docker

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

func mustNewTag(t *testing.T, s string) name.Tag {
	tag, err := name.NewTag(s, name.WeakValidation)
	if err != nil {
		t.Fatalf("NewTag(%v) = %v", s, err)
	}
	return tag
}

var fakeDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

func TestHead(t *testing.T) {

	mediaType := types.DockerManifestSchema1Signed

	tags := []name.Tag{mustNewTag(t, "mariadb:latest"), mustNewTag(t, "docker.getoutline.com/outlinewiki/outline:latest")}

	for _, tag := range tags {
		// Head should succeed even for invalid json. We don't parse the response.
		desc, err := remote.Head(tag)
		if err != nil {
			t.Fatalf("Head(%s) = %v", tag, err)
		}

		if desc.MediaType != mediaType {
			t.Errorf("Descriptor.MediaType = %q, expected %q", desc.MediaType, mediaType)
		}

		if desc.Digest.String() != fakeDigest {
			t.Errorf("Descriptor.Digest = %q, expected %q", desc.Digest, fakeDigest)
		}
	}
}
