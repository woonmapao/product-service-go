package validations

import (
	"github.com/gin-gonic/gin"
)

func IsProductNameDuplicate(yo string) bool {
	return false
}

// ValidateStockQuantity checks if stock quantity is greater than or equal to reorder level
func ValidateStockQuantity(c *gin.Context, stock, reorder int) bool {
	return stock >= reorder
}
