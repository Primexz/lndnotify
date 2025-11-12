package main

import (
	"flag"
	"runtime"

	"github.com/Primexz/lndnotify/internal/app"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version and Go version")
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat:  "2006/01/02 - 15:04:05",
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		SpacePadding:     45,
	})
	log.SetReportCaller(true)

	log.WithFields(log.Fields{
		"commit":     commit,
		"runtime":    runtime.Version(),
		"arch":       runtime.GOARCH,
		"build_date": date,
	}).Infof("⚡️ lndnotify %s", version)

	if *versionFlag {
		return
	}

	app.Run(*configPath, version)
}
