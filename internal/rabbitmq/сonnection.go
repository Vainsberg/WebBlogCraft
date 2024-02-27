package rabbitmq

import (
	"log"

	config "github.com/Vainsberg/WebBlogCraft/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(cfg config.Сonfigurations) (*amqp.Connection, error) {

	conn, err := amqp.Dial(cfg.RABBITMQ_URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil, err
	}

	log.Println("Successfully connected to RabbitMQ")
	return conn, nil
}
