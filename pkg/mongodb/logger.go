package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OperationLog 操作日志结构
type OperationLog struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	Action    string    `bson:"action"`
	Details   string    `bson:"details"`
	CreatedAt time.Time `bson:"created_at"`
}

// LogService 日志服务
type LogService struct {
	collection *mongo.Collection
}

// NewLogService 创建日志服务
func NewLogService(ctx context.Context, mongoURI string) (*LogService, error) {
	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 获取集合
	collection := client.Database("tdl").Collection("operation_logs")

	// 创建TTL索引，30天后过期
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(60 * 60 * 24 * 30),
	})
	if err != nil {
		return nil, err
	}

	return &LogService{collection: collection}, nil
}

// LogOperation 记录操作
func (s *LogService) LogOperation(ctx context.Context, userID, action, details string) error {
	log := OperationLog{
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, log)
	return err
}
