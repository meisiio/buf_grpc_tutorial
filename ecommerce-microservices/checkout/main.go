package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderEvent struct {
	OrderID int    `json:"order_id"`
	Reason  string `json:"reason"`
	Status  string `json:"status"`
}

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

	// 2. Declare the Queue for publishing
	q, err := ch.QueueDeclare(
		"order_events", // queue name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// 3. STARTUP LOGIC: Setup the Failure Consumer EXACTLY ONCE
	_, err = ch.QueueDeclare(
		"inventory_failures",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare failure queue: %v", err)
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

	// Run the consumer in the background forever
	go func() {
		for msg := range msgs {
			var failureEvent OrderEvent
			err := json.Unmarshal(msg.Body, &failureEvent)
			if err != nil {
				log.Printf("failed to unmarshal the event message %v", err)
				continue
			}
			log.Printf("COMPENSATING TRANSACTION: Marking order as CANCELLED! Refunding customer for order id: %d with reason: %s", failureEvent.OrderID, failureEvent.Reason)
		}
	}()

	// 4. PER-REQUEST LOGIC: The HTTP Endpoint
	http.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		// Generate a random dynamic Order ID
		fakeID := int(time.Now().Unix())
		fakeOrder := OrderEvent{
			OrderID: fakeID,
			Status:  "PAID",
		}

		body, err := json.Marshal(fakeOrder)
		if err != nil {
			log.Printf("failed to marshal the event message %v", err)
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(body),
			})
		if err != nil {
			log.Printf("Failed to publish a message: %v", err)
			http.Error(w, "Failed to publish order", http.StatusInternalServerError)
			return
		}

		log.Printf(" [x] Received HTTP request! Successfully shouted event into RabbitMQ: %s", body)
		w.Write([]byte("Order accepted into the Saga!\n"))
	})

	// 5. STARTUP LOGIC: Start the Web Server (This blocks forever, keeping the program alive)
	log.Printf("Checkout Service HTTP Server running on port 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
