package RabbitMQ

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQProducer 负责发送消息到 RabbitMQ
type RabbitMQProducer struct {
	channel   *amqp.Channel
	exchange  string
	delayType string // 延迟交换机类型
}

type ReminderConsumer struct {
	channel     *amqp.Channel
	queue       string
	mongoClient *mongo.Client
}

// 创建新的 Consumer
func NewReminderConsumer(conn *amqp.Connection, queue string, mongoClient *mongo.Client) (*ReminderConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &ReminderConsumer{channel: ch, queue: queue, mongoClient: mongoClient}, nil
}
func (c *ReminderConsumer) Start() error {
	// 声明队列
	_, err := c.channel.QueueDeclare(
		"task_reminders", // 队列名称
		true,             // durable: 队列持久化
		false,            // autoDelete: 消费者断开时是否删除队列
		false,            // exclusive: 是否独占队列
		false,            // noWait: 是否等待服务端响应
		nil,              // arguments: 其他参数
	)
	if err != nil {
		fmt.Println(err, 51)
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// 绑定队列到交换机
	err = c.channel.QueueBind(
		"task_reminders", // 队列名称
		"task_reminders", // routing key: 空字符串表示绑定所有消息（适用于 fanout 交换机）
		"remainder",      // 交换机名称
		false,            // noWait: 是否等待服务端响应
		nil,              // arguments: 其他参数
	)
	if err != nil {
		fmt.Println(err, 52)
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}
	msgs, err := c.channel.Consume(
		"task_reminders",
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// 处理消息
	for msg := range msgs {
		var task map[string]interface{}
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		// 触发提醒通知
		if err := c.handleReminder(task); err != nil {
			log.Printf("Failed to handle reminder: %v", err)
		}
	}

	return nil
}

func (c *ReminderConsumer) SendEmail() error {
	// QQ 邮箱 SMTP 配置
	smtpServer := "smtp.qq.com"    // QQ 邮箱的 SMTP 服务器地址
	port := "465"                  // 使用 SSL 加密的端口号
	username := "690722590@qq.com" // 发件人邮箱地址
	password := "cffikkgbpraebbjg" // QQ 邮箱授权码

	// 邮件内容
	to := "690722590@qq.com" // 收件人邮箱地址
	subject := "Subject: Task is Out of time\n"
	body := "Task is Out of time\n"

	// 邮件头部信息（必须符合 RFC 标准）
	headers := make(map[string]string)
	headers["From"] = `"Your Name" <690722590@qq.com>` // 发件人信息（昵称+邮箱地址）
	headers["To"] = to                                 // 收件人地址
	headers["Subject"] = subject                       // 邮件主题
	fmt.Println(headers, 55)

	// 构造完整的邮件内容
	var message strings.Builder
	for k, v := range headers {
		message.WriteString(k + ": " + v + "\r\n")
	}
	message.WriteString("\r\n" + body)

	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}

	// 建立连接并发送邮件
	conn, err := tls.Dial("tcp", smtpServer+":"+port, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// 进行身份验证
	auth := smtp.PlainAuth("", username, password, smtpServer)
	if err := client.Auth(auth); err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	// 设置发件人和收件人
	if err := client.Mail(username); err != nil {
		log.Fatalf("Failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		log.Fatalf("Failed to set recipient: %v", err)
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		log.Fatalf("Failed to start data command: %v", err)
	}
	_, err = writer.Write([]byte(message.String()))
	if err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}

	log.Println("Email sent successfully!")
	return nil
}

// 处理提醒通知
func (c *ReminderConsumer) handleReminder(task map[string]interface{}) error {
	// 模拟发送短信/邮件通知
	log.Printf("Sending notification for task: %v", task)
	c.SendEmail()

	// 记录到 MongoDB
	collection := c.mongoClient.Database("tdl").Collection("notifications")
	_, err := collection.InsertOne(nil, bson.M{
		"task_id":     task["task_id"],
		"user_id":     task["user_id"],
		"title":       task["title"],
		"reminder_at": task["reminder_at"],
		"notified_at": time.Now(),
	})
	return err
}

// 创建新的 Producer
func NewRabbitMQProducer(conn *amqp.Connection, exchange string) (*RabbitMQProducer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 声明延迟交换机
	err = ch.ExchangeDeclare(
		"remainder",
		"x-delayed-message",                    // 延迟交换机类型
		true,                                   // durable
		false,                                  // auto-delete
		false,                                  // internal
		false,                                  // no-wait
		amqp.Table{"x-delayed-type": "direct"}, // 延迟参数
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQProducer{channel: ch, exchange: exchange, delayType: "direct"}, nil
}

// 发布消息到延迟队列
func (p *RabbitMQProducer) Publish(queue string, message map[string]interface{}, reminderAt time.Time) error {
	// 计算延迟时间
	delay := int(reminderAt.Sub(time.Now()).Milliseconds())
	if delay < 0 {
		delay = 0
	}

	// 序列化消息
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	fmt.Println(string(body), 55, delay)
	// 发布消息
	return p.channel.Publish(
		"remainder",
		"task_reminders",
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     amqp.Table{"x-delay": int32(10000)}, // 指定延迟时间
		},
	)
}
