package validations

import (
	"github.com/woonmapao/product-service-go/models"
	"gorm.io/gorm"
)

func IsProductNameDuplicate(name string, tx *gorm.DB) bool {
	var product models.Product
	err := tx.Where("name = ?", name).First(&product).Error
	return err == nil
}

func ValidateStockQuantity(stockQuantity, reorderLevel int, tx *gorm.DB) bool {
	return stockQuantity >= reorderLevel
}
