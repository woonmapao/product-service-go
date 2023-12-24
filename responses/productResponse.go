package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/models"
)

func CreateSuccessResponse(product *models.Product) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Product added successfully",
		"data": gin.H{
			"product": gin.H{
				"id":            product.ID,
				"name":          product.Name,
				"category":      product.Category,
				"price":         product.Price,
				"description":   product.Description,
				"stockQuantity": product.StockQuantity,
				"reorderLevel":  product.ReorderLevel,
			},
		},
	}
}

func DeleteSuccessResponse() gin.H {
	return gin.H{
		"status":  "success",
		"message": "Product added successfully",
		"data":    gin.H{},
	}
}

// CreateErrorResponseForProduct formats the error response for product services
func CreateErrorResponse(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "Validation failed",
		"data": gin.H{
			"errors": errors,
		},
	}
}

func CreateSuccessResponseForMultipleProducts(products []models.Product) gin.H {
	productList := make([]map[string]interface{}, len(products))
	for _, product := range products {
		productList = append(productList, gin.H{
			"id":            product.ID,
			"name":          product.Name,
			"category":      product.Category,
			"price":         product.Price,
			"description":   product.Description,
			"stockQuantity": product.StockQuantity,
			"reorderLevel":  product.ReorderLevel,
		})
	}

	return gin.H{
		"status":  "success",
		"message": "Products retrieved successfully",
		"data": gin.H{
			"products": productList,
		},
	}
}
