// Package sql internal/repository/sql/task_repository.go
package sql

import (
	"gorm.io/gorm"
	"tdl/types"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *types.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) Update(task *types.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Model(&types.Task{}).
		Where("id = ?", id).
		Update("status", types.TaskStatusDeleted).Error
}

func (r *TaskRepository) FindByID(id uint) (*types.Task, error) {
	var task types.Task
	err := r.db.Where("id = ? AND status != ?", id, types.TaskStatusDeleted).First(&task).Error
	return &task, err
}

func (r *TaskRepository) FindByUserID(userID uint) ([]types.Task, error) {
	var tasks []types.Task
	err := r.db.Where("user_id = ? AND status != ?", userID, types.TaskStatusDeleted).Find(&tasks).Error
	return tasks, err
}
