package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/models"
)

func CreateSuccess() gin.H {
	return gin.H{
		"status":  "success",
		"message": "product added successfully",
	}
}

func GetError(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "failed to fetch",
		"data": gin.H{
			"errors": errors,
		},
	}
}

func GetSuccess(product *models.Product) gin.H {
	return gin.H{
		"status":  "success",
		"message": "product retrieved successfully",
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

func UpdateSuccess() gin.H {
	return gin.H{
		"status":  "success",
		"message": "product updated successfully",
	}
}

func DeleteSuccessResponse(product *models.Product) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
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

func GetProductsSuccess(productList []models.Product) gin.H {

	products := make([]gin.H, len(productList))
	for i, product := range productList {
		products[i] = map[string]interface{}{
			"id":            product.ID,
			"name":          product.Name,
			"category":      product.Category,
			"price":         product.Price,
			"description":   product.Description,
			"stockQuantity": product.StockQuantity,
			"reorderLevel":  product.ReorderLevel,
		}
	}
	return gin.H{
		"status":  "success",
		"message": "products retrieved successfully",
		"data": gin.H{
			"products": products,
		},
	}
}
