package main

type configuration struct {
	Broker      string `ini:"broker"`
	Username    string `ini:"username"`
	Password    string `ini:"password"`
	InsecureSSL bool   `ini:"insecure_ssl"`
	CACert      string `ini:"ca_cert"`
	NodeTopic   string `ini:"node_topic"`
	JobTopic    string `ini:"job_topic"`
}
