//go:build integration

package registry

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetManifestDigestFromGoogle(t *testing.T) {
	ctx := context.Background()

	r, err := New(ctx, "europe-docker.pkg.dev", "", "", nil)
	require.NoError(t, err)
	digest, err := r.ManifestDigest(ctx, "/px-service-gcr/midgard/presence", "latest")
	require.NoError(t, err)
	require.Len(t, digest, 64)
}

func TestGetManifestDigestFromDockerHub(t *testing.T) {
	ctx := context.Background()

	r, err := New(ctx, "index.docker.io", os.Getenv("DOCKER_USERNAME"), os.Getenv("DOCKER_PASSWORD"), nil)
	require.NoError(t, err)
	digest, err := r.ManifestDigest(ctx, "/library/alpine", "latest")
	require.NoError(t, err)
	require.Len(t, digest, 64)
}
