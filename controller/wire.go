//go:build wireinject

package controller

import (
	"github.com/google/wire"
	"tdl/svc"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(User), "*"),
	wire.Struct(new(Task), "*"),
	svc.ProviderSet,
	wire.Struct(new(Controllers), "*"),
	NewGinServer,
)
