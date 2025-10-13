package events

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/internal/uploader"
	"github.com/lightningnetwork/lnd/lnrpc"
	"golang.org/x/text/language"
)

type BackupMultiEvent struct {
	Backup    *lnrpc.MultiChanBackup
	timestamp time.Time
}

type BackupMultiTemplate struct {
	ChanPoints    []string
	NumChanPoints int
	Filename      string
	Sha256Sum     string
}

func NewBackupMultiEvent(backup *lnrpc.MultiChanBackup) *BackupMultiEvent {
	return &BackupMultiEvent{
		Backup:    backup,
		timestamp: time.Now(),
	}
}

func (e *BackupMultiEvent) Type() EventType {
	return Event_BACKUP_MULTI
}

func (e *BackupMultiEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *BackupMultiEvent) GetTemplateData(lang language.Tag) interface{} {
	var chanPoints []string
	for _, cp := range e.Backup.ChanPoints {
		txHex := hex.EncodeToString(cp.GetFundingTxidBytes())
		chanPoint := fmt.Sprintf("%s:%d", txHex, cp.OutputIndex)
		chanPoints = append(chanPoints, chanPoint)
	}

	hash := sha256.Sum256(e.Backup.MultiChanBackup)
	sha256sum := hex.EncodeToString(hash[:])

	return &BackupMultiTemplate{
		NumChanPoints: len(e.Backup.ChanPoints),
		ChanPoints:    chanPoints,
		Filename:      e.getFileName(),
		Sha256Sum:     sha256sum,
	}
}

func (e *BackupMultiEvent) ShouldProcess(cfg *config.Config) bool {
	return cfg.Events.BackupEvents
}

func (e *BackupMultiEvent) getFileName() string {
	timestamp := e.timestamp.Format("20060102_150405")
	return "channel_backup_" + timestamp + ".backup"
}

func (e *BackupMultiEvent) GetFile() *uploader.File {
	return &uploader.File{
		Data:     e.Backup.MultiChanBackup,
		Filename: e.getFileName(),
	}
}
