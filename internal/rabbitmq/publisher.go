package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishMessage(ch *amqp.Channel, queueName string, body string) error {
	_, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return err
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
		return err
	}

	log.Printf(" [x] Sent %s", body)
	return nil
}
