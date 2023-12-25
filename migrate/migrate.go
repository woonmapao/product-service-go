package main

import (
	"log"

	"github.com/woonmapao/product-service-go/initializer"
	"github.com/woonmapao/product-service-go/models"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {

	err := initializer.DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to perform auto migration: &v", err)
	}
}
