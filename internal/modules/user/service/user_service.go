package service

import (
	"context"
	"github.com/alternative/backend/internal/modules/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}
