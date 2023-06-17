package repository

import "os"

type AWSS3PutObjectInput struct {
	Body   *os.File
	Bucket *string
	Key    *string
}