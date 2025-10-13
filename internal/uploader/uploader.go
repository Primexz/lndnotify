package uploader

import (
	"net/url"

	"github.com/nicholas-fedor/shoutrrr/pkg/services/push/ntfy"
	log "github.com/sirupsen/logrus"
)

func NewUploader(provider string, shoutrrrUrl *url.URL) Uploader {
	if shoutrrrUrl == nil {
		log.Error("No URL provided for uploader")
		return nil
	}

	switch provider {
	case "ntfy":
		uploader, err := newNtfyUploader(shoutrrrUrl)
		if err != nil {
			log.WithError(err).Error("error creating ntfy uploader")
			return nil
		}
		return uploader
	default:
		log.WithField("provider", provider).Warn("file upload not supported for this provider")
		return nil
	}
}

func newNtfyUploader(shoutrrrUrl *url.URL) (*httpUploader, error) {
	service := ntfy.Service{}
	err := service.Initialize(shoutrrrUrl, nil)
	if err != nil {
		return nil, err
	}

	ntfyUrl, err := url.Parse(service.Config.GetAPIURL())
	if err != nil {
		return nil, err
	}
	return &httpUploader{Url: ntfyUrl}, nil
}
