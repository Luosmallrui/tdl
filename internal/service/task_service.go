// Package service internal/service/task_service.go
package service

import (
	"fmt"
	"tdl/internal/domain"
	"tdl/internal/repository/RabbitMQ"
	"tdl/internal/repository/es"
	"tdl/internal/repository/mongodb"
	"tdl/internal/repository/redis"
	"tdl/internal/repository/sql"
	"time"
)

type TaskService struct {
	taskRepo         *sql.TaskRepository
	taskCache        *redis.TaskCache
	taskEsRepo       *es.TaskRepository
	logRepo          *mongodb.LogRepository
	reminderProducer *RabbitMQ.RabbitMQProducer
}

func NewTaskService(
	taskRepo *sql.TaskRepository,
	taskCache *redis.TaskCache,
	taskEsRepo *es.TaskRepository,
	logRepo *mongodb.LogRepository,
	reminderProducer *RabbitMQ.RabbitMQProducer,
) *TaskService {
	return &TaskService{
		taskRepo:         taskRepo,
		taskCache:        taskCache,
		taskEsRepo:       taskEsRepo,
		logRepo:          logRepo,
		reminderProducer: reminderProducer,
	}
}

func (s *TaskService) CreateTask(task *domain.Task) error {
	// 1. 保存到MySQL
	if err := s.taskRepo.Create(task); err != nil {
		return err
	}

	// 2. 索引到ES
	if err := s.taskEsRepo.IndexTask(task); err != nil {
		fmt.Println(err)
		// 记录错误但不阻止流程
		// TODO: 可以考虑使用消息队列重试
	}

	// 3. 清除用户任务缓存
	if err := s.taskCache.InvalidateUserTasks(task.UserID); err != nil {
		// 记录错误但不阻止流程
	}

	// 4. 记录操作日志
	s.logRepo.AddLog(&mongodb.OperationLog{
		UserID:   task.UserID,
		Action:   "create",
		Target:   "task",
		TargetID: task.ID,
		Details:  map[string]string{"title": task.Title},
	})

	// 5. 如果有设置提醒，发送到消息队列
	if !task.ReminderAt.IsZero() {
		err := s.scheduleReminder(task)
		if err != nil {
			fmt.Println(err, 55)
		}
	}

	return nil
}

func (s *TaskService) UpdateTask(task *domain.Task) error {
	// 1. 保存到MySQL
	if err := s.taskRepo.Update(task); err != nil {
		return err
	}

	// 2. 更新ES索引
	if err := s.taskEsRepo.IndexTask(task); err != nil {
		// 记录错误但不阻止流程
	}

	// 3. 清除用户任务缓存
	if err := s.taskCache.InvalidateUserTasks(task.UserID); err != nil {
		// 记录错误但不阻止流程
	}

	// 4. 记录操作日志
	s.logRepo.AddLog(&mongodb.OperationLog{
		UserID:   task.UserID,
		Action:   "update",
		Target:   "task",
		TargetID: task.ID,
		Details:  map[string]string{"title": task.Title},
	})

	// 5. 如果提醒时间变更，重新调度提醒
	if !task.ReminderAt.IsZero() {
		s.scheduleReminder(task)
	}

	return nil
}

func (s *TaskService) DeleteTask(taskID uint, userID uint) error {
	// 1. 获取任务
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return err
	}

	// 验证所有者
	if task.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// 2. 从MySQL标记为删除
	if err := s.taskRepo.Delete(taskID); err != nil {
		return err
	}

	// 3. 从ES删除
	if err := s.taskEsRepo.DeleteTask(taskID); err != nil {
		// 记录错误但不阻止流程
	}

	// 4. 清除用户任务缓存
	if err := s.taskCache.InvalidateUserTasks(userID); err != nil {
		// 记录错误但不阻止流程
	}

	// 5. 记录操作日志
	s.logRepo.AddLog(&mongodb.OperationLog{
		UserID:   userID,
		Action:   "delete",
		Target:   "task",
		TargetID: taskID,
	})

	return nil
}

func (s *TaskService) GetUserTasks(userID uint) ([]domain.Task, error) {
	// 1. 尝试从缓存获取
	tasks, err := s.taskCache.GetUserTasks(userID)
	if err == nil {
		return tasks, nil
	}

	// 2. 缓存未命中，从MySQL获取
	tasks, err = s.taskRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 3. 更新缓存
	if err := s.taskCache.CacheUserTasks(userID, tasks); err != nil {
		// 记录错误但不阻止流程
	}

	return tasks, nil
}

func (s *TaskService) SearchTasks(query string, status domain.TaskStatus, userID uint) ([]domain.Task, error) {
	// 直接调用ES进行搜索
	return s.taskEsRepo.Search(query, status, userID)
}

// 调度任务提醒
func (s *TaskService) scheduleReminder(task *domain.Task) error {
	location, _ := time.LoadLocation("Asia/Shanghai")
	// 解析时间字符串

	timeInLocal := task.ReminderAt.In(location)
	fmt.Println("Parsed time (+8 timezone):", timeInLocal)
	message := map[string]interface{}{
		"task_id":     task.ID,
		"user_id":     task.UserID,
		"title":       task.Title,
		"reminder_at": timeInLocal,
	}

	// 发送到Rabbitmq
	err := s.reminderProducer.Publish("task_reminders", message, task.ReminderAt)
	if err != nil {
		fmt.Println(err, "ok")
		return err
	}
	fmt.Println("ok2")
	return nil
}
