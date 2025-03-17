// Package sql internal/repository/sql/task_repository.go
package sql

import (
	"gorm.io/gorm"
	"tdl/internal/domain"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Model(&domain.Task{}).
		Where("id = ?", id).
		Update("status", domain.TaskStatusDeleted).Error
}

func (r *TaskRepository) FindByID(id uint) (*domain.Task, error) {
	var task domain.Task
	err := r.db.Where("id = ? AND status != ?", id, domain.TaskStatusDeleted).First(&task).Error
	return &task, err
}

func (r *TaskRepository) FindByUserID(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("user_id = ? AND status != ?", userID, domain.TaskStatusDeleted).Find(&tasks).Error
	return tasks, err
}
