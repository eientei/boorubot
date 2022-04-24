// Package pleroma provides pleroma API integration
package pleroma

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/eientei/boorubot/integration/util/http/middleware"
)

// Config for pleroma client
type Config struct {
	HTTPClient *http.Client
	APIKey     string
	URL        string
}

// NewClient returns new pleroma client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{}
	}

	if config.APIKey != "" {
		config.HTTPClient.Transport = middleware.NewStaticHeadersMiddleware(
			config.HTTPClient.Transport, map[string]string{
				"Authorization": "Bearer " + config.APIKey,
			},
		)
	}

	base, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config: *config,
		base:   *base,
	}, nil
}

// Client instance
type Client struct {
	Config
	base url.URL
}

func (client *Client) exchange(
	ctx context.Context,
	method, url, ctype string,
	body io.Reader,
	data interface{},
) (err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return
	}

	if ctype != "" && body != nil {
		req.Header.Set("content-type", ctype)
	}

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if data != nil {
		err = json.NewDecoder(resp.Body).Decode(data)
	}

	return
}
