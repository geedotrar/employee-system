package service

import (
	"context"

	"main.go/internal/models"
	"main.go/internal/repository"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)

	CreateUser(ctx context.Context, createUser models.UserRequest) (models.UserResponse, error)
	UpdateUser(ctx context.Context, id uint64, user models.UserRequest) (models.UserResponse, error)

	DeleteUser(ctx context.Context, id uint64) (models.UserResponse, error)
}

type userServiceImpl struct {
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (u *userServiceImpl) GetUserByID(ctx context.Context, id uint64) (models.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userServiceImpl) CreateUser(ctx context.Context, createUser models.UserRequest) (models.UserResponse, error) {
	user, err := u.repo.CreateUser(ctx, createUser)
	if err != nil {
		return models.UserResponse{}, err
	}
	return models.UserResponse{
		Status:  200,
		Message: "User created successfully",
		Data:    &user,
		Error:   false,
	}, nil
}

func (u *userServiceImpl) UpdateUser(ctx context.Context, id uint64, user models.UserRequest) (models.UserResponse, error) {
	updatedUser, err := u.repo.UpdateUser(ctx, id, user)
	if err != nil {
		return models.UserResponse{
			Status:  400,
			Message: "Failed to delete user",
			Data:    nil,
			Error:   true,
		}, err
	}
	return models.UserResponse{
		Status:  200,
		Message: "User updated successfully",
		Data:    &updatedUser,
		Error:   false,
	}, nil
}

func (u *userServiceImpl) DeleteUser(ctx context.Context, id uint64) (models.UserResponse, error) {
	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return models.UserResponse{
			Status:  400,
			Message: "Failed to delete user",
			Data:    nil,
			Error:   true,
		}, err
	}

	return models.UserResponse{
		Status:  200,
		Message: "User deleted successfully",
		Data:    nil,
		Error:   false,
	}, nil
}
