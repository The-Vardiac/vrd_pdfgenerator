package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	config "github.com/williamluisan/vrd_pdfgenerator/config"
	routes "github.com/williamluisan/vrd_pdfgenerator/routes"
)

func init() {
	var config config.Config

	// initialize godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// initialize rabbitmq
	config.InitRabbitmq()
}

func main() {
	defer config.RabbitmqChPubl.Close()
	defer config.RabbitmqChCons.Close()

	// initialize gin
	router := gin.Default()
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}