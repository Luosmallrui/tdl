package domain

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusDeleted    TaskStatus = "deleted"
)

// Task 表
type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"size:255;not null"`
	Description string     `json:"description" gorm:"type:text"`
	Status      TaskStatus `json:"status" gorm:"size:20;default:'pending'"`
	UserID      uint       `json:"user_id" gorm:"index"`
	DueDate     time.Time  `json:"due_date"`
	ReminderAt  time.Time  `json:"reminder_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Tags        string     `json:"tags" gorm:"column:tags"` // 多个标签，逗号分隔
}

type CreateTask struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description" `
	Status      TaskStatus `json:"status" `
	UserID      uint       `json:"user_id"`
	DueDate     time.Time  `json:"due_date"`
	ReminderAt  time.Time  `json:"reminder_at"`
	CreatedAt   time.Time  `json:"created_at" `
	UpdatedAt   time.Time  `json:"updated_at"`
	Tags        []string   `json:"tags" `
}

type ListTask struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description" `
	Status      TaskStatus `json:"status" `
	UserID      uint       `json:"user_id"`
	DueDate     time.Time  `json:"due_date"`
	ReminderAt  time.Time  `json:"reminder_at"`
	CreatedAt   time.Time  `json:"created_at" `
	UpdatedAt   time.Time  `json:"updated_at"`
	Tags        []string   `json:"tags" `
}
