package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var configFile = flag.String("config", defaultConfigurationFile, "Configuration file")
	var help = flag.Bool("help", false, "Show help text")
	var version = flag.Bool("version", false, "Show version information")
	var _mode = flag.String("mode", "", "Set mode to node or job")
	var mode int
	var down = flag.Bool("down", false, "Mode down")
	var drain = flag.Bool("drain", false, "Mode drain")
	var idle = flag.Bool("idle", false, "Mode idle")
	var up = flag.Bool("up", false, "Mode up")
	var fini = flag.Bool("fini", false, "Mode fini")
	var _time = flag.Bool("time", false, "Mode time")
	var retain = flag.Bool("retain", false, "Retain MQTT message")

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

	if *_mode == "" {
		fmt.Fprint(os.Stderr, "Error: Missing --mode argument is mandatory\n\n")
		showHelp()
		os.Exit(1)
	}

	if strings.ToLower(*_mode) == "node" {
		mode = ModeNode
	} else if strings.ToLower(*_mode) == "job" {
		mode = ModeJob
	} else {
		fmt.Fprint(os.Stderr, "Error: Invalid value for mode\n\n")
		showHelp()
		os.Exit(1)
	}

	// validate mode and their flags
	if mode == ModeNode {
		if *fini {
			fmt.Fprintf(os.Stderr, "Error: Option --fini is not valid for node mode\n")
			os.Exit(1)
		}
		if *_time {
			fmt.Fprintf(os.Stderr, "Error: Option --time is not valid for node mode\n")
			os.Exit(1)
		}

		if !*down && !*drain && !*idle && !*up {
			fmt.Fprint(os.Stderr, "Error: Neither --down, --drain, --idle nor --up are set\n")
			os.Exit(1)
		}

		var _tmp int
		if *down {
			_tmp++
		}
		if *drain {
			_tmp++
		}
		if *idle {
			_tmp++
		}
		if *up {
			_tmp++
		}

		if _tmp > 1 {
			fmt.Fprint(os.Stderr, "Error: --down, --drain, --idle nor --up are mutually exclusive\n")
			os.Exit(1)
		}
	} else {
		if *down {
			fmt.Fprint(os.Stderr, "Error: Option --down is not valid for job mode\n")
			os.Exit(1)
		}
		if *drain {
			fmt.Fprint(os.Stderr, "Error: Option --drain is not valid for job mode\n")
			os.Exit(1)
		}
		if *idle {
			fmt.Fprint(os.Stderr, "Error: Option --idle is not valid for job mode\n")
			os.Exit(1)
		}
		if *up {
			fmt.Fprint(os.Stderr, "Error: Option --up is not valid for job mode\n")
			os.Exit(1)
		}

		if !*fini && !*_time {
			fmt.Fprint(os.Stderr, "Error: Neither --fini nor --time are set\n")
			os.Exit(1)
		}

		if *fini && *_time {
			fmt.Fprint(os.Stderr, "Error: --fini and --time are mutually exclusive\n")
			os.Exit(1)
		}
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

	if mode == ModeJob && cfg.JobTopic == "" {
		fmt.Fprint(os.Stderr, "Error: Job mode set but no job_topic defined\n")
		os.Exit(1)
	}

	if mode == ModeNode && cfg.NodeTopic == "" {
		fmt.Fprint(os.Stderr, "Error: Node mode set but no node_topic defined\n")
		os.Exit(1)
	}

	cfg.down = *down
	cfg.drain = *drain
	cfg.idle = *idle
	cfg.up = *up
	cfg.fini = *fini
	cfg.time = *_time
	cfg.mode = mode
	cfg.retain = *retain

	mqttClient, err := mqttConnect(cfg)
	if err != nil {
		log.WithFields(log.Fields{
			"broker":       cfg.Broker,
			"error":        err,
			"insecure_ssl": cfg.InsecureSSL,
			"password":     "***redacted***",
			"username":     cfg.Username,
			"qos":          cfg.QoS,
		}).Error("Can't connect to MQTT broker")
		os.Exit(1)
	}

	trailing := flag.Args()
	for _, arg := range trailing {
		err = mqttPublish(cfg, mqttClient, arg)
		if err != nil {
			log.WithFields(log.Fields{
				"broker":       cfg.Broker,
				"error":        err,
				"insecure_ssl": cfg.InsecureSSL,
				"password":     "***redacted***",
				"username":     cfg.Username,
				"qos":          cfg.QoS,
				"retain":       cfg.retain,
			}).Error("Can't publish data to MQTT broker")
			mqttClient.Disconnect(1000)
			os.Exit(1)
		}
	}
	mqttClient.Disconnect(1000)
}
