package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"product/models"
	"product/repository"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userRepo *repository.UserRepository
}

func NewUserController(userRepo *repository.UserRepository) *UserController {
	return &UserController{userRepo}
}

// GetUser retrieves a user by ID
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := uc.userRepo.GetUser(objectID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

// GetUsers retrieves all users
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	users, err := uc.userRepo.GetUsers()
	if err != nil {
		http.Error(w, "Error finding users", http.StatusInternalServerError)
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

// CreateUser creates a new user
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user models.User

	// Decode the request body into the user struct
	json.NewDecoder(r.Body).Decode(&user)

	err := uc.userRepo.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(uj)
}

// UpdateUser updates a user's details
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "invalid user id for updation", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = uc.userRepo.UpdateUser(objectID, &user)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated user with ID: %s\n", id)
}

// DeleteUser deletes a user
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = uc.userRepo.DeleteUser(objectID)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user with ID: %s\n", id)
}

func (uc *UserController) Validateuser(id string) bool {
	// id := p.ByName("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// http.Error(w, "Invalid user ID", http.StatusBadRequest)
		// fmt.Fprintf(w, "false")
		return false
	}

	user, err := uc.userRepo.GetUser(objectID)
	if err != nil {
		// http.Error(w, "no user found with this id", http.StatusInternalServerError)
		// fmt.Fprintf(w, "false")
		return false
	}

	_, err = json.Marshal(user)
	if err != nil {
		// http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return false
	}
	return true
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusFound)
	// fmt.Fprintf(w, "true")
}
