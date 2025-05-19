// internal/repository/sql/user_repository.go
package sql

import (
	"gorm.io/gorm"
	"tdl/types"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *types.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *types.User) error {
	return r.db.Save(user).Error
}

// FindByID 通过ID查找用户
func (r *UserRepository) FindByID(id uint) (*types.User, error) {
	var user types.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

// FindByUsername 通过用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*types.User, error) {
	var user types.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// FindByEmail 通过邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*types.User, error) {
	var user types.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
