package models

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID       int    `json:"id"`
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Brand    string `json:"brand"`
}

type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) GetByID(id int) (*Product, error) {
	product := &Product{}

	query := `SELECT id, sku, name, brand
	          FROM products
	          WHERE id = $1`

	err := m.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Brand,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return product, nil
}