package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	Description   string  `json:"description"`
	StockQuantity int     `json:"stock_quantity"`
	ReorderLevel  int     `json:"reorder_level"`
}

type ProductRequest struct {
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	Description   string  `json:"description"`
	StockQuantity int     `json:"stock_quantity"`
	ReorderLevel  int     `json:"reorder_level"`
}

type StockRequest struct {
	Quantity int `json:"quantity"`
}
