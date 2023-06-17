package repository

import "os"

type AWSS3PutObjectInput struct {
	ACL	*string
	Body   *os.File
	Bucket *string
	Key    *string
}