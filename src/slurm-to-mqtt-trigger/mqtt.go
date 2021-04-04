package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/nu7hatch/gouuid"
)

func mqttConnect(cfg configuration) (mqtt.Client, error) {
	var client mqtt.Client

	_uuid, err := uuid.NewV4()
	if err != nil {
		return client, err
	}

	// MQTT 3.1 standard limits the client ID to 23 bytes
	clientID := fmt.Sprintf("strigger-%s", _uuid.String()[0:13])

	copts := mqtt.NewClientOptions()
	copts.AddBroker(cfg.Broker)
	// AutoReconnect is rather pointless if QoS is 0
	if cfg.qos > 0 {
		copts.SetAutoReconnect(true)
	}
	copts.SetCleanSession(true)

	// Enforce MQTT 3.1.1
	copts.SetProtocolVersion(4)

	copts.SetUsername(cfg.Username)
	copts.SetPassword(cfg.Password)

	copts.SetClientID(clientID)

	if cfg.useTLS {
		tlsConfig := &tls.Config{}

		if cfg.CACert != "" {
			caPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(cfg.CACert)
			if err != nil {
				return client, err
			}

			ok := caPool.AppendCertsFromPEM(ca)
			if !ok {
				return client, fmt.Errorf("Unable to add CA certificate to CA pool")
			}

			tlsConfig.RootCAs = caPool
		}

		if cfg.InsecureSSL {
			tlsConfig.InsecureSkipVerify = true
		}

		copts.SetTLSConfig(tlsConfig)
	}

	client = mqtt.NewClient(copts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	return client, nil
}

func mqttPublish(cfg configuration, client mqtt.Client, nodeOrJob string) error {
	var nodeState string
	var jobState string

	if cfg.mode == ModeNode {
		if cfg.down {
			nodeState = "down"
		} else if cfg.drain {
			nodeState = "drained"
		} else if cfg.idle {
			nodeState = "idle"
		} else if cfg.up {
			nodeState = "up"
		}
		nodeStateTopic := cfg.NodeTopic + "/state/" + nodeState
		nodeNodeTopic := cfg.NodeTopic + "/node/" + nodeOrJob

		token := client.Publish(nodeStateTopic, cfg.qos, cfg.retain, nodeOrJob)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}

		token = client.Publish(nodeNodeTopic, cfg.qos, cfg.retain, nodeState)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	} else if cfg.mode == ModeJob {
		if cfg.fini {
			jobState = "finished"
		} else if cfg.time {
			jobState = "timelimit"
		}
		jobStateTopic := cfg.JobTopic + "/state/" + jobState
		jobJobTopic := cfg.JobTopic + "/job/" + nodeOrJob

		token := client.Publish(jobStateTopic, cfg.qos, cfg.retain, nodeOrJob)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}

		token = client.Publish(jobJobTopic, cfg.qos, cfg.retain, jobState)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}
