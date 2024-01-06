package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/initializer"
	"github.com/woonmapao/product-service-go/models"
	"github.com/woonmapao/product-service-go/responses"
	"gorm.io/gorm"
)

func BindAndValidate(c *gin.Context, body *models.ProductRequest) error {

	err := c.ShouldBindJSON(&body)
	if err != nil {
		return errors.New(
			"invalid request format",
		)
	}
	if body.Name == "" ||
		body.Category == "" ||
		body.Price == 0 ||
		body.Description == "" ||
		body.StockQuantity == 0 ||
		body.ReorderLevel == 0 {
		return errors.New(
			"missing fields",
		)
	}
	return nil
}

func StartTrx(c *gin.Context) (*gorm.DB, error) {

	tx := initializer.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
func AddProduct(product *models.ProductRequest, tx *gorm.DB) error {

	adding := models.Product{
		Name:          product.Name,
		Category:      product.Category,
		Price:         product.Price,
		Description:   product.Description,
		StockQuantity: product.StockQuantity,
		ReorderLevel:  product.ReorderLevel,
	}
	err := tx.Create(&adding).Error
	if err != nil {
		tx.Rollback()
		return errors.New("failed to add product")
	}
	return nil
}

func CommitTrx(c *gin.Context, tx *gorm.DB) error {

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to commit transaction",
				err.Error(),
			}))
		return err
	}
	return nil
}

func GetProducts(db *gorm.DB) (*[]models.Product, error) {

	var products []models.Product
	err := db.Find(&products).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return &products, errors.New("failed to fetch user")
	}
	if err == gorm.ErrRecordNotFound {
		return &products, errors.New("no products found")
	}
	return &products, nil
}
