package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		AvailableRoutes []string `json:"available_routes"`
	}{
		AvailableRoutes: []string{
			"list",
			"/latest/{product}",
			"/versions/{product}",
		},
	}
	json.NewEncoder(w).Encode(resp)
}

func handleListProducts(w http.ResponseWriter, r *http.Request) {
	products := index.ListProducts()
	json.NewEncoder(w).Encode(products)
}

func handleProductLatest(w http.ResponseWriter, r *http.Request) {
	product := chi.URLParam(r, "product")
	latest := index.LatestVersion(product)
	if latest == "" {
		http.Error(w, product+" not found", http.StatusNotFound)
		return
	}
	resp := struct {
		Product string `json:"product"`
		Latest  string `json:"latest_version"`
	}{
		Product: product,
		Latest:  latest,
	}
	json.NewEncoder(w).Encode(resp)
}

func handleListVersions(w http.ResponseWriter, r *http.Request) {
	product := chi.URLParam(r, "product")
	versions := index.ListVersions(product)
	if versions == nil {
		http.Error(w, product+" not found", http.StatusNotFound)
		return
	}
	resp := struct {
		Product  string   `json:"product"`
		Versions []string `json:"versions"`
	}{
		Product:  product,
		Versions: versions,
	}
	json.NewEncoder(w).Encode(resp)
}
