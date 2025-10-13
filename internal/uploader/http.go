package uploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type httpUploader struct {
	Url *url.URL
}

func (h *httpUploader) Upload(message string, file *File) error {
	if file == nil || len(file.Data) == 0 {
		return fmt.Errorf("no file data to upload")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadUrl := *h.Url
	if message != "" {
		query := uploadUrl.Query()
		query.Set("message", message)
		uploadUrl.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "POST", uploadUrl.String(), bytes.NewReader(file.Data))
	if err != nil {
		return err
	}

	if user := h.Url.User; user != nil {
		if pass, ok := user.Password(); ok {
			req.SetBasicAuth(user.Username(), pass)
		}
	}

	if file.Filename != "" {
		req.Header.Set("Filename", file.Filename)
	}

	contentType := http.DetectContentType(file.Data)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("http error: status %d", resp.StatusCode)
	}
	return nil
}
