// internal/handler/routes.go
package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, taskHandler *TaskHandler, userHandler *UserHandler) {
	// 公共API
	public := r.Group("/api")
	{
		public.POST("/login", userHandler.Login)
		public.POST("/register", userHandler.Register)
	}

	// 需要认证的API
	authorized := r.Group("/api")
	{
		// 任务相关
		authorized.GET("/tasks", taskHandler.GetUserTasks)
		authorized.POST("/tasks", taskHandler.CreateTask)
		authorized.PUT("/tasks/:id", taskHandler.UpdateTask)
		authorized.DELETE("/tasks/:id", taskHandler.DeleteTask)
		authorized.GET("/tasks/search", taskHandler.SearchTasks)

		// 用户相关
		authorized.GET("/user/profile", userHandler.GetProfile)
		authorized.PUT("/user/profile", userHandler.UpdateProfile)
	}
}
