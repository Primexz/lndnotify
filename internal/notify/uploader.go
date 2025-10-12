package notify

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Primexz/lndnotify/pkg/file"
	log "github.com/sirupsen/logrus"
)

type Uploader struct {
	Url *url.URL
}

func (u *Uploader) Upload(message string, file *file.File) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	log.WithField("url", u.Url.Hostname()).Info("Uploading file")

	req, err := http.NewRequest("POST", u.Url.String(), bytes.NewReader(file.Content))
	if err != nil {
		return err
	}

	if user := u.Url.User; user != nil {
		if pass, ok := user.Password(); ok {
			req.SetBasicAuth(user.Username(), pass)
		}
	}

	if file.Filename != "" {
		req.Header.Set("Filename", file.Filename)
	}

	if message != "" {
		req.Header.Set("Message", message)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy error: status %d", resp.StatusCode)
	}
	return nil
}
