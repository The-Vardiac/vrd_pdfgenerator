package main

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
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

	driver, err := mongodb.WithInstance(config.MongoDBClient, &mongodb.Config{
		DatabaseName:         "thevardiac",
		MigrationsCollection: "x-migrations-collection",
		Locking:              "x-advisory-lock-collection",
	})
	if err != nil {
		log.Fatal("golang-migrate db driver: " + err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"../../db/migrations",
		"mongodb", driver)
	if err != nil {
		log.Fatal("golang-migrate: " + err.Error())
	}
	log.Fatal(m)

	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
}

func migrateCollections() {
	// aws_s3_bucket
	err := config.MongoTheVardiacDB.CreateCollection(context.TODO(), "aws_s3_bucket", nil)
	if err != nil {
		log.Println("Create colletion: " + err.Error())
	}
}
