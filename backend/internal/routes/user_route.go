package routes

import (
	"github.com/gin-gonic/gin"
	"main.go/config"
	"main.go/internal/handlers"
	"main.go/internal/middleware"
)

type UserRouter interface {
	Mount()
}

type userRouterImpl struct {
	v       *gin.RouterGroup
	handler handlers.UserHandler
	db      config.GormPostgres
}

func NewUserRouter(v *gin.RouterGroup, handler handlers.UserHandler, db config.GormPostgres) UserRouter {
	return &userRouterImpl{v: v, handler: handler, db: db}
}

func (u *userRouterImpl) Mount() {
	u.v.POST("/login", u.handler.LoginUser)
	u.v.POST("/logout", u.handler.LogoutUser)

	u.v.Use(middleware.AuthMiddleware(u.db.GetConnection()))
	u.v.GET("/", u.handler.GetUsers)
	u.v.GET("/:id", u.handler.GetUserByID)
	u.v.POST("/", u.handler.CreateUser)
	u.v.PUT("/:id", u.handler.UpdateUser)
	u.v.DELETE("/:id", u.handler.DeleteUser)
}
