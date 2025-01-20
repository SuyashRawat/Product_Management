package repository

import (
	"context"
	"time"

	"product/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db string = "service2"
var c string = "product"

// ProductRepository interface defines the methods we need for DB operations
type ProductRepository struct {
	client *mongo.Client
}

// NewProductRepository creates and returns a new ProductRepository instance
func NewProductRepository(client *mongo.Client) *ProductRepository {
	return &ProductRepository{client}
}

// Createproduct inserts a new product into the database
func (ur *ProductRepository) CreateProduct(product *models.Product) error {
	collection := ur.client.Database(db).Collection(c)

	// Generate a new ObjectID if not already set
	if product.ID.IsZero() {
		product.ID = primitive.NewObjectID()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, product)
	return err
}

// Getproduct fetches a product by ID from the database
func (ur *ProductRepository) Getproduct(id primitive.ObjectID) (*models.Product, error) {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var product models.Product
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Getproducts fetches all products from the database
func (ur *ProductRepository) Getproducts() ([]models.Product, error) {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Updateproduct updates a product's details by ID
func (ur *ProductRepository) Updateproduct(id primitive.ObjectID, product *models.Product) error {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product.ID = id
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": product})
	return err
}

// Deleteproduct deletes a product by ID
func (ur *ProductRepository) Deleteproduct(id primitive.ObjectID) error {
	collection := ur.client.Database(db).Collection(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
