package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Goreng", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "Silver Queen", Harga: 15000, Stok: 15},
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Product belum ada", http.StatusNotFound)
}

func updateProdukByID(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updatedProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			updatedProduk.ID = id
			produk[i] = updatedProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduk)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

func deleteProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode((map[string]string{
				"message": "sukses delete",
			}))
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProdukByID(w, r)
		} else if r.Method == "DELETE" {
			deleteProdukByID(w, r)
		}
	})

	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
			}
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
		}

	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API Running",
		})
	})

	fmt.Println("Starting server on :8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
