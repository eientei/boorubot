package pleroma

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path"
)

// MediaUpload returns media id for provided file path
func (client *Client) MediaUpload(ctx context.Context, name string, reader io.Reader) (string, error) {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", path.Base(name))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(fw, reader)
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	base := client.base

	base.Path += "/api/v1/media"

	resp := struct {
		ID    string `json:"id"`
		Error string `json:"error"`
	}{}

	err = client.exchange(ctx, http.MethodPost, base.String(), w.FormDataContentType(), &b, &resp)
	if err != nil {
		return "", err
	}

	if resp.Error != "" {
		return "", errors.New(resp.Error)
	}

	return resp.ID, nil
}
