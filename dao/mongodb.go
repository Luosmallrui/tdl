package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OperationLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    uint               `bson:"user_id"`
	Action    string             `bson:"action"`
	Target    string             `bson:"target"`
	TargetID  uint               `bson:"target_id"`
	Details   interface{}        `bson:"details,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}

type LogRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(dbClient *mongo.Client) *LogRepository {
	db := dbClient.Database("k")
	collection := db.Collection("operation_logs")

	// 检查是否已经存在TTL索引
	ctx := context.Background()
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		panic(err)
	}

	var indices []bson.M
	if err = cursor.All(ctx, &indices); err != nil {
		panic(err)
	}

	// 查找已存在的TTL索引
	ttlIndexExists := false
	for _, index := range indices {
		if name, ok := index["name"].(string); ok && name == "created_at_1" {
			ttlIndexExists = true
			break
		}
	}

	// 如果TTL索引不存在，创建新索引
	if !ttlIndexExists {
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{"created_at", 1}},
			Options: options.Index().SetExpireAfterSeconds(7 * 24 * 60 * 60),
		}

		_, err := collection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			panic(err)
		}
	}

	return &LogRepository{collection: collection}
}

func (r *LogRepository) AddLog(log *OperationLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	_, err := r.collection.InsertOne(context.Background(), log)
	return err
}

func (r *LogRepository) GetUserLogs(userID uint, limit int64) ([]OperationLog, error) {
	findOptions := options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(limit)

	cursor, err := r.collection.Find(
		context.Background(),
		bson.M{"user_id": userID},
		findOptions,
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var logs []OperationLog
	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
