package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"main.go/config"
	"main.go/internal/models"
)

type UserQuery interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)
	CreateUser(ctx context.Context, user models.UserRequest) (models.User, error)
	UpdateUser(ctx context.Context, id uint64, user models.UserRequest) (models.User, error)

	DeleteUser(ctx context.Context, id uint64) error
}

type userQueryImpl struct {
	db config.GormPostgres
}

func NewUserQuery(db config.GormPostgres) UserQuery {
	return &userQueryImpl{db: db}
}

func (u *userQueryImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	db := u.db.GetConnection()
	users := []models.User{}
	if err := db.
		WithContext(ctx).
		Preload("Role").
		Preload("Position").
		Table("users").
		Where("deleted_at IS NULL").
		Find(&users).Error; err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (u *userQueryImpl) GetUserByID(ctx context.Context, id uint64) (models.User, error) {
	db := u.db.GetConnection()
	users := models.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Preload("Role").
		Preload("Position").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Find(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return users, nil
}

func (u *userQueryImpl) CreateUser(ctx context.Context, user models.UserRequest) (models.User, error) {
	var existingUser models.User
	db := u.db.GetConnection()

	if err := db.WithContext(ctx).Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return models.User{}, errors.New("email already exists")
	}

	var role models.Role
	if err := db.WithContext(ctx).Where("id = ?", user.RoleID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("role not found")
		}
		return models.User{}, err
	}

	var position models.Position
	if err := db.WithContext(ctx).Where("id = ?", user.PositionID).First(&position).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("position not found")
		}
		return models.User{}, err
	}

	newUser := models.User{
		Firstname:  user.Firstname,
		Lastname:   user.Lastname,
		Email:      user.Email,
		Password:   user.Password,
		RoleID:     user.RoleID,
		PositionID: user.PositionID,
	}

	if err := db.WithContext(ctx).Create(&newUser).Error; err != nil {
		return models.User{}, err
	}

	if err := db.WithContext(ctx).
		Preload("Role").
		Preload("Position").
		First(&newUser, newUser.Id).Error; err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (u *userQueryImpl) UpdateUser(ctx context.Context, id uint64, user models.UserRequest) (models.User, error) {
	db := u.db.GetConnection()

	var existingUser models.User
	if err := db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("user not found or already deleted")
		}
		return models.User{}, err
	}

	var emailCheck models.User
	if err := db.WithContext(ctx).Where("email = ? AND id != ?", user.Email, id).First(&emailCheck).Error; err == nil {
		return models.User{}, errors.New("email already exists")
	}

	var role models.Role
	if err := db.WithContext(ctx).Where("id = ?", user.RoleID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("role not found")
		}
		return models.User{}, err
	}

	var position models.Position
	if err := db.WithContext(ctx).Where("id = ?", user.PositionID).First(&position).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, errors.New("position not found")
		}
		return models.User{}, err
	}

	updatedUser := existingUser
	updatedUser.Firstname = user.Firstname
	updatedUser.Lastname = user.Lastname
	updatedUser.Email = user.Email
	updatedUser.Password = user.Password
	updatedUser.RoleID = user.RoleID
	updatedUser.PositionID = user.PositionID

	if err := db.WithContext(ctx).Save(&updatedUser).Error; err != nil {
		return models.User{}, err
	}

	if err := db.WithContext(ctx).
		Preload("Role").
		Preload("Position").
		First(&updatedUser, updatedUser.Id).Error; err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

func (u *userQueryImpl) DeleteUser(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()

	if err := db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error; err != nil {
		return err
	}

	return nil
}
