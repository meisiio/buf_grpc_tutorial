package main

import (
	"context"
	"fmt"
	"log"

	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderEvent struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	//connect to the rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
				log.Printf("failed to fetch the json :%v", err)
				continue
			}

			log.Printf("Inventory System processing order %d with status %s", event.OrderID, event.Status)

			log.Printf("Uh oh! Item for order %d is OUT OF STOCK!", event.OrderID)

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
				log.Printf("failed to publish a message : %v", err)
			} else {
				log.Printf("Successfully published failure event for order %d", event.OrderID)
			}
		}
	}()
	log.Printf(" [*] Inventory Service is waiting for orders. To exit press CTRL+C")
	<-forever
}
