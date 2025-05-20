//go:build wireinject

package dao

import (
	"github.com/google/wire"
	"tdl/dao/cache"
)

var ProviderSet = wire.NewSet(
	NewDbRepository,
	NewUserRepo,
	NewEsRepo,
	cache.NewTaskCache,
	NewLogRepository,
	NewReminderConsumer,
)
