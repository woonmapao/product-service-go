package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/controllers"
	"github.com/woonmapao/product-service-go/initializer"
	"github.com/woonmapao/product-service-go/models"
	"github.com/woonmapao/product-service-go/responses"
	"github.com/woonmapao/product-service-go/validations"
)

func AddProduct(c *gin.Context) {

	// Get data from request body
	var body models.ProductRequest
	err := controllers.BindAndValidate(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := controllers.StartTrx(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"there is a problem",
			}))
		return
	}

	// Check if product valid
	valid, err := validations.IsValidProduct(body.Name, body.StockQuantity, body.ReorderLevel, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"item not valid",
			}))
		return
	}

	// Add product to database
	err = controllers.AddProduct(&body, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = controllers.CommitTrx(c, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccess(),
	)
}

// Retrieve all products from the database
func GetProductsHandler(c *gin.Context) {

	products, err := controllers.GetProducts(initializer.DB)
	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.GetError([]string{
				err.Error(),
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.GetProductsSuccess(*products),
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
				err.Error(),
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
				err.Error(),
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
				err.Error(),
			}))
		return
	}

	// Extract updated product data from the request body
	var updateData struct {
		Name          string  `json:"name" binding:"required"`
		Category      string  `json:"category" binding:"required"`
		Price         float64 `json:"price" binding:"required,gt=0"`
		Description   string  `json:"description" binding:"required"`
		StockQuantity int     `json:"stockQuantity" binding:"required,gte=0"`
		ReorderLevel  int     `json:"reorderLevel" binding:"required,gte=0,ltfield=StockQuantity"`
	}
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
				err.Error(),
			}))
		return
	}

	// Check for empty values
	if updateData.Name == "" || updateData.Category == "" ||
		updateData.Price == 0.0 || updateData.Description == "" ||
		updateData.StockQuantity == 0 || updateData.ReorderLevel == 0 {
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

	// Check if the product with the given ID exists
	var product models.Product
	err = tx.First(&product, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch product",
				err.Error(),
			}))
		return
	}
	if product == (models.Product{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Product not found:",
			}))
		return
	}

	// Check for duplicate product name
	if validations.IsProductNameDuplicate(updateData.Name, tx) {
		c.JSON(http.StatusConflict,
			responses.CreateErrorResponse([]string{
				"Product name is already taken",
			}))
		return
	}

	// Validate stock quantity
	if !validations.ValidateStockQuantity(updateData.StockQuantity, updateData.ReorderLevel) {
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
	err = tx.Save(&product).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update product",
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
		responses.UpdateSuccessResponse(&product),
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
				err.Error(),
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

	// Check if the product with the given ID exists
	var product models.Product
	err = tx.First(&product, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch product",
				err.Error(),
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
	err = tx.Delete(&models.Product{}, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to delete product",
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
		responses.DeleteSuccessResponse(&product),
	)
}

func UpdateStock(c *gin.Context) {

	// Extract product ID from the request parameters
	productID := c.Param("id")

	// Convert product ID to integer (validations)
	id, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid product ID",
				err.Error(),
			}))
		return
	}

	// Extract updated product data from the request body
	var body struct {
		Quantity int `json:"quantity"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
				err.Error(),
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

	// Check if the product with the given ID exists
	var product models.Product
	err = tx.First(&product, id).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch product",
				err.Error(),
			}))
		return
	}
	if product == (models.Product{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Product not found:",
			}))
		return
	}

	newStock := product.StockQuantity - body.Quantity

	product.StockQuantity = newStock

	// Save the updated product to the database
	err = tx.Save(&product).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update product",
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
		responses.UpdateSuccessResponse(&product),
	)

}
