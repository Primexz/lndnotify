package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Primexz/lndnotify/internal/config"
	"github.com/Primexz/lndnotify/internal/events"
	"github.com/Primexz/lndnotify/internal/lnd"
	"github.com/Primexz/lndnotify/internal/notify"
	log "github.com/sirupsen/logrus"
)

func Run(configPath string) {
	log.SetLevel(log.DebugLevel)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.WithError(err).Fatal("failed to load config file")
	}

	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("invalid log level")
	}
	log.SetLevel(level)

	// Create LND client
	lndClient := lnd.NewClient(cfg)

	// Connect to LND
	if err := lndClient.Connect(); err != nil {
		log.Fatalf("failed to connect to LND: %v", err)
	}
	defer func() {
		if err := lndClient.Disconnect(); err != nil {
			log.WithError(err).Error("error disconnecting from LND")
		}
	}()

	// Create notification manager
	notifier := notify.NewManager(&notify.ManagerConfig{
		Providers: cfg.Notifications.Providers,
		Templates: cfg.Notifications.Templates,
		Batching:  cfg.Notifications.Batching,
	})

	// Subscribe to events
	eventChan, err := lndClient.SubscribeEvents()
	if err != nil {
		log.Fatalf("failed to subscribe to events: %v", err)
	}

	if cfg.Events.StatusEvents {
		notifier.SendNotification("🟢 lndnotify connected")
		defer notifier.SendNotification("🔴 lndnotify disconnected")
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("started lndnotify")

	// Main event loop
	for {
		select {
		case event := <-eventChan:
			logger := log.WithField("event", event.Type())

			logger.Info("received event")

			if !event.ShouldProcess(cfg) {
				logger.Debug("event filtered, skipping")
				continue
			}

			msg, err := notifier.RenderTemplate(event.Type().String(), event.GetTemplateData(cfg.Notifications.Formatting.Locale.Tag))
			if err != nil {
				logger.WithError(err).Error("error rendering template")
				continue
			}

			if source, ok := event.(events.FileSource); ok {
				notifier.SendNotificationWithFile(msg, source.GetFile())
				continue
			}
			notifier.SendNotification(msg)

		case <-sigChan:
			log.Info("received shutdown signal")

			// Stop the notification manager first to flush any pending batches
			notifier.Stop()

			if err := lndClient.Disconnect(); err != nil {
				log.WithError(err).Error("error disconnecting from LND")
			}

			log.Info("shutdown complete")
			return
		}
	}
}
