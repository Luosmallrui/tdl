package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"tdl/svc"
	"tdl/types"
)

type Task struct {
	TaskService svc.ITaskService
}

func (u *Task) RegisterRouter(r gin.IRouter) {
	task := r.Group("/")
	task.GET("/tasks", u.GetUserTasks)
	task.POST("/tasks", u.CreateTask)
	task.PUT("/tasks/:id", u.UpdateTask)
	task.DELETE("/tasks/:id", u.DeleteTask)
	task.GET("/tasks/search", u.SearchTasks)
}

func (u *Task) CreateTask(c *gin.Context) {

	var task types.CreateTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t := types.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserID:      task.UserID,
		DueDate:     task.DueDate,
		ReminderAt:  task.ReminderAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Tags:        strings.Join(task.Tags, ","),
	}

	// 从JWT中获取用户ID
	//userID, _ := c.Get("userID")
	t.UserID = 1
	if err := u.TaskService.CreateTask(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (u *Task) UpdateTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var task types.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = uint(id)

	// 从JWT中获取用户ID
	userID, _ := c.Get("userID")
	task.UserID = userID.(uint)

	if err := u.TaskService.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (u *Task) DeleteTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// 从JWT中获取用户ID
	//userID, _ := c.Get("userID")

	if err := u.TaskService.DeleteTask(uint(id), 1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (u *Task) GetUserTasks(c *gin.Context) {
	// 从JWT中获取用户ID
	//userID, _ := c.Get("userID")

	tasks, err := u.TaskService.GetUserTasks(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (u *Task) SearchTasks(c *gin.Context) {
	query := c.Query("q")
	status := types.TaskStatus(c.Query("status"))

	// 从JWT中获取用户ID
	//userID, _ := c.Get("userID")

	tasks, err := u.TaskService.SearchTasks(query, status, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := make([]types.ListTask, 0)

	for _, task := range tasks {
		res = append(res, types.ListTask{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			UserID:      task.UserID,
			DueDate:     task.DueDate,
			ReminderAt:  task.ReminderAt,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Tags:        strings.Split(task.Tags, ","),
		})
	}
	c.JSON(http.StatusOK, res)
}
