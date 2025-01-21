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

	_, err = uc.productRepo.Getproduct(objectID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
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

	log.Println(id, " ", quantity)
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	var product models.Product
	objectID, errr := primitive.ObjectIDFromHex(id)
	if errr != nil {
		return false
	}
	dbproduct, err := uc.productRepo.Getproduct(objectID)
	if err != nil {
		return false
	}
	if dbproduct.Quantity < quantity {
		return false
	}
	dbproduct.Quantity -= quantity
	err = uc.productRepo.Updateproduct(objectID, dbproduct)
	if err != nil {
		return false
	}
	_, err = json.Marshal(product)
	return err == nil
}
