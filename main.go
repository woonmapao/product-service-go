package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/initializer"
	"github.com/woonmapao/product-service-go/routes"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {
	r := gin.Default()

	// Setup routes
	routes.SetupProductRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)

}
