package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var configFile = flag.String("config", defaultConfigurationFile, "Configuration file")
	var help = flag.Bool("help", false, "Show help text")
	var version = flag.Bool("version", false, "Show version information")

	flag.Usage = showHelp
	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	if *version {
		showVersion()
		os.Exit(0)
	}

	// Logging setup
	var logFmt = new(log.TextFormatter)
	logFmt.FullTimestamp = true
	logFmt.TimestampFormat = time.RFC3339
	log.SetFormatter(logFmt)

	cfg, err := loadConfiguration(*configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"config_file": *configFile,
			"error":       err,
		}).Error("Can't load configuration from configuration file")
		os.Exit(1)
	}

	fmt.Printf("%+v\n", cfg)

}
