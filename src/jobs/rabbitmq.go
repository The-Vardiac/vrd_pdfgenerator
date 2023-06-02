package jobs

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Queue	 amqp.Queue
)

type RabbitmqJob struct{}

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