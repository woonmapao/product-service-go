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

func GetProductHandler(c *gin.Context) {

	id, err := controllers.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.GetError([]string{
				err.Error(),
			}))
		return
	}

	product, err := controllers.GetProduct(id, initializer.DB)
	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.GetError([]string{
				err.Error(),
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.GetSuccess(product),
	)
}

func UpdateProductHandler(c *gin.Context) {

	id, err := controllers.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	var updData models.ProductRequest
	err = controllers.BindAndValidate(c, &updData)
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
		return
	}

	// Find the updating user (validation)
	exist, err := controllers.GetProduct(id, tx)
	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Check if product valid
	valid, err := validations.IsValidProduct(updData.Name, updData.StockQuantity, updData.ReorderLevel, tx)
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

	err = controllers.UpdateProduct(&updData, exist, tx)
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
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.UpdateSuccess(),
	)
}

func DeleteProduct(c *gin.Context) {

	id, err := controllers.GetID(c)
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
		return
	}

	// Find the updating product (validation)
	_, err = controllers.GetProduct(id, tx)
	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	err = controllers.DeleteProduct(id, tx)
	if err != nil {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = controllers.CommitTrx(c, tx)
	if err != nil {
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.DeleteSuccess(),
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
