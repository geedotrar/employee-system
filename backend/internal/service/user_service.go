package service

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"main.go/internal/auth"
	"main.go/internal/models"
	"main.go/internal/repository"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	CreateUser(ctx context.Context, createUser models.UserRequest) (models.User, error)
	UpdateUser(ctx context.Context, id uint64, user models.UserRequest) (models.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	Login(ctx context.Context, email string, password string) (models.AuthResponse, error)
	Logout(ctx context.Context, token string) error
}

type userServiceImpl struct {
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	user, err := u.repo.GetUsers(ctx)
	if err != nil {
		return []models.User{}, err
	}
	return user, nil
}

func (u *userServiceImpl) CreateUser(ctx context.Context, createUser models.UserRequest) (models.User, error) {
	existingUser, err := u.repo.GetUserByEmail(ctx, createUser.Email)
	if err != nil {
		return models.User{}, err
	}
	if existingUser.Id != 0 {
		return models.User{}, errors.New("email already exists")
	}

	role, err := u.repo.GetRoleByID(ctx, uint64(createUser.RoleID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("role not found")
		}
		return models.User{}, err
	}

	position, err := u.repo.GetPositionByID(ctx, uint64(createUser.PositionID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("position not found")
		}
		return models.User{}, err
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(createUser.Password)
	if err != nil {
		return models.User{}, errors.New("error hashing password")
	}

	newUser := models.User{
		Firstname:  createUser.Firstname,
		Lastname:   createUser.Lastname,
		Email:      createUser.Email,
		Password:   hashedPassword,
		RoleID:     role.ID,
		PositionID: position.ID,
	}

	createdUser, err := u.repo.CreateUser(ctx, newUser)
	if err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (u *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *userServiceImpl) UpdateUser(ctx context.Context, id uint64, updateUser models.UserRequest) (models.User, error) {
	existingUser, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	if existingUser.Id == 0 {
		return models.User{}, errors.New("user not found")
	}

	userWithSameEmail, err := u.repo.GetUserByEmail(ctx, updateUser.Email)
	if err != nil {
		return models.User{}, err
	}
	if userWithSameEmail.Id != 0 && userWithSameEmail.Id != int(id) {
		return models.User{}, errors.New("email already exists")
	}

	role, err := u.repo.GetRoleByID(ctx, uint64(updateUser.RoleID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("role not found")
		}
		return models.User{}, err
	}

	position, err := u.repo.GetPositionByID(ctx, uint64(updateUser.PositionID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("position not found")
		}
		return models.User{}, err
	}

	updatedPassword := existingUser.Password
	if updateUser.Password != "" && updateUser.Password != existingUser.Password {
		hashedPassword, err := auth.HashPassword(updateUser.Password)
		if err != nil {
			return models.User{}, errors.New("error hashing password")
		}
		updatedPassword = hashedPassword
	}

	updatedUser := models.User{
		Id:         existingUser.Id,
		Firstname:  updateUser.Firstname,
		Lastname:   updateUser.Lastname,
		Email:      updateUser.Email,
		Password:   updatedPassword,
		RoleID:     role.ID,
		PositionID: position.ID,
	}

	savedUser, err := u.repo.UpdateUser(ctx, id, updatedUser)
	if err != nil {
		return models.User{}, err
	}

	return savedUser, nil
}

func (u *userServiceImpl) DeleteUser(ctx context.Context, id uint64) error {
	err := u.repo.DeleteUser(ctx, id)
	return err
}

func (u *userServiceImpl) Login(ctx context.Context, email string, password string) (models.AuthResponse, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil || user.Id == 0 {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}
	if !auth.CheckPasswordHash(password, user.Password) {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}
	if !auth.CheckPasswordHash(password, user.Password) {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	token, err := auth.GenerateJWT(user.Email)
	if err != nil {
		return models.AuthResponse{}, err
	}
	return models.AuthResponse{
		Status:  200,
		Message: "Login successful",
		Token:   token,
		Error:   false,
	}, nil
}

func (u *userServiceImpl) GetUserByID(ctx context.Context, id uint64) (models.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	if user.Id == 0 {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func (u *userServiceImpl) Logout(ctx context.Context, token string) error {

	return nil
}
