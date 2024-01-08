package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/controllers"
	"github.com/woonmapao/product-service-go/initializer"
	"github.com/woonmapao/product-service-go/models"
	"github.com/woonmapao/product-service-go/responses"
	"github.com/woonmapao/product-service-go/validations"
)

func AddProductHandler(c *gin.Context) {

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
	defer tx.Rollback()
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
	defer tx.Rollback()
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
	defer tx.Rollback()
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

func UpdateStockHandler(c *gin.Context) {

	id, err := controllers.GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Extract updated product data from the request body
	var body models.StockRequest
	err = controllers.BindAndValidate(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Start a transaction
	tx, err := controllers.StartTrx(c)
	defer tx.Rollback()
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

	enough, err := validations.EnoughStock(exist.StockQuantity, body.Quantity)
	if !enough {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	newStock := exist.StockQuantity - body.Quantity

	err = controllers.UpdateStock(newStock, exist, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest,
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
		responses.StockSuccess(),
	)

}
