package config

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitmqChPubl	*amqp.Channel
	RabbitmqChCons	*amqp.Channel

	RMQMainExchange = "vardiac1"

	RMQPdfGeneratorQueue = "vrdpdfgeneratorqueue"
	RMQPdfGeneratorQueueKey = "vrdpdfgeneratorqueuekey"

	RMQMailerQueue = "vrdmailerqueue"
	RMQMailerQueueKey = "vrdmailerqueuekey"
)

type RabbitmqConf struct{}

func (cfg *RabbitmqConf) RabbitmqMakeConn() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_CONN"))
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	chPubl, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a RabbitMQ publisher channel", err)
	}
	chCons, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a RabbitMQ consumer channel", err)
	}

	RabbitmqChPubl = chPubl
	RabbitmqChCons = chCons
}