// cmd/server/main.go
package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"tdl/pkg/core"
)

func main() {

	ctx := &cli.Context{}
	App := NewInjector()
	App.RegisterRoutes()
	if err := core.Run(ctx, App); err != nil {
		fmt.Println(err)
	}
}
