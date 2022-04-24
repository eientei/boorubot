package pleroma

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

// StatusCreateRequest represents fediverse status (post) creation parameters
type StatusCreateRequest struct {
	Status      string   `json:"status,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	InReplyToID string   `json:"in_reply_to_id,omitempty"`
	MediaIDs    []string `json:"media_ids,omitempty"`
	Sensitive   bool     `json:"sensitive,omitempty"`
}

// Status model
type Status struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Error string `json:"error"`
}

var (
	symregex = regexp.MustCompile(`[^\pL\pN_]+`)
	numregex = regexp.MustCompile(`^\pN+$`)
)

// MakeTag returns corresponding tag for provided string
func MakeTag(s string) string {
	if numregex.MatchString(s) {
		s = "_" + s
	}

	return "#" + symregex.ReplaceAllString(s, "_")
}

// StatusCreate creates a new status
func (client *Client) StatusCreate(ctx context.Context, req *StatusCreateRequest) (status *Status, err error) {
	bs, err := json.Marshal(req)
	if err != nil {
		return
	}

	base := client.base

	base.Path += "/api/v1/statuses"

	status = &Status{}

	err = client.exchange(ctx, http.MethodPost, base.String(), "application/json", bytes.NewReader(bs), status)
	if err != nil {
		return nil, err
	}

	if status.Error != "" {
		return nil, errors.New(status.Error)
	}

	return status, nil
}
