package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// 2. Tell RabbitMQ we want to consume messages from the "order_events" queue
	msgs, err := ch.Consume(
		"order_events", // queue name
		"",             // consumer name
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan struct{})
	go func(msgs <-chan amqp.Delivery) {
		for events := range msgs {
			log.Printf("Shipping Service preparing to ship Order: %s", string(events.Body))
		}
	}(msgs)

	log.Printf(" [*] Shipping Service is waiting for orders. To exit press CTRL+C")
	<-forever
}
