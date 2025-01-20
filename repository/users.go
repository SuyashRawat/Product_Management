package repository

import (
	"context"
	"time"

	"product/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository interface defines the methods we need for DB operations
type UserRepository struct {
	client *mongo.Client
}

// NewUserRepository creates and returns a new UserRepository instance
func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{client}
}

// CreateUser inserts a new user into the database
func (ur *UserRepository) CreateUser(user *models.User) error {
	collection := ur.client.Database(db).Collection("Users")

	// Generate a new ObjectID if not already set
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, user)
	return err
}

// GetUser fetches a user by ID from the database
func (ur *UserRepository) GetUser(id primitive.ObjectID) (*models.User, error) {
	collection := ur.client.Database(db).Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers fetches all users from the database
func (ur *UserRepository) GetUsers() ([]models.User, error) {
	collection := ur.client.Database(db).Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates a user's details by ID
func (ur *UserRepository) UpdateUser(id primitive.ObjectID, user *models.User) error {
	collection := ur.client.Database(db).Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = id
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": user})
	return err
}

// DeleteUser deletes a user by ID
func (ur *UserRepository) DeleteUser(id primitive.ObjectID) error {
	collection := ur.client.Database(db).Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
