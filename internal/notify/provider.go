package notify

import (
	"github.com/Primexz/lndnotify/pkg/uploader"
	"github.com/nicholas-fedor/shoutrrr"
	log "github.com/sirupsen/logrus"
)

// initializeProviders sets up all notification providers
func (m *Manager) initializeProviders() {
	for _, p := range m.cfg.Providers {
		sender, err := shoutrrr.CreateSender(p.URL)
		if err != nil {
			log.WithField("provider", p.Name).WithError(err).Error("error creating sender")
			continue
		}

		name, url, err := sender.ExtractServiceName(p.URL)
		if err != nil {
			log.WithField("provider", p.Name).WithError(err).Error("cannot initialize uploader, invalid URL")
			m.providers[p.Name] = Provider{Sender: sender, Uploader: nil}
			continue
		}
		upl, err := uploader.NewUploader(name, url)
		if err != nil {
			log.WithField("provider", p.Name).WithError(err).Warn("error creating uploader")
			m.providers[p.Name] = Provider{Sender: sender, Uploader: nil}
			continue
		}
		m.providers[p.Name] = Provider{Sender: sender, Uploader: upl}
	}
}
