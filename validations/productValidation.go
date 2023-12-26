package validations

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IsProductNameDuplicate(prodName string, tx *gorm.DB) bool {
	return false
}

// ValidateStockQuantity checks if stock quantity is greater than or equal to reorder level
func ValidateStockQuantity(c *gin.Context, stock, reorder int) bool {
	return stock >= reorder
}
