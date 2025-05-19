package svc

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(TaskService), "*"),
	wire.Bind(new(ITaskService), new(*TaskService)),
)
