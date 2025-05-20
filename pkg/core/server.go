package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tdl/controller"

	//"tdl/controller"
	"time"
)

var (
	once sync.Once
	// 服务唯一ID
	serverId string
)

func init() {
	once.Do(func() {
		id, err := gonanoid.Generate("0123456789abcdefghjklmnpqrstuvwxyz", 10)
		if err != nil {
			panic(err)
		}

		serverId = id
	})
}

func (app *AppProvider) RegisterRoutes() {
	app.Controllers.RegisterRouters(app.Engine)
}

type AppProvider struct {
	Engine      *gin.Engine
	Controllers *controller.Controllers
}

func GetServerId() string {
	return serverId
}

func Run(ctx context.Context, app *AppProvider) error {

	eg, groupCtx := errgroup.WithContext(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	//log.Printf("Server ID   :%s", server.ID())
	log.Printf("HTTP Listen Port :%d", 3002)
	log.Printf("HTTP Server Pid  :%d", os.Getpid())

	return run(c, eg, groupCtx, app)
}

func run(c chan os.Signal, eg *errgroup.Group, ctx context.Context, app *AppProvider) error {
	serv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 3002),
		Handler: app.Engine,
	}

	// 启动 http 服务
	eg.Go(func() error {
		err := serv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		defer func() {
			log.Println("Shutting down serv...")

			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer timeCancel()

			if err := serv.Shutdown(timeCtx); err != nil {
				log.Fatalf("HTTP Server Shutdown Err: %s", err)
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return nil
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("HTTP Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")

	return nil
}
