package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderEvent struct {
	TraceID string `json:"trace_id"`
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	//connect to the rabbitmq
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabitmq failed to open and opening encountered failure errors : %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel failed to open and opening encountered failure errors : %v", err)
	}
	defer ch.Close()

	failureQueue, err := ch.QueueDeclare(
		"inventory_failures",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare failure queue %v", err)
	}

	msgs, err := ch.Consume(
		"order_events",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("consume failed to open and opening encountered failure errors : %v", err)
	}
	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			var event OrderEvent

			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				log.Printf("[TraceID : %s ]failed to fetch the json :%v", event.TraceID, err)
				continue
			}

			log.Printf("[TraceID : %s ] Inventory System processing order %d with status %s", event.TraceID, event.OrderID, event.Status)

			log.Printf("[TraceID : %s ]Uh oh! Item for order %d is OUT OF STOCK!", event.TraceID, event.OrderID)

			failureBody := fmt.Sprintf(`{"order_id": %d, "reason": "out_of_stock"}`, event.OrderID)

			err = ch.PublishWithContext(
				context.Background(),
				"",
				failureQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(failureBody),
				})
			if err != nil {
				log.Printf("[TraceID : %s ]failed to publish a message : %v", event.TraceID, err)
			} else {
				log.Printf("[TraceID : %s ]Successfully published failure event for order %d", event.TraceID, event.OrderID)
			}
		}
	}()
	log.Printf(" [*] Inventory Service is waiting for orders. To exit press CTRL+C")
	<-forever
}
