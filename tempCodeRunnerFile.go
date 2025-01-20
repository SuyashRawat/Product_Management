package main

import (
	"log"
	"net/http"
	"time"

	"product/config"
	"product/consumer"
	"product/controllers"
	"product/repository"

	"github.com/julienschmidt/httprouter"
)

var userRepo *repository.UserRepository
var productRepo *repository.ProductRepository
var userController *controllers.UserController
var productController *controllers.ProductController

func main() {
	// Connect to MongoDB
	client, err := config.ConnectToDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	time.Sleep(100 * time.Millisecond)
	// Create a new router
	router := httprouter.New()

	// Initialize the controller
	userRepo = repository.NewUserRepository(client)
	userController = controllers.NewUserController(userRepo)
	productRepo = repository.NewProductRepository(client)
	productController = controllers.NewProductController(productRepo)

	// user routes
	router.GET("/users/:id", userController.GetUser)
	router.GET("/users", userController.GetUsers)
	router.POST("/users", userController.CreateUser)
	router.PUT("/users/:id", userController.UpdateUser)
	router.DELETE("/users/:id", userController.DeleteUser)
	router.GET("/validateuser/:id", userController.Validateuser)

	//product routes
	router.GET("/products/:id", productController.GetProduct)
	router.GET("/products", productController.GetProducts)
	router.POST("/products", productController.CreateProduct)
	router.PUT("/products/:id", productController.UpdateProduct)
	router.DELETE("/products/:id", productController.DeleteProduct)

	log.Println("Star :9000")
	conn, channel, err := consumer.MQConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer channel.Close()
	log.Println("connected to rabbit mq")
	// // Start consuming messages
	err = consumer.MQConsume(channel)
	if err != nil {
		log.Fatal(err)
	}
	// Start the server
	log.Println("Starting server on :9000")
	log.Fatal(http.ListenAndServe(":9000", router))

	//consumer code for rabbit mq
}
