package services

import (
	"context"
	"infinite-experiment/infinite-experiment-backend/internal/db/repositories"
	"infinite-experiment/infinite-experiment-backend/internal/models/entities"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, user *entities.User) error {
	return s.repo.InsertUser(ctx, user)
}
