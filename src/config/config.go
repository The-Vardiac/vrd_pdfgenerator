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
	rabbitmqJob.DeclareExchange(RabbitmqChPubl, RMQMainExchange, "direct")
	// declare queue and bind #1 (vrdpdfgenerator)
	rabbitmqJob.DeclareQueue(RabbitmqChPubl, RMQPdfGeneratorQueue)
	rabbitmqJob.BindQueue(RabbitmqChPubl, jobs.Queue.Name, RMQPdfGeneratorQueueKey, RMQMainExchange)
	// declare queue and bind #2 (vrdmailer)
	rabbitmqJob.DeclareQueue(RabbitmqChPubl, RMQMailerQueue)
	rabbitmqJob.BindQueue(RabbitmqChPubl, jobs.Queue.Name, RMQMailerQueueKey, RMQMainExchange)
}

func (cfg *Config) InitAWSS3() {
	var awsS3Conf AmazonS3Conf

	awsS3Conf.Configure()
}

func (cfg *Config) InitMongoDB() {
	var mongoDB MongoDB

	mongoDB.MongoDBMakeConn()
}
