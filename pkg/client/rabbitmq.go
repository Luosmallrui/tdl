package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func NewRabbitmqClient() *amqp.Connection {
	rabbitURL := fmt.Sprintf("amqp://%s:%s@localhost:5672/",
		"guest", "guest")
	// 创建连接
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("RabbitMQ connection established")
	return conn
}
