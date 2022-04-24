// Package danbooru providers danbooru API integration
package danbooru

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/eientei/boorubot/integration/util/http/middleware"
)

// Config for danbooru client
type Config struct {
	HTTPClient *http.Client
	APIKey     string
	Login      string
	URL        string
}

// NewClient returns new client instance
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{}
	}

	if config.APIKey != "" && config.Login != "" {
		bs := ([]byte)(config.Login + ":" + config.APIKey)

		config.HTTPClient.Transport = middleware.NewStaticHeadersMiddleware(
			config.HTTPClient.Transport, map[string]string{
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString(bs),
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
	method, url string,
	body io.Reader,
	data interface{},
) (err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return
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
