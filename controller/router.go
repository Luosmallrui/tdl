package controller

import "github.com/gin-gonic/gin"

// Controller 接口定义所有 controller 必须实现的方法
type Controller interface {
	RegisterRouter(r gin.IRouter)
}

// RegisterRouters 注册所有路由
func (c *Controllers) RegisterRouters(r gin.IRouter) {
	c.User.RegisterRouter(r)
	c.Task.RegisterRouter(r)
	// ... 注册其他 controller 的路由
}

// Controllers 存储所有的 controller
type Controllers struct {
	User *User
	Task *Task
	// ... 添加其他 controller
}

func NewGinServer() *gin.Engine {
	return gin.Default()
}
