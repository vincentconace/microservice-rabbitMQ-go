package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RAbbtiMQ connection.
	connectionRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectionRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectionRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subcribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance.
	app := fiber.New()

	// Add middleware to the app.
	app.Use(
		logger.New(), // Add simple logger.
	)

	// Add for send message to Service 1.
	app.Get("/send", func(c *fiber.Ctx) error {
		// Create a message to publish.
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("msg")),
		}

		// Attenpt to publish the message to the queue.
		if err := channelRabbitMQ.Publish(
			"",              // exchange
			"QueueService1", // queue container_name
			false,           // mandatory
			false,           // immediate
			message,         // message to publishing
		); err != nil {
			return err
		}

		return nil
	})

	// Start Fiber API server.
	log.Fatal(app.Listen(":3000"))
}
