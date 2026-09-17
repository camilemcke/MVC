package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Analytics struct {
	TotalProducts     int `json:"total_products"`
	TotalStock        int `json:"total_stock"`
	LowStockProducts  int `json:"low_stock_products"`
	OutOfStockProducts int `json:"out_of_stock_products"`
	TotalEntries      int `json:"total_entries"`
	TotalExits        int `json:"total_exits"`
}

type BrandStock struct {
	Brand string `json:"brand"`
	Stock int    `json:"stock"`
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/analytics/brands", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT
				p.brand,
				COALESCE(SUM(i.stock), 0)
			FROM products p
			JOIN inventory i ON p.id = i.product_id
			GROUP BY p.brand
			ORDER BY SUM(i.stock) DESC
		`)

		if err != nil {
			http.Error(w, "could not calculate stock by brand", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		brands := []BrandStock{}

		for rows.Next() {
			brand := BrandStock{}

			err := rows.Scan(
				&brand.Brand,
				&brand.Stock,
			)

			if err != nil {
				http.Error(w, "could not read stock by brand", http.StatusInternalServerError)
				return
			}

			brands = append(brands, brand)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(brands)
	})

	http.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		analytics := Analytics{}

		err := db.QueryRow(`
			SELECT
				COUNT(*),
				COALESCE(SUM(stock), 0),
				COUNT(*) FILTER (WHERE stock BETWEEN 1 AND 5),
				COUNT(*) FILTER (WHERE stock = 0),
				COALESCE(SUM(entries), 0),
				COALESCE(SUM(exits), 0)
			FROM inventory
		`).Scan(
			&analytics.TotalProducts,
			&analytics.TotalStock,
			&analytics.LowStockProducts,
			&analytics.OutOfStockProducts,
			&analytics.TotalEntries,
			&analytics.TotalExits,
		)

		if err != nil {
			http.Error(w, "could not calculate analytics", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(analytics)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("analytics service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}