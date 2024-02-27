package rabbitmq

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Vainsberg/WebBlogCraft/internal/response"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
)

type RepositoryRabbitMQ struct {
	ch   *amqp.Channel
	conn *amqp.Connection
}

func NewRepositoryRabbitMQ(ch *amqp.Channel, conn *amqp.Connection) *RepositoryRabbitMQ {
	return &RepositoryRabbitMQ{ch: ch, conn: conn}
}

func (rab *RepositoryRabbitMQ) ConsumeMessages(queueName string) {
	var message response.RabbitMQMessage
	msgs, err := rab.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				return
			}
			code := strconv.Itoa(message.Code)

			m := gomail.NewMessage()
			m.SetHeader("From", "content.blog@mail.ru")
			m.SetHeader("To", message.Email)
			m.SetHeader("Subject", "Your code")
			m.SetBody("text/html", "Code <b>"+code+"</b>")
			d := gomail.NewDialer("smtp.mail.ru", 465, "content.blog@mail.ru", "1Y4f8WfSmxZzb0XPmFQ5")

			if err := d.DialAndSend(m); err != nil {
				panic(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (rab *RepositoryRabbitMQ) PublishMessage(queueName string, body []byte) error {
	_, err := rab.ch.QueueDeclare(
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

	err = rab.ch.Publish(
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
