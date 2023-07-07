package main

import (
	"context"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	args := os.Args
	argsLen := len(args)
	if argsLen < 2 {
		log.Fatal("Usage: go run main.go (up|down)")
	}
	if (args[1] != "up") && (args[1] != "down") {
		log.Fatal("Usage: go run main.go (up|down)")
	}
	migrateCommand := args[1]

	migrateCollections()

	driver, err := mongodb.WithInstance(config.MongoDBClient, &mongodb.Config{
		DatabaseName: "thevardiac",
	})
	if err != nil {
		log.Fatal("golang-migrate db driver: " + err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../db/migrations",
		"mongodb", driver)
	if err != nil {
		log.Fatal("golang-migrate: " + err.Error())
	}

	if migrateCommand == "up" {
		m.Up()
	}
	if migrateCommand == "down" {
		m.Down()
	}
}

func migrateCollections() {
	// aws_s3_bucket
	err := config.MongoTheVardiacDB.CreateCollection(context.TODO(), "aws_s3_bucket", nil)
	if err != nil {
		log.Println("Create collection: " + err.Error())
	}
}
