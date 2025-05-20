//go:build wireinject

package main

import (
	"github.com/google/wire"
	"tdl/controller"
	"tdl/dao"
	"tdl/pkg/client"

	"tdl/pkg/core"
)

func NewInjector() (*core.AppProvider, error) {
	panic(
		wire.Build(
			controller.ProviderSet,
			core.ProviderSet,
			client.ProviderSet,
			dao.ProviderSet,
		),
	)
}
