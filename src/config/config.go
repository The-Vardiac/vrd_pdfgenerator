package config

import (
	jobs "github.com/williamluisan/vrd_pdfgenerator/jobs"
)

type Config struct{}

func (cfg *Config) InitRabbitmq() {
	var rabbitmqJob jobs.RabbitmqJob

	// make connection
	var rabbitmqConf RabbitmqConf
	rabbitmqConf.RabbitmqMakeConn()

	// declare rabbitmq exchange
	rabbitmqJob.DeclareExchange(RabbitmqChPubl, "vardiac1", "direct")
	// declare queue and bind #1 (vrdpdfgenerator)
	rabbitmqJob.DeclareQueue(RabbitmqChPubl, "vrdpdfgeneratorqueue")
	rabbitmqJob.BindQueue(RabbitmqChPubl, jobs.Queue.Name, "vrdpdfgeneratorqueuekey", "vardiac1")
	// declare queue and bind #2 (vrdmailer)
	rabbitmqJob.DeclareQueue(RabbitmqChPubl, "vrdmailerqueue")
	rabbitmqJob.BindQueue(RabbitmqChPubl, jobs.Queue.Name, "vrdmailerqueuekey", "vardiac1")
}