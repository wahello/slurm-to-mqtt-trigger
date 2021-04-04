package main

import (
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/ini.v1"
)

func loadConfiguration(file string) (configuration, error) {
	var result configuration

	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, file)
	if err != nil {
		return result, err
	}

	mqttCfgSection, err := cfg.GetSection("mqtt")
	if err != nil {
		return result, err
	}

	err = mqttCfgSection.MapTo(&result)
	if err != nil {
		return result, err
	}

	err = validateConfiguration(result)
	if err != nil {
		return result, err
	}

	_url, err := url.Parse(result.Broker)
	if err != nil {
		return result, err
	}

	result.Broker = fmt.Sprintf("%s://%s", _url.Scheme, _url.Host)
	result.NodeTopic = strings.TrimSuffix(strings.TrimPrefix(result.NodeTopic, "/"), "/")
	result.JobTopic = strings.TrimSuffix(strings.TrimPrefix(result.JobTopic, "/"), "/")

	if _url.Scheme == "ssl" || _url.Scheme == "tls" || _url.Scheme == "mqtts" {
		result.useTLS = true
	}

	result.qos = byte(result.QoS)

	return result, nil
}

func validateConfiguration(config configuration) error {
	if config.Broker == "" {
		return fmt.Errorf("Missing broker configuration")
	}

	_url, err := url.Parse(config.Broker)
	if err != nil {
		return err
	}

	if _url.Scheme != "tcp" && _url.Scheme != "ssl" && _url.Scheme != "tls" && _url.Scheme != "mqtt" && _url.Scheme != "mqtts" {
		return fmt.Errorf("Invalid scheme in broker URL")
	}

	if _url.Host == "" {
		return fmt.Errorf("Invalid or missing hostname in broker URL")
	}

	if config.NodeTopic == "" && config.JobTopic == "" {
		return fmt.Errorf("Neither node_topic nor job_topic are defined")
	}

	if config.QoS > 2 {
		return fmt.Errorf("Invalid QoS value")
	}

	return nil
}
