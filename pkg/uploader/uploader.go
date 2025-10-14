package uploader

import (
	"fmt"
	"net/url"
)

func NewUploader(provider string, shoutrrrUrl *url.URL) (Uploader, error) {
	if shoutrrrUrl == nil {
		return nil, fmt.Errorf("no URL provided for uploader")
	}

	switch provider {
	case "ntfy":
		uploader, err := NewNtfyUploader(shoutrrrUrl)
		if err != nil {
			return nil, fmt.Errorf("error creating ntfy uploader: %w", err)
		}
		return uploader, nil
	default:
		return nil, fmt.Errorf("file upload not supported for provider")
	}
}
