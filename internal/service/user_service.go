// internal/service/user_service.go
package service

import (
	"tdl/internal/domain"
	"tdl/internal/repository/sql"
)

type UserService struct {
	userRepo *sql.UserRepository
}

func NewUserService(userRepo *sql.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.FindByUsername(username)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(email)
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(user *domain.User) error {
	return s.userRepo.Create(user)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *domain.User) error {
	return s.userRepo.Update(user)
}
