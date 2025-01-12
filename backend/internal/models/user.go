package models

import "time"

type UsersResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    *[]User `json:"data"`
	Error   bool    `json:"error"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    *User  `json:"data"`
	Error   bool   `json:"error"`
}

type User struct {
	Id         int       `json:"id"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	RoleID     int       `json:"role_id"`
	PositionID int       `json:"position_id"`
	Role       Role      `json:"role" gorm:"foreignKey:RoleID"`
	Position   Position  `json:"position" gorm:"foreignKey:PositionID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Position struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type UserRequest struct {
	Firstname  string `json:"firstname" binding:"required"`
	Lastname   string `json:"lastname" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RoleID     int    `json:"role_id" binding:"required"`
	PositionID int    `json:"position_id" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
