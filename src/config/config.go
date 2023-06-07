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

	// declare rabbitmq queue publisher
	rabbitmqJob.DeclareExchange(RabbitmqChPubl, "vardiac1", "direct")
	rabbitmqJob.DeclareQueue(RabbitmqChPubl, "pdfgeneratorqueue")
	rabbitmqJob.BindQueue(RabbitmqChPubl, jobs.Queue.Name, "pdfgeneratorqueuekey", "vardiac1")
}