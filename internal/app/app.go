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
		log.WithError(err).Fatal("Invalid log level")
	}
	log.SetLevel(level)

	// Create LND client
	lndClient := lnd.NewClient(cfg)

	// Connect to LND
	if err := lndClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to LND: %v", err)
	}
	defer lndClient.Disconnect()

	// Create notification manager
	notifier := notify.NewManager(&notify.ManagerConfig{
		Providers: cfg.Notifications.Providers,
		Templates: cfg.Notifications.Templates,
	})

	// Subscribe to events
	eventChan, err := lndClient.SubscribeEvents()
	if err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	if cfg.Events.StatusEvents {
		notifier.Send("🟢 lndnotify connected")
		defer notifier.Send("🔴 lndnotify disconnected")
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("started lndnotify")

	// Main event loop
	for {
		select {
		case event := <-eventChan:
			if !event.ShouldProcess(cfg) {
				log.WithField("event_type", event.Type()).Debug("event filtered, skipping")
				continue
			}

			msg, err := notifier.RenderTemplate(event.Type().String(), event.GetTemplateData(cfg.Notifications.Formatting.Locale.Tag))
			if err != nil {
				log.WithError(err).Error("error rendering template")
				continue
			}

			if source, ok := event.(events.FileSource); ok {
				notifier.UploadFile(msg, source.GetFile())
				continue
			}
			notifier.Send(msg)

		case <-sigChan:
			log.Info("received shutdown signal")
			if err := lndClient.Disconnect(); err != nil {
				log.WithError(err).Error("error disconnecting from LND")
			}

			log.Info("shutdown complete")
			return
		}
	}
}
