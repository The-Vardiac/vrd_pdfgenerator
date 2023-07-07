package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-pdf/fpdf"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/williamluisan/vrd_pdfgenerator/config"
	"github.com/williamluisan/vrd_pdfgenerator/repository"
	serviceMail "github.com/williamluisan/vrd_pdfgenerator/services/mail"
	"github.com/williamluisan/vrd_pdfgenerator/utils"
)

type Generate repository.Generate

func (generate *Generate) GeneratePDF() (extension string, err error) {
	extension = ".pdf"

	_, current_folder, _, _ := runtime.Caller(0)
	config_path := filepath.Dir(current_folder)
	absolute_path := config_path + "/../../files/pdfs/"

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)
	_, lineHt := pdf.GetFontSize()
	pdf.Write(lineHt, generate.Text)
	err = pdf.OutputFileAndClose(absolute_path + generate.Filename + extension)
	if err != nil {
		return extension, err
	}
	
	return extension, nil
}

func (generate *Generate) RMQConsumer() {
	msgs, err := config.RabbitmqChCons.Consume(
		config.RMQPdfGeneratorQueue, // queue
		"",                          // consumer
		false,                       // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to consume a message", err)
	}

	var readFile Readfile
	var generatePdf Generate
	
	go func() {
		var syncWaitGroup sync.WaitGroup

		counter := 0
		ctx, cancel := context.WithCancel(context.TODO())
		for d := range msgs {
			log.Println("Processing pdf ..")
			// readFile.Filename = "EdStatsData1.txt"
			readFile.Filename = "EdStatsData2.txt"

			// generate the pdf
			pdfGeneratedNameWithExtension := ""
			syncWaitGroup.Add(1)
			go func() {
				defer syncWaitGroup.Done()

				generatePdf.Filename = string(d.Body) + "_" + strconv.Itoa(counter)
				generatePdf.Text = readFile.ReadFile()
				ext, err := generatePdf.GeneratePDF()
				if err != nil {
					log.Println("PDF: failed to generate the pdf file " + generatePdf.Filename)

					// ...
					// retry to generate or
					// send email for fail generation
					// ...
				}
				pdfGeneratedNameWithExtension = generatePdf.Filename + ext
				log.Println("PDF: " + pdfGeneratedNameWithExtension + " done.")
			}()
			syncWaitGroup.Wait()

			syncWaitGroup.Add(1)
			go func() {
				defer syncWaitGroup.Done()

				// upload the pdf to S3 bucket
				log.Println("S3 Bucket: uploading ..")
				readFile.Filename = pdfGeneratedNameWithExtension
				file := readFile.GetPdfFileWithPath()
				fileToUpload, err := os.Open(file)
				if err != nil {
					log.Printf("S3 Bucket: Couldn't open file %v to upload: %v\n", pdfGeneratedNameWithExtension, err)
				}
				var AWSS3PutObjectInput AWSS3PutObjectInput
				AWSS3PutObjectInput.Body = fileToUpload
				AWSS3PutObjectInput.Bucket = aws.String(config.AwsS3MainBucket)
				AWSS3PutObjectInput.Key = aws.String(pdfGeneratedNameWithExtension)
				putObjectOutput, err := AWSS3PutObjectInput.PutObject()
				if err != nil {
					log.Println("Failed to upload to S3: " + err.Error())
				}
				log.Println("S3 Bucket: " + pdfGeneratedNameWithExtension + " uploaded.")

				// store bucket data to mongodb
				coll := config.MongoTheVardiacDB.Collection(config.COLL_AWS_S3_BUCKET)
				_, err = coll.InsertOne(context.TODO(), struct{
					Filename string
					Etag *string
					Versionid *string
				}{
					Filename: pdfGeneratedNameWithExtension,
					Etag: putObjectOutput.ETag,
					Versionid: putObjectOutput.VersionId,
				})
				if err != nil {
					log.Println("Failed to store to mongodb: " + err.Error())
				}
			}()
			syncWaitGroup.Wait()
			
			// send to mailer queue
			var AWSS3Object AWSS3Object
			var generateFileEmailTemplate serviceMail.GenerateFileEmailTemplate
			generateFileEmailTemplate.UserName = "Lunba5th" 
			generateFileEmailTemplate.Filename = pdfGeneratedNameWithExtension
			generateFileEmailTemplate.FileUrl = AWSS3Object.GetObjectUrl() + pdfGeneratedNameWithExtension
			generateFileEmailTemplate.GetGenerateFileHTMLEmailTemplate()
			emailTemplate := generateFileEmailTemplate.GetGenerateFileHTMLEmailTemplate()
			emailTemplateRendered := utils.HTMLTemplateParser(emailTemplate)
			vrdMailerData := repository.Vrd_mailer{
				Subject: "The Vardiac - Your PDF " + generatePdf.Filename,
				Body:    emailTemplateRendered,
				MailTo:  "lunba5th@gmail.com",
			}
			vrdMailerDataJson, _ := json.Marshal(vrdMailerData)
			if err != nil {
				log.Printf("%s: %s", "Failed to convert json", err)
			}
			body := string(vrdMailerDataJson)
			err = config.RabbitmqChPubl.PublishWithContext(ctx,
				config.RMQMainExchange,   // exchange
				config.RMQMailerQueueKey, // routing key
				false,                    // mandatory
				false,            // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				},
			)
			if err != nil {
				log.Panicf("%s: %s", "Failed to publish a message", err)
			}

			// acknowledge message
			if err = d.Ack(false); err != nil {
				log.Fatal("RabbitMQ: failed to acknowledge message in queue: " + string(d.Body))
			}

			counter++
		}
		cancel()
	}()
}