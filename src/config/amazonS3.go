package config

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	AwsS3MainBucket = "the-vardiac-bucket"

	AwsS3Session *session.Session
	AwsS3Client *s3.S3
)

type AmazonS3Conf struct{}

func (cfg *AmazonS3Conf) Configure() {
	aws_region := os.Getenv("AWS_REGION")
	aws_s3_access_key := os.Getenv("AWS_S3_ACCESS_KEY")
	aws_s3_secret_access_key := os.Getenv("AWS_S3_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(aws_region),
		Credentials: credentials.NewStaticCredentials(aws_s3_access_key, aws_s3_secret_access_key, ""),
	})
	if err != nil {
		log.Panicf("%s: %s", "Failed to configure an AWS S3 session: ", err.Error())
	}
	AwsS3Session = sess

	svc := s3.New(sess)
	AwsS3Client = svc
}