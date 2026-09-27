package main

import (
	"context"
	"encoding/json"
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
	newOrder := OrderEvent{OrderID: 888, Status: "PAID"}
	body, err := json.Marshal(newOrder)
	if err != nil {
		log.Fatalf("failed to marshal the event message %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	forever := make(chan struct{})

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
	_, err = ch.QueueDeclare(
		"inventory_failures",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	msgs, err := ch.Consume(
		"inventory_failures",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume: %v", err)
	}
	go func() {
		for msg := range msgs {
			var failureEvent OrderEvent
			err := json.Unmarshal(msg.Body, &failureEvent)

			if err != nil {
				log.Fatalf("failed to unmarshal the event message %v", err)
			}

			log.Printf("COMPENSATING TRANSACTION: Marking order as CANCELLED! Refunding customer for order id : %d with reason :%s", failureEvent.OrderID, failureEvent.Reason)
		}
	}()

	log.Printf(" [x] Successfully shouted event into RabbitMQ: %s\n", body)
	<-forever
}

type OrderEvent struct {
	OrderID int    `json:"order_id"`
	Reason  string `json:"reason"`
	Status  string `json:"status"`
}
