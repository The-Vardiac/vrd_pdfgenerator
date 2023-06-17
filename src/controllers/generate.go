package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	config "github.com/williamluisan/vrd_pdfgenerator/config"
	"github.com/williamluisan/vrd_pdfgenerator/repository"
	"github.com/williamluisan/vrd_pdfgenerator/services"
)

func Request(c *gin.Context) {
	var readFile services.Readfile
	readFile.Filename = "EdStatsData2.txt"
	file := readFile.ReadFile()
	var AWSS3PutObjectInput services.AWSS3PutObjectInput
	AWSS3PutObjectInput.Body = aws.ReadSeekCloser(strings.NewReader(file))
	AWSS3PutObjectInput.Bucket = aws.String("the-vardiac-bucket")
	AWSS3PutObjectInput.Key = aws.String(readFile.Filename)
	AWSS3PutObjectInput.PutObject()
	os.Exit(1)

	currentTime := time.Now()
	timeString := currentTime.Format("20060102150405.000")

	// send to rabbitmq
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	body := timeString
	err := config.RabbitmqChPubl.PublishWithContext(ctx, 
		"vardiac1",     // exchange
		"vrdpdfgeneratorqueuekey", // routing key
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
		"vrdpdfgeneratorqueue", // queue
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
		ctx, cancel := context.WithCancel(context.TODO())
		for d := range msgs {
			log.Println("Processing pdf ..")
			// readFile.Filename = "EdStatsData1.txt"
			readFile.Filename = "EdStatsData2.txt"

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

			// send email to vrd_mailer (via rest)
			// var sv_Vrd_mailer services.Vrd_mailer
			// sv_Vrd_mailer.Subject = "The Vardiac: your pdf document"
			// sv_Vrd_mailer.Body = "Filename " + generatePdf.Filename
			// sv_Vrd_mailer.MailTo = "lunba5th@gmail.com"
			// resp, err := sv_Vrd_mailer.Send()
			// if err != nil {
			// 	log.Println(resp + " | " + err.Error())
			// }

			// send to mailer queue
			vrdMailerData := repository.Vrd_mailer{
				Subject: "The Vardiac - Your PDF " + generatePdf.Filename,
				Body: generatePdf.Filename,
				MailTo: "lunba5th@gmail.com",
			}
			vrdMailerDataJson, _ := json.Marshal(vrdMailerData)
			if err != nil {
				log.Printf("%s: %s", "Failed to convert json", err)
			}
			body := string(vrdMailerDataJson)
			err = config.RabbitmqChPubl.PublishWithContext(ctx, 
				"vardiac1",     // exchange
				"vrdmailerqueuekey", // routing key
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

			counter++
		}
		cancel()
	}()
}