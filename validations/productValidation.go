package validations

import (
	"errors"

	"github.com/woonmapao/product-service-go/models"
	"gorm.io/gorm"
)

func IsValidProduct(name string, stock, reorder int, tx *gorm.DB) (bool, error) {

	var p models.Product
	err := tx.Where("name = ?", name).First(&p).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, errors.New("failed to validate")
	}
	// Return if stock is valid and no duplicate name
	return stock > reorder && err == gorm.ErrRecordNotFound, nil

}

func EnoughStock(want int, have int) (bool, error) {

	if want > have {
		return false, errors.New("not enough stock")
	}
	return true, nil
}
