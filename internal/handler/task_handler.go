// internal/handler/task_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"tdl/internal/domain"
	"tdl/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task domain.CreateTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t := domain.Task{
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

	if err := h.taskService.CreateTask(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = uint(id)

	// 从JWT中获取用户ID
	userID, _ := c.Get("userID")
	task.UserID = userID.(uint)

	if err := h.taskService.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// 从JWT中获取用户ID
	//userID, _ := c.Get("userID")

	if err := h.taskService.DeleteTask(uint(id), 1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	// 从JWT中获取用户ID
	//userID, _ := c.Get("userID")

	tasks, err := h.taskService.GetUserTasks(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) SearchTasks(c *gin.Context) {
	query := c.Query("q")
	status := domain.TaskStatus(c.Query("status"))

	// 从JWT中获取用户ID
	//userID, _ := c.Get("userID")

	tasks, err := h.taskService.SearchTasks(query, status, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := make([]domain.ListTask, 0)

	for _, task := range tasks {
		res = append(res, domain.ListTask{
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
