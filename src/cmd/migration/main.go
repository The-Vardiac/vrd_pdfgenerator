package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	config "github.com/williamluisan/vrd_pdfgenerator/config"
)

func init() {
	// initialize godotenv
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config.Config
	cfg.InitMongoDB()
}

func main() {
	migrateCollections()
}

func migrateCollections() {
	// aws_s3_bucket
	err := config.MongoTheVardiacDB.CreateCollection(context.TODO(), "aws_s3_bucket", nil)
	if err != nil {
		log.Println("Create colletion: " + err.Error())
	}
}
