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
	userRepo := repository.NewUserRepository(client)
	userController := controllers.NewUserController(userRepo)
	productRepo := repository.NewProductRepository(client)
	productController := controllers.NewProductController(productRepo)

	// user routes
	router.GET("/users/:id", userController.GetUser)
	router.GET("/users", userController.GetUsers)
	router.POST("/users", userController.CreateUser)
	router.PUT("/users/:id", userController.UpdateUser)
	router.DELETE("/users/:id", userController.DeleteUser)

	//product routes
	router.GET("/products/:id", productController.GetProduct)
	router.GET("/products", productController.GetProducts)
	router.POST("/products", productController.CreateProduct)
	router.PUT("/products/:id", productController.UpdateProduct)
	router.DELETE("/products/:id", productController.DeleteProduct)

	//connecting rabbitmq
	go consumer.StartRabbitMQ()
	// Start the server
	log.Print("Starting server on :9000\n\n")
	log.Fatal(http.ListenAndServe(":9000", router))
}
