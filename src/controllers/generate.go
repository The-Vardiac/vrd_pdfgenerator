package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	config "github.com/williamluisan/vrd_pdfgenerator/config"
	jobs "github.com/williamluisan/vrd_pdfgenerator/jobs"
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
		"",     // exchange
		jobs.Queue.Name, // routing key
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
	consumer()

	c.JSON(http.StatusOK, gin.H{
		"message": "done",
	})
}

func consumer() {
	msgs, err := config.RabbitmqChCons.Consume(
		jobs.Queue.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to consume a message", err)
	}

	var readFile services.Readfile
	var generatePdf services.Generate

	go func() {
		counter := 0
		for d := range msgs {
			log.Println("Processing pdf ..")
			readFile.Filename = "EdStatsData1.txt"

			generatePdf.Filename = string(d.Body) + "_" + strconv.Itoa(counter)
			generatePdf.Text = readFile.ReadFile()
			err := generatePdf.GeneratePDF()
			if err != nil {
				log.Println("PDF: failed to generate the pdf file " + generatePdf.Filename)

				// ...
				// retry to generate or
				// send email for fail generation
				// ...
			}
			
			if err = d.Ack(false); err != nil {
				log.Fatal("RabbitMQ: failed to acknowledge message in queue: " + string(d.Body))
			}
			
			log.Println("PDF: " + generatePdf.Filename + " done.")

			counter++
		}
	}()
}