// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"tdl/config"
	"tdl/internal/handler"
	"tdl/internal/repository/RabbitMQ"
	"tdl/internal/repository/es"
	"tdl/internal/repository/mongodb"
	"tdl/internal/service"
	"tdl/internal/types"
	"tdl/pkg/elasticsearch"

	"github.com/gin-gonic/gin"
	r "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tdl/internal/repository/redis"
	"tdl/internal/repository/sql"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化MySQL
	dsn := cfg.GetMySQLDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&types.User{}, &types.Task{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化Redis
	rdb := r.NewClient(&r.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// 测试Redis连接
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化Elasticsearch
	esClient, err := elasticsearch.NewSearchService(cfg.ES.URL)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// 初始化MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	mongoDB := mongoClient.Database(cfg.MongoDB.Database)

	taskRepo := sql.NewTaskRepository(db)
	userRepo := sql.NewUserRepository(db)
	taskCache := redis.NewTaskCache(rdb)
	rabbitURL := fmt.Sprintf("amqp://%s:%s@localhost:5672/", cfg.RabbitMq.UserName, cfg.RabbitMq.Password)
	// 创建连接
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("RabbitMQ connection established")

	// 初始化生产者
	exchange := "remainder" // 交换机名称
	consumer, err := RabbitMQ.NewReminderConsumer(conn, "task_reminders", mongoDB.Client())
	if err != nil {
		log.Fatalf("Failed to create ReminderConsumer: %v", err)
	}
	rabbitmq, err := RabbitMQ.NewRabbitMQProducer(conn, exchange)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ: %v", err)
	}
	taskEsRepo := es.NewTaskRepository(esClient.Client)
	logRepo := mongodb.NewLogRepository(mongoDB)
	// 初始化服务
	taskService := service.NewTaskService(taskRepo, taskCache, taskEsRepo, logRepo, rabbitmq)
	userService := service.NewUserService(userRepo)

	// 初始化处理器
	taskHandler := handler.NewTaskHandler(taskService)
	userHandler := handler.NewUserHandler(userService)

	go consumer.Start()
	// 设置路由
	s := gin.Default()
	s.Use(CORSMiddleware())
	handler.RegisterRoutes(s, taskHandler, userHandler)
	// 启动服务器
	if err := s.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置 CORS 头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// 对于 OPTIONS 请求，直接返回 204
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
