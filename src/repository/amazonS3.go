package repository

import "github.com/aws/aws-sdk-go/aws"

type AWSS3PutObjectInput struct {
	Body aws.ReaderSeekerCloser
	Bucket *string
	Key	*string
}