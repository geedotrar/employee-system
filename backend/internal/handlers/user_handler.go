package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/models"
	"main.go/internal/service"
)

type UserHandler interface {
	GetUsers(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)

	LoginUser(ctx *gin.Context)
	LogoutUser(ctx *gin.Context)
}

type userHandlerImpl struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandlerImpl{svc: svc}
}

func (u *userHandlerImpl) GetUsers(ctx *gin.Context) {
	users, err := u.svc.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.UsersResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get users",
			Data:    nil,
			Error:   true,
		})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, models.UsersResponse{
			Status:  http.StatusNotFound,
			Message: "Users not found",
			Data:    nil,
			Error:   false,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.UsersResponse{
		Status:  http.StatusOK,
		Message: "Success to get users",
		Data:    &users,
		Error:   false,
	})
}

func (u *userHandlerImpl) GetUserByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	user, err := u.svc.GetUserByID(ctx, uint64(id))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, models.UserResponse{
				Status:  http.StatusNotFound,
				Message: "User not found",
				Data:    nil,
				Error:   true,
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, models.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to get user",
				Data:    nil,
				Error:   true,
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "Success to get user",
		Data:    &user,
		Error:   false,
	})
}

func (u *userHandlerImpl) CreateUser(ctx *gin.Context) {
	var createUserRequest models.UserRequest
	if err := ctx.ShouldBindJSON(&createUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request data",
			Data:    nil,
			Error:   true,
		})
		return
	}
	userResponse, err := u.svc.CreateUser(ctx, createUserRequest)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" ||
			err.Error() == "role not found" ||
			err.Error() == "position not found" {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, models.UserResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "User created successfully",
		Data:    &userResponse,
		Error:   false,
	})
}

func (u *userHandlerImpl) UpdateUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	var updateUserRequest models.UserRequest
	if err := ctx.ShouldBindJSON(&updateUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request data",
			Data:    nil,
			Error:   true,
		})
		return
	}

	userResponse, err := u.svc.UpdateUser(ctx, uint64(id), updateUserRequest)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" ||
			err.Error() == "role not found" ||
			err.Error() == "position not found" {
			statusCode = http.StatusBadRequest
		} else if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, models.UserResponse{
			Status:  statusCode,
			Message: err.Error(),
			Data:    nil,
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data:    &userResponse,
		Error:   false,
	})
}

func (u *userHandlerImpl) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, models.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid or missing ID parameter",
			Data:    nil,
			Error:   true,
		})
		return
	}

	user, _ := u.svc.GetUserByID(ctx, uint64(id))
	if user.Id == 0 {
		ctx.JSON(http.StatusNotFound, models.UserResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Error:   false,
		})
		return
	}

	err = u.svc.DeleteUser(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.UserResponse{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
		Data:    nil,
		Error:   false,
	})
}

func (u *userHandlerImpl) LoginUser(ctx *gin.Context) {
	var authUserRequest models.AuthRequest

	if err := ctx.ShouldBindJSON(&authUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := u.svc.Login(ctx, authUserRequest.Email, authUserRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (u *userHandlerImpl) LogoutUser(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	err := u.svc.Logout(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
