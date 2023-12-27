package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/controllers"
)

func SetupProductRoutes(router *gin.Engine) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("/", controllers.AddProduct)
		productGroup.GET("/", controllers.GetAllProducts)
		productGroup.GET("/:id", controllers.GetProductByID)
		productGroup.PUT("/:id", controllers.UpdateProduct)
		productGroup.DELETE("/:id", controllers.DeleteProduct)
		productGroup.PUT("/update-stock/:id", controllers.UpdateStock)
	}
}
