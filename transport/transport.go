package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type (
	HTTPClient interface {
		Get(ctx context.Context, url string) (*http.Response, error)
	}

	HTTPClientImpl struct {
		client *http.Client
	}
)

func NewHTTPClient() HTTPClient {
	return &HTTPClientImpl{
		client: http.DefaultClient,
	}
}

func (c HTTPClientImpl) Get(ctx context.Context, urlString string) (*http.Response, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	if parsedURL.Scheme == "" {
		return nil, errors.New("no scheme provided")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
