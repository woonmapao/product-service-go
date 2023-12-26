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

func AddProduct(c *gin.Context) {

	// Extract product data from the request body
	var productData struct {
		Name          string  `json:"name" binding:"required"`
		Category      string  `json:"category" binding:"required"`
		Price         float64 `json:"price" binding:"required"`
		Description   string  `json:"description" binding:"required"`
		StockQuantity int     `json:"stockQuantity" binding:"required"`
		ReorderLevel  int     `json:"reorderLevel" binding:"required"`
	}
	err := c.ShouldBindJSON(&productData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
				err.Error(),
			}))
		return
	}

	// Check for empty values
	if productData.Name == "" || productData.Category == "" ||
		productData.Price == 0.0 || productData.Description == "" ||
		productData.StockQuantity == 0 || productData.ReorderLevel == 0 {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Name, Category, Price, Description, StockQuantity, and ReorderLevel are required fields",
			}))
		return
	}

	// Start a transaction
	tx := initializer.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to begin transaction",
				tx.Error.Error(),
			}))
		return
	}

	// Check for duplicate product name
	if validations.IsProductNameDuplicate(productData.Name, tx) {
		tx.Rollback()
		c.JSON(http.StatusConflict,
			responses.CreateErrorResponse([]string{
				"Product name is already taken",
			}))
		return
	}

	// Validate stock quantity
	if !validations.ValidateStockQuantity(productData.StockQuantity, productData.ReorderLevel, tx) {
		tx.Rollback()
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

	err = tx.Create(&product).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to create product",
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to commit transaction",
				err.Error(),
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

	// Extract updated product data from the request body
	var updateData struct {
		Name          string  `json:"name" binding:"required"`
		Category      string  `json:"category" binding:"required"`
		Price         float64 `json:"price" binding:"required,gt=0"`
		Description   string  `json:"description"`
		StockQuantity int     `json:"stockQuantity" binding:"required,gte=0"`
		ReorderLevel  int     `json:"reorderLevel" binding:"required,gte=0,ltfield=StockQuantity"`
	}

	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
			}))
		return
	}

	// Check if the product with the given ID exists
	var product models.Product
	err = initializer.DB.First(&product, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch product",
			}))
		return
	}

	// Check if the updated product name is unique
	if product == (models.Product{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Product not found:",
			}))
		return
	}

	// Validate stock quantity
	if !validations.ValidateStockQuantity(c, updateData.StockQuantity, updateData.ReorderLevel) {
		c.JSON(http.StatusConflict,
			responses.CreateErrorResponse([]string{
				"Stock quantity must be greater than or equal to reorder level",
			}))
		return
	}

	// Update product fields
	product.Name = updateData.Name
	product.Category = updateData.Category
	product.Price = updateData.Price
	product.Description = updateData.Description
	product.StockQuantity = updateData.StockQuantity
	product.ReorderLevel = updateData.ReorderLevel

	// Save the updated product to the database
	err = initializer.DB.Save(&product).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update product",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&product),
	)
}

func DeleteProduct(c *gin.Context) {
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

	// Check if the product with the given ID exists
	var product models.Product
	err = initializer.DB.First(&product, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch product",
			}))
		return
	}

	if product == (models.Product{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Product not found",
			}))
		return
	}

	// Delete the product
	err = initializer.DB.Delete(&models.Product{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to delete product",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.DeleteSuccessResponse(),
	)
}
