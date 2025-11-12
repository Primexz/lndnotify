package notify

import (
	"fmt"
	"strings"
	"time"

	"github.com/Primexz/lndnotify/pkg/uploader"
	log "github.com/sirupsen/logrus"
)

// addToBatch adds a notification to the batch queue
func (m *Manager) addToBatch(message string, file *uploader.File) {
	m.batchMu.Lock()
	defer m.batchMu.Unlock()

	log.WithFields(log.Fields{
		"message":    message,
		"batch_size": len(m.batchQueue) + 1,
	}).Debug("adding notification to batch")
	m.batchQueue = append(m.batchQueue, QueuedNotification{
		Message: message,
		File:    file,
	})

	// Start or reset the flush timer
	if m.flushTimer != nil {
		m.flushTimer.Stop()
	}
	m.flushTimer = time.AfterFunc(m.cfg.Batching.FlushInterval, func() {
		m.flushBatch()
	})

	// Check if we've reached max batch size
	if len(m.batchQueue) >= m.cfg.Batching.MaxSize {
		if m.flushTimer != nil {
			m.flushTimer.Stop()
		}
		go m.flushBatch()
	}
}

// flushBatch sends all queued notifications
func (m *Manager) flushBatch() {
	m.batchMu.Lock()
	defer m.batchMu.Unlock()

	if len(m.batchQueue) == 0 {
		return
	}

	log.WithField("batch_size", len(m.batchQueue)).Debug("flushing notification batch")

	var messages []string
	var fileUploads []QueuedNotification

	for _, notification := range m.batchQueue {
		if notification.File != nil {
			fileUploads = append(fileUploads, notification)
		} else {
			messages = append(messages, notification.Message)
		}
	}

	// Send batched regular messages
	if len(messages) > 0 {
		m.sendBatch(messages)
	}

	// Send file uploads individually (they can't be batched)
	for _, upload := range fileUploads {
		m.uploadFile(upload.Message, upload.File)
	}

	m.batchQueue = m.batchQueue[:0]
}

// sendBatch sends multiple notifications as a batch with improved formatting
func (m *Manager) sendBatch(messages []string) {
	if len(messages) == 0 {
		return
	}

	count := len(messages)

	var batchMessage string
	if count > 1 {
		batchMessage = fmt.Sprintf("📢 %d Notifications\n", count)
		batchMessage += strings.Repeat("═", 10) + "\n"
	}

	for i, msg := range messages {
		if msg == "" {
			continue
		}

		if count > 1 {
			batchMessage += fmt.Sprintf("%d. %s", i+1, msg)
		} else {
			batchMessage += msg
		}

		if i < count-1 {
			batchMessage += "\n" + strings.Repeat("─", 5) + "\n"
		}
	}

	m.send(batchMessage)
}
