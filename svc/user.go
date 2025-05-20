package svc

import (
	"tdl/dao"
	"tdl/types"
)

var _ IUserService = (*UserService)(nil)

type IUserService interface {
	GetUserByID(id uint) (*types.User, error)
	GetUserByUsername(username string) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	CreateUser(user *types.User) error
	UpdateUser(user *types.User) error
}
type UserService struct {
	UserRepo *dao.UserRepo
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*types.User, error) {
	return s.UserRepo.FindByID(id)
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*types.User, error) {
	return s.UserRepo.FindByUsername(username)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*types.User, error) {
	return s.UserRepo.FindByEmail(email)
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(user *types.User) error {
	return s.UserRepo.Create(user)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *types.User) error {
	return s.UserRepo.Update(user)
}
