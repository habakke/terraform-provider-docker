package registry

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/habakke/terraform-provider-docker/internal/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"strings"
)

type LogfCallback func(ctx context.Context, format string, args ...interface{})

/*
 * Discard log messages silently.
 */
func Quiet(format string, args ...interface{}) {
	/* discard logs */
}

type Registry struct {
	URL    string
	Client *http.Client
	Logger util.Logger
}

/*
 * Create a new Registry with the given URL and credentials, then Ping()s it
 * before returning it to verify that the registry is available.
 *
 * You can, alternately, construct a Registry manually by populating the fields.
 * This passes http.DefaultTransport to WrapTransport when creating the
 * http.Client.
 */
func New(ctx context.Context, registryURL, username, password string, logger util.Logger) (*Registry, error) {
	transport := http.DefaultTransport

	return newFromTransport(ctx, registryURL, username, password, transport, logger)
}

/*
 * Create a new Registry, as with New, using an http.Transport that disables
 * SSL certificate verification.
 */
func NewInsecure(ctx context.Context, registryURL, username, password string) (*Registry, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// TODO: Why?
			InsecureSkipVerify: true, //nolint:gosec
		},
	}

	return newFromTransport(ctx, registryURL, username, password, transport, util.NewZerologLogger())
}

var GCRScopes = []string{"https://www.googleapis.com/auth/cloud-platform"}

func wrapOauth2Transport(ctx context.Context, transport http.RoundTripper) http.RoundTripper {
	creds, err := google.FindDefaultCredentials(ctx, GCRScopes...)
	if err != nil {
		return transport
	}
	return &oauth2.Transport{
		Base:   transport,
		Source: creds.TokenSource,
	}
}

func wrapBasicAuthTransport(username, password, url string, transport http.RoundTripper) http.RoundTripper {
	if username == "" {
		return transport
	}
	return &BasicTransport{
		Transport: transport,
		URL:       url,
		Username:  username,
		Password:  password,
	}
}

func wrapTokenTransport(username, password string, transport http.RoundTripper) http.RoundTripper {
	if username == "" {
		return transport
	}
	return &TokenTransport{
		Transport: transport,
		Username:  username,
		Password:  password,
	}
}

func wrapErrorTransport(transport http.RoundTripper) http.RoundTripper {
	return &ErrorTransport{
		Transport: transport,
	}
}

/*
 * Given an existing http.RoundTripper such as http.DefaultTransport, build the
 * transport stack necessary to authenticate to the Docker registry API. This
 * adds in support for OAuth bearer tokens and HTTP Basic auth, and sets up
 * error handling this library relies on.
 */
func WrapTransport(ctx context.Context, transport http.RoundTripper, url, username, password string) http.RoundTripper {
	return wrapErrorTransport(wrapBasicAuthTransport(username, password, url, wrapTokenTransport(username, password, wrapOauth2Transport(ctx, util.NewLoggingRoundTripper(ctx, transport, util.NewTerraformLogger())))))
}

func newFromTransport(ctx context.Context, registryURL, username, password string, transport http.RoundTripper, logger util.Logger) (*Registry, error) {
	url := strings.TrimSuffix(registryURL, "/")
	transport = WrapTransport(ctx, transport, url, username, password)
	registry := &Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
		Logger: logger,
	}

	if err := registry.Ping(); err != nil {
		return nil, err
	}

	return registry, nil
}

func NewFromClient(registryURL string, client *http.Client, logger util.Logger) (*Registry, error) {
	url := strings.TrimSuffix(registryURL, "/")
	registry := &Registry{
		URL:    url,
		Client: client,
		Logger: logger,
	}

	if err := registry.Ping(); err != nil {
		return nil, err
	}
	return registry, nil
}

func (r *Registry) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

func (r *Registry) Ping() error {
	url := r.url("/v2/")
	r.Logger.Infof(context.TODO(), "registry.ping url=%s", url)
	resp, err := r.Client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}
