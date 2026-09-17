package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type Inventory struct {
	ProductID int `json:"product_id"`
	Entries   int `json:"entries"`
	Exits     int `json:"exits"`
	Stock     int `json:"stock"`
}

type InventoryItem struct {
	ProductID int    `json:"product_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Entries   int    `json:"entries"`
	Exits     int    `json:"exits"`
	Stock     int    `json:"stock"`
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
	http.HandleFunc("/inventory/all", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT
				p.id,
				p.sku,
				p.name,
				p.brand,
				i.entries,
				i.exits,
				i.stock
			FROM products p
			JOIN inventory i ON p.id = i.product_id
			ORDER BY p.id
		`)

		if err != nil {
			http.Error(w, "could not get inventory", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := []InventoryItem{}

		for rows.Next() {
			item := InventoryItem{}

			err := rows.Scan(
				&item.ProductID,
				&item.SKU,
				&item.Name,
				&item.Brand,
				&item.Entries,
				&item.Exits,
				&item.Stock,
			)

			if err != nil {
				http.Error(w, "could not read inventory", http.StatusInternalServerError)
				return
			}

			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	http.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")

		productID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		inventory := Inventory{}

		err = db.QueryRow(
			"SELECT product_id, entries, exits, stock FROM inventory WHERE product_id = $1",
			productID,
		).Scan(
			&inventory.ProductID,
			&inventory.Entries,
			&inventory.Exits,
			&inventory.Stock,
		)

		if err != nil {
			http.Error(w, "inventory not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(inventory)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("inventory service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}