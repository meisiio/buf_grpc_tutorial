package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 2. Declare the Queue
	// We ensure the queue exists before we publish to it.
	q, err := ch.QueueDeclare(
		"order_events", // queue name
		true,           // durable (survives broker restarts)
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 3. Simulate a user checking out!
	log.Println("User clicked 'Checkout'...")
	// (Pretend we just saved Order 123 to our PostgreSQL database here)

	// 4. Publish the Event
	body := `{"order_id": 123, "status": "PAID"}`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange (default)
		q.Name, // routing key (the queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	log.Printf(" [x] Successfully shouted event into RabbitMQ: %s\n", body)
}
