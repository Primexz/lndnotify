package events

import (
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/pkg/lndversion"
	"golang.org/x/text/language"
)

type LndUpdateAvailableEvent struct {
	LatestVersion  *lndversion.LndVersion
	CurrentVersion *lndversion.LndVersion
	timestamp      time.Time
}

type LndUpdateAvailableTemplate struct {
	LatestVersion  *lndversion.LndVersion
	CurrentVersion *lndversion.LndVersion
}

func NewLndUpdateAvailableEvent(latestVersion *lndversion.LndVersion, currentVersion *lndversion.LndVersion) *LndUpdateAvailableEvent {
	return &LndUpdateAvailableEvent{
		LatestVersion:  latestVersion,
		CurrentVersion: currentVersion,
		timestamp:      time.Now(),
	}
}

func (e *LndUpdateAvailableEvent) Type() EventType {
	return Event_LND_UPDATE_AVAILABLE
}

func (e *LndUpdateAvailableEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *LndUpdateAvailableEvent) GetTemplateData(lang language.Tag) interface{} {
	return &LndUpdateAvailableTemplate{
		LatestVersion:  e.LatestVersion,
		CurrentVersion: e.CurrentVersion,
	}
}

func (e *LndUpdateAvailableEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.LndUpdateEvents
}
