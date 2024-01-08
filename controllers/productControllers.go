package controllers

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/product-service-go/initializer"
	"github.com/woonmapao/product-service-go/models"
	"github.com/woonmapao/product-service-go/responses"
	"gorm.io/gorm"
)

func BindAndValidate[T any](c *gin.Context, body *T) error {

	err := c.ShouldBindJSON(&body)
	if err != nil {
		return errors.New(
			"invalid request format",
		)
	}
	if reflect.ValueOf(body).IsNil() {
		return errors.New("missing fields")
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

func GetID(c *gin.Context) (int, error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errors.New("invalid product id")
	}
	return id, nil
}

func GetProduct(id int, db *gorm.DB) (*models.Product, error) {

	var product models.Product
	err := db.First(&product, id).Error
	if err == gorm.ErrRecordNotFound {
		return &product, errors.New("product not found")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return &product, errors.New("something went wrong")
	}
	return &product, nil
}

func UpdateProduct(update *models.ProductRequest, exist *models.Product, tx *gorm.DB) error {

	exist = &models.Product{
		Name:          update.Name,
		Category:      update.Category,
		Price:         update.Price,
		Description:   update.Description,
		StockQuantity: update.StockQuantity,
		ReorderLevel:  update.ReorderLevel,
	}
	err := tx.Save(&exist).Error
	if err != nil {
		tx.Rollback()
		return errors.New("failed to update product")
	}
	return nil
}

func UpdateStock(newStock int, exist *models.Product, tx *gorm.DB) error {

	exist.StockQuantity = newStock

	err := tx.Save(&exist).Error
	if err != nil {
		tx.Rollback()
		return errors.New("failed to update stock")
	}
	return nil
}

func DeleteProduct(id int, tx *gorm.DB) error {

	err := tx.Delete(&models.Product{}, id).Error
	if err != nil {
		tx.Rollback()
		return errors.New("failed to delete product")
	}
	return nil
}
