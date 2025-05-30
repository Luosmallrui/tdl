package svc

import (
	"fmt"
	"tdl/dao"
	"tdl/dao/cache"
	"tdl/types"
	"time"
)

var _ ITaskService = (*TaskService)(nil)

type TaskService struct {
	TaskRepo         *dao.DbRepo
	TaskCache        *cache.TaskCache
	TaskEsRepo       *dao.EsRepo
	LogRepo          *dao.LogRepository
	ReminderProducer *dao.RabbitMQProducer
}

type ITaskService interface {
	CreateTask(task *types.Task) error
	UpdateTask(task *types.Task) error
	DeleteTask(taskID uint, userID uint) error
	GetUserTasks(userID uint) ([]types.Task, error)
	SearchTasks(query string, status types.TaskStatus, userID uint) ([]types.Task, error)
}

func (s *TaskService) CreateTask(task *types.Task) error {
	// 1. 保存到MySQL
	if err := s.TaskRepo.Create(task); err != nil {
		return err
	}

	// 2. 索引到ES
	if err := s.TaskEsRepo.IndexTask(task); err != nil {
		fmt.Println(err)
		// 记录错误但不阻止流程
		// TODO: 可以考虑使用消息队列重试
	}

	// 3. 清除用户任务缓存
	if err := s.TaskCache.InvalidateUserTasks(task.UserID); err != nil {
		// 记录错误但不阻止流程
	}

	// 4. 记录操作日志
	s.LogRepo.AddLog(&dao.OperationLog{
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

func (s *TaskService) UpdateTask(task *types.Task) error {
	// 1. 保存到MySQL
	if err := s.TaskRepo.Update(task); err != nil {
		return err
	}

	// 2. 更新ES索引
	if err := s.TaskEsRepo.IndexTask(task); err != nil {
		// 记录错误但不阻止流程
	}

	// 3. 清除用户任务缓存
	if err := s.TaskCache.InvalidateUserTasks(task.UserID); err != nil {
		// 记录错误但不阻止流程
	}

	// 4. 记录操作日志
	s.LogRepo.AddLog(&dao.OperationLog{
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
	task, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		return err
	}

	// 验证所有者
	if task.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// 2. 从MySQL标记为删除
	if err := s.TaskRepo.Delete(taskID); err != nil {
		return err
	}

	// 3. 从ES删除
	if err := s.TaskEsRepo.DeleteTask(taskID); err != nil {
		// 记录错误但不阻止流程
	}

	// 4. 清除用户任务缓存
	if err := s.TaskCache.InvalidateUserTasks(userID); err != nil {
		// 记录错误但不阻止流程
	}

	// 5. 记录操作日志
	s.LogRepo.AddLog(&dao.OperationLog{
		UserID:   userID,
		Action:   "delete",
		Target:   "task",
		TargetID: taskID,
	})

	return nil
}

func (s *TaskService) GetUserTasks(userID uint) ([]types.Task, error) {
	// 1. 尝试从缓存获取
	tasks, err := s.TaskCache.GetUserTasks(userID)
	if err == nil {
		return tasks, nil
	}

	// 2. 缓存未命中，从MySQL获取
	tasks, err = s.TaskRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 3. 更新缓存
	if err := s.TaskCache.CacheUserTasks(userID, tasks); err != nil {
		// 记录错误但不阻止流程
	}

	return tasks, nil
}

func (s *TaskService) SearchTasks(query string, status types.TaskStatus, userID uint) ([]types.Task, error) {
	// 直接调用ES进行搜索
	return s.TaskEsRepo.Search(query, status, userID)
}

// 调度任务提醒
func (s *TaskService) scheduleReminder(task *types.Task) error {
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
	err := s.ReminderProducer.Publish("task_reminders", message, task.ReminderAt)
	if err != nil {
		fmt.Println(err, "ok")
		return err
	}
	fmt.Println("ok2")
	return nil
}
