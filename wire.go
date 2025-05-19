//go:build wireinject

package main

import (
	"github.com/google/wire"
	"tdl/controller"

	"tdl/pkg/core"
)

func NewInjector() *core.AppProvider {
	panic(
		wire.Build(
			controller.ProviderSet,
			core.ProviderSet,
		),
	)
}
