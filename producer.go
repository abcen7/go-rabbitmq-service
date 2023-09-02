package main

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	utils "go-rabbitmq-service/utils"
	"log"
	"time"
)

type HubInteractionMessage struct {
	Id           uint32 `json:"id"`
	Certificate  string `json:"certificate"`
	NewState     string `json:"newState"`
	CurrentState string `json:"currentState"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Готовим синтетические данные
	body := HubInteractionMessage{
		Id:           1,
		Certificate:  "6c0c486749cb36b195bb59a279fa0793ba454a0b6d38e54d21cccd22da29b383",
		NewState:     "OPEN",
		CurrentState: "CLOSE",
	}

	// Преобразываем структуру в JSON-строку
	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Ошибка при преобразовании в JSON: %v", err)
	}

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
