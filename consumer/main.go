package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"product/config"
	"product/controllers"
	"product/repository"
	"time"

	"github.com/streadway/amqp"
)

// Define the message structure
type Order struct {
	ProductId string `json:"productID"`
	Quantity  int    `json:"quantity"`
	UserId    string `json:"userID"`
}

// Function to connect to RabbitMQ
func MQConnect() (*amqp.Connection, *amqp.Channel, error) {
	// Connect to RabbitMQ
	url := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	// Declare the same queue that the producer sends messages to
	_, err = channel.QueueDeclare(
		"email_queue", // Queue name
		true,          // Durable
		false,         // Delete when unused
		false,         // Exclusive
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		return nil, nil, err
	}

	return conn, channel, nil
}

// Function to consume messages from RabbitMQ
func MQConsume(channel *amqp.Channel) error {
	// Start consuming messages from the queue
	msgs, err := channel.Consume(
		"email_queue", // Queue name
		"",            // Consumer name (empty means random)
		true,          // Auto-acknowledge
		false,         // Exclusive
		false,         // No-local
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		return err
	}
	client, err := config.ConnectToDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Loop over messages and process them
	i := 0

	productRepo := repository.NewProductRepository(client)
	productController := controllers.NewProductController(productRepo)

	userRepo := repository.NewUserRepository(client)
	userController := controllers.NewUserController(userRepo)
	for msg := range msgs {
		var regData Order
		err := json.Unmarshal(msg.Body, &regData)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// Log received message
		fmt.Println("Received message:")
		fmt.Printf("product: %s ", regData.ProductId)
		fmt.Printf("Quantity: %d ", regData.Quantity)
		fmt.Printf("Userid: %s \n", regData.UserId)
		if !userController.Validateuser(regData.UserId) {
			log.Print("user updation not being done\n\n")
			continue
		}
		go func() {
			val := productController.ValidateProduct(regData.ProductId, regData.Quantity)
			if val {
				log.Print("value in db updated\n\n")
			} else {
				log.Print("no update was done for this query\n\n")
			}
		}()
		i++
		time.Sleep(1 * time.Second)
	}

	return nil
}
func StartRabbitMQ() {
	conn, channel, err := MQConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer channel.Close()
	log.Print("connected to rabbit mq\n\n")

	err = MQConsume(channel)
	if err != nil {
		log.Fatal(err)
	}
}
