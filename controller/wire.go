//go:build wireinject

package controller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(User), "*"),
	wire.Struct(new(Task), "*"),
)
