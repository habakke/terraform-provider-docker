package registry

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/habakke/terraform-provider-docker/internal/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
)

var GCRScopes = []string{"https://www.googleapis.com/auth/cloud-platform"}

type Registry interface {
	ManifestDigest(name string, tag string) (string, error)
}

type registry struct {
	dockerHost string
	Client     *http.Client
	Logger     util.Logger
	username   string
	password   string
}

func New(ctx context.Context, dockerHost, username, password string, logger util.Logger) (*registry, error) {
	var httpClient *http.Client
	googleCredentials, _ := google.FindDefaultCredentials(ctx, GCRScopes...)
	httpClient, err := createHTTPClient(ctx, googleCredentials, logger)
	if err != nil {
		return nil, err
	}

	return &registry{
		dockerHost: dockerHost,
		Client:     httpClient,
		Logger:     logger,
		username:   username,
		password:   password,
	}, nil
}

func (r registry) ManifestDigest(ctx context.Context, name string, tag string) (string, error) {
	b, err := r.dockerRegistryGet(ctx, fmt.Sprintf("/%s/manifests/%s", name, tag))
	if err != nil {
		return "", err
	}
	return dockerManifestDigest(b), nil
}

func createHTTPClient(ctx context.Context, googleCreds *google.Credentials, logger util.Logger) (*http.Client, error) {
	if googleCreds == nil {
		return &http.Client{
			Transport: util.NewLoggingRoundTripper(ctx, http.DefaultTransport, logger),
		}, nil
	}
	client := &http.Client{
		Transport: util.NewLoggingRoundTripper(ctx, &oauth2.Transport{
			Base:   http.DefaultTransport,
			Source: googleCreds.TokenSource,
		}, logger),
	}
	return client, nil
}

func dockerManifestDigest(manifest []byte) string {
	h := sha256.New()
	h.Write(manifest)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (r registry) dockerRegistryGet(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://"+r.dockerHost+"/v2"+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")

	if r.username != "" && r.password != "" {
		req.SetBasicAuth(r.username, r.password)
	}

	res, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("expected OK status, got %d from docker registry '%s'", res.StatusCode, path)
	}

	return io.ReadAll(res.Body)
}
