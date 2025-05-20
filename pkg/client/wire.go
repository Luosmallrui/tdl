package client

import (
	"github.com/google/wire"
	"tdl/dao"
)

var ProviderSet = wire.NewSet(
	NewRedisClient,
	NewEsClient,
	NewMySQLClient,
	NewMongoDbClient,
	NewRabbitmqClient,
	dao.NewRabbitMQProducer,
)
