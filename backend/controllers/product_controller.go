package controllers

import (
	"encoding/json"
	"os"
	"net/http"
	"strconv"

	"backend/models"
)

type ProductController struct {
	ProductModel *models.ProductModel
}

func (c *ProductController) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	product, err := c.ProductModel.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Instance", os.Getenv("INSTANCE_NAME"))
	json.NewEncoder(w).Encode(product)
}