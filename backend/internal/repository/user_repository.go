package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"main.go/config"
	"main.go/internal/models"
)

type UserQuery interface {
	// GET
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetRoleByID(ctx context.Context, id uint64) (models.Role, error)
	GetPositionByID(ctx context.Context, id uint64) (models.Position, error)

	// POST
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateUser(ctx context.Context, id uint64, user models.User) (models.User, error)
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

func (u *userQueryImpl) GetRoleByID(ctx context.Context, id uint64) (models.Role, error) {
	db := u.db.GetConnection()
	var role models.Role

	if err := db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Role{}, errors.New("role not found")
		}
		return models.Role{}, err
	}
	return role, nil
}

func (u *userQueryImpl) GetPositionByID(ctx context.Context, id uint64) (models.Position, error) {
	db := u.db.GetConnection()
	var position models.Position

	if err := db.WithContext(ctx).Where("id = ?", id).First(&position).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Position{}, errors.New("position not found")
		}
		return models.Position{}, err
	}

	return position, nil
}

func (u *userQueryImpl) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	db := u.db.GetConnection()

	user := models.User{}
	if err := db.WithContext(ctx).Where("email = ?", email).
		Where("deleted_at is NULL").
		First(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	db := u.db.GetConnection()
	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		return models.User{}, err
	}

	if err := db.WithContext(ctx).
		Preload("Role").
		Preload("Position").
		First(&user, user.Id).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *userQueryImpl) UpdateUser(ctx context.Context, id uint64, user models.User) (models.User, error) {
	db := u.db.GetConnection()

	tx := db.Begin()

	var existingUser models.User
	if err := tx.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&existingUser).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found or already deleted")
		}
		return models.User{}, err
	}

	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.Email = user.Email
	existingUser.Password = user.Password
	existingUser.RoleID = user.RoleID
	existingUser.PositionID = user.PositionID

	if err := tx.WithContext(ctx).Save(&existingUser).Error; err != nil {
		tx.Rollback()

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return models.User{}, errors.New("email already exists")
		}

		return models.User{}, err
	}

	tx.Commit()

	if err := db.WithContext(ctx).
		Preload("Role").
		Preload("Position").
		First(&existingUser, existingUser.Id).Error; err != nil {
		return models.User{}, err
	}

	return existingUser, nil
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
