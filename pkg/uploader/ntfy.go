package uploader

import (
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/push/ntfy"
)

func NewNtfyUploader(shoutrrrUrl *url.URL) (*HttpUploader, error) {
	service := ntfy.Service{}
	err := service.Initialize(shoutrrrUrl, nil)
	if err != nil {
		return nil, err
	}

	ntfyUrl, err := url.Parse(service.Config.GetAPIURL())
	if err != nil {
		return nil, err
	}
	return &HttpUploader{Url: ntfyUrl}, nil
}
