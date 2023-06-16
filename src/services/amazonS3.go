package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/williamluisan/vrd_pdfgenerator/config"
	"github.com/williamluisan/vrd_pdfgenerator/repository"
)

type AWSS3PutObjectInput repository.AWSS3PutObjectInput

func (obj *AWSS3PutObjectInput) PutObject() (s3.PutObjectOutput, error) {
	var client = config.AwsS3Client

	input := &s3.PutObjectInput{
		Body:   obj.Body,
		Bucket: obj.Bucket,
		Key:    obj.Key,
	}

	result, err := client.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return s3.PutObjectOutput{}, nil
	}

	fmt.Println(result)

	return s3.PutObjectOutput{}, nil
}