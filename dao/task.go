package dao

import (
	"gorm.io/gorm"
	"tdl/types"
)

type DbRepo struct {
	db *gorm.DB
}

func NewDbRepository(db *gorm.DB) *DbRepo {
	return &DbRepo{db: db}
}

func (r *DbRepo) Create(task *types.Task) error {
	return r.db.Create(task).Error
}

func (r *DbRepo) Update(task *types.Task) error {
	return r.db.Save(task).Error
}

func (r *DbRepo) Delete(id uint) error {
	return r.db.Model(&types.Task{}).
		Where("id = ?", id).
		Update("status", types.TaskStatusDeleted).Error
}

func (r *DbRepo) FindByID(id uint) (*types.Task, error) {
	var task types.Task
	err := r.db.Where("id = ? AND status != ?", id, types.TaskStatusDeleted).First(&task).Error
	return &task, err
}

func (r *DbRepo) FindByUserID(userID uint) ([]types.Task, error) {
	var tasks []types.Task
	err := r.db.Where("user_id = ? AND status != ?", userID, types.TaskStatusDeleted).Find(&tasks).Error
	return tasks, err
}
