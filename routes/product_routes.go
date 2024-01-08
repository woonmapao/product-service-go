package routes

import (
	"github.com/gin-gonic/gin"
	h "github.com/woonmapao/product-service-go/handlers"
)

func SetupProductRoutes(router *gin.Engine) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("/", h.AddProductHandler)

		productGroup.GET("/", h.GetProductsHandler)
		productGroup.GET("/:id", h.GetProductHandler)

		productGroup.PUT("/:id", h.UpdateProductHandler)
		productGroup.PUT("/stock/:id", h.UpdateStockHandler)

		productGroup.DELETE("/:id", h.DeleteProduct)
	}
}
