// Package middleware provides generic http client midleware
package middleware

import "net/http"

// StaticHeadersMiddleware adds fixed set of headers to every outgoing request
type StaticHeadersMiddleware struct {
	parent  http.RoundTripper
	headers map[string]string
}

// RoundTrip implementation
func (mw *StaticHeadersMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range mw.headers {
		req.Header.Set(k, v)
	}

	return mw.parent.RoundTrip(req)
}

// NewStaticHeadersMiddleware returns new StaticHeadersMiddleware instance, parent may be nil
func NewStaticHeadersMiddleware(parent http.RoundTripper, headers map[string]string) *StaticHeadersMiddleware {
	if parent == nil {
		parent = http.DefaultTransport
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	return &StaticHeadersMiddleware{
		parent:  parent,
		headers: headers,
	}
}
