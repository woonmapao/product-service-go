package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/initializer"
	"github.com/woonmapao/product-service-go/models"
	"github.com/woonmapao/product-service-go/responses"
	"github.com/woonmapao/product-service-go/validations"
)

func CreateProduct(c *gin.Context) {
	// Extract product data from the request body
	var productData struct {
		Name          string  `json:"name"`
		Category      string  `json:"category"`
		Price         float64 `json:"price"`
		Description   string  `json:"description"`
		StockQuantity int     `json:"stockQuantity"`
		ReorderLevel  int     `json:"reorderLevel"`
	}

	err := c.ShouldBindJSON(&productData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
			}))
		return
	}

	// Check for duplicate product name
	if validations.IsProductNameDuplicate(productData.Name) {
		c.JSON(http.StatusConflict,
			responses.CreateErrorResponse([]string{
				"Product name is already taken",
			}))
		return
	}

	// Validate stock quantity
	if !validations.ValidateStockQuantity(c, productData.StockQuantity, productData.ReorderLevel) {
		c.JSON(http.StatusConflict,
			responses.CreateErrorResponse([]string{
				"Stock quantity must be greater than or equal to reorder level",
			}))
		return
	}

	// Create product in the database
	product := models.Product{
		Name:          productData.Name,
		Category:      productData.Category,
		Price:         productData.Price,
		Description:   productData.Description,
		StockQuantity: productData.StockQuantity,
		ReorderLevel:  productData.ReorderLevel,
	}

	err = initializer.DB.Create(&product).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to create product",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&product),
	)
}

// Retrieve all products from the database
func GetAllProducts(c *gin.Context) {

	var products []models.Product
	err := initializer.DB.Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch products",
			}))
		return
	}

	// Check if no products were found
	if len(products) == 0 {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"No products found",
			}))
		return
	}

	// Return a JSON response with the list of products
	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForMultipleProducts(products),
	)
}

func GetProductByID(c *gin.Context) {
	// Extract product ID from the request parameters
	productID := c.Param("id")

	// Convert product ID to integer (validations)
	id, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid product ID",
			}))
		return
	}

	// Get the product from the database
	var product models.Product
	err = initializer.DB.First(&product, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch product",
			}))
		return
	}

	// Check if the product was not found
	if product == (models.Product{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Product not found",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&product),
	)
}

func UpdateProduct(c *gin.Context) {
	// Extract product ID from the request parameters
	id := c.Param("id")

	// Extract updated product data from the request body
	var updatedProductData struct {
		Name          string  `json:"name" binding:"required"`
		Category      string  `json:"category" binding:"required"`
		Price         float64 `json:"price" binding:"required,gt=0"`
		Description   string  `json:"description"`
		StockQuantity int     `json:"stockQuantity" binding:"required,gte=0"`
		ReorderLevel  int     `json:"reorderLevel" binding:"required,gte=0,ltfield=StockQuantity"`
	}

	err := c.ShouldBindJSON(&updatedProductData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate the input data

	// Check if the product exists
	var existingProduct models.Product
	err = initializer.DB.First(&existingProduct, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	// Check if the updated product name is unique
	isDuplicate := validators.IsProductNameDuplicate(updatedProductData.Name)
	if isDuplicate {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product with this name already exists",
		})
		return
	}

	// Check if stock quantity is greater than or equal to reorder level
	if updatedProductData.StockQuantity < updatedProductData.ReorderLevel {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock quantity must be greater than or equal to reorder level",
		})
		return
	}

	// Update the product in the database
	initializer.DB.Model(&existingProduct).Updates(models.Product{
		Name:          updatedProductData.Name,
		Category:      updatedProductData.Category,
		Price:         updatedProductData.Price,
		Description:   updatedProductData.Description,
		StockQuantity: updatedProductData.StockQuantity,
		ReorderLevel:  updatedProductData.ReorderLevel,
	})

	// Return a JSON response with the updated product
	c.JSON(http.StatusOK, gin.H{
		"updatedProduct": existingProduct,
	})
}

func DeleteProduct(c *gin.Context) {
	// Extract product ID from the request parameters
	id := c.Param("id")

	// Delete the product from the database
	err := initializer.DB.Delete(&models.Product{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}

	// Return a JSON response indicating success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
