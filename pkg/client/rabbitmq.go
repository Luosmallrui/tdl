package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"tdl/internal/repository/RabbitMQ"
)

func NewRabbitmqClient() *RabbitMQ.RabbitMQProducer {
	rabbitURL := fmt.Sprintf("amqp://%s:%s@localhost:5672/",
		"root", "123")
	// 创建连接
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	exchange := "remainder"
	rabbitmq, err := RabbitMQ.NewRabbitMQProducer(conn, exchange)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ: %v", err)
	}
	log.Println("RabbitMQ connection established")
	return rabbitmq
}
