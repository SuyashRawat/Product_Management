package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"product/models"
	"product/repository"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductController struct {
	productRepo *repository.ProductRepository
}

func NewProductController(productRepo *repository.ProductRepository) *ProductController {
	return &ProductController{productRepo}
}

func (uc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := uc.productRepo.Getproduct(objectID)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

func (uc *ProductController) GetProducts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	users, err := uc.productRepo.Getproducts()
	if err != nil {
		http.Error(w, "Error finding products", http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

func (uc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var product models.Product

	json.NewDecoder(r.Body).Decode(&product)

	err := uc.productRepo.CreateProduct(&product)
	if err != nil {
		http.Error(w, "Error inserting product", http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(product)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(uj)
}

func (uc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "invalid product id for updation", http.StatusBadRequest)
		return
	}

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = uc.productRepo.Updateproduct(objectID, &product)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated product with ID: %s\n", id)
}

func (uc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = uc.productRepo.Deleteproduct(objectID)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted product with ID: %s\n", id)
}

func (uc *ProductController) ValidateProduct(id string, quantity int) bool {
	// id := p.ByName("id")
	log.Println(id, " ", quantity)
	// a := 1
	_, err := primitive.ObjectIDFromHex(id)
	// log.Println(a)
	if err != nil {
		return false
	}
	var product models.Product

	// json.NewDecoder(r.Body).Decode(&product)

	// log.Println(a)
	// quantity := product.Quantity
	objectID, errr := primitive.ObjectIDFromHex(id)
	if errr != nil {
		// http.Error(w, "product not found", http.StatusNotFound)
		// fmt.Fprintf(w, "false")
		return false
	}
	// log.Println(a)/
	dbproduct, err := uc.productRepo.Getproduct(objectID)
	if err != nil {
		// http.Error(w, "product not found", http.StatusNotFound)
		// fmt.Fprintf(w, "false")
		return false
	}
	// log.Println(a)
	if dbproduct.Quantity < quantity {
		// http.Error(w, "not enough items in stock", http.StatusInternalServerError)
		// fmt.Fprintf(w, "false")
		return false
	}
	// log.Println(a)
	dbproduct.Quantity -= quantity
	err = uc.productRepo.Updateproduct(objectID, dbproduct)
	if err != nil {
		// http.Error(w, "Error updating product", http.StatusInternalServerError)
		// fmt.Fprintf(w, "false")
		return false
	}

	// log.Println(a)
	_, err = json.Marshal(product)
	if err != nil {
		// http.Error(w, "Error encoding response", http.StatusInternalServerError)
		// fmt.Fprintf(w, "false")
		return false
	}
	// log.Println(a)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusFound)

	// fmt.Fprintf(w, "true")
	return true
}
