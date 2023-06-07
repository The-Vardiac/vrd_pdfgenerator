package jobs

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Queue	 amqp.Queue
)

type RabbitmqJob struct{}

func (job *RabbitmqJob) DeclareExchange(ch *amqp.Channel, exchangeName string, exchangeType string) {
	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to declare an exchange", err)
	}
}

func (job *RabbitmqJob) DeclareQueue(ch *amqp.Channel, queueName string) {
	var err error

	Queue, err = ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to declare a queue", err)
	}
}

func (job *RabbitmqJob) BindQueue(ch *amqp.Channel, queueName string, routingKey string, exchangeName string) {
	err := ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to bind exchange and queue", err)
	}
}