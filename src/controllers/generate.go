package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	config "github.com/williamluisan/vrd_pdfgenerator/config"
	"github.com/williamluisan/vrd_pdfgenerator/services"
)

func Request(c *gin.Context) {
	currentTime := time.Now()
	timeString := currentTime.Format("20060102150405.000")

	// send to rabbitmq
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	body := timeString
	err := config.RabbitmqChPubl.PublishWithContext(ctx, 
		config.RMQMainExchange,     // exchange
		config.RMQPdfGeneratorQueueKey, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to publish a message", err)
	}

	// generate the pdfs
	var generateService services.Generate
	generateService.RMQConsumer()

	c.JSON(http.StatusOK, gin.H{
		"message": "done",
	})
}