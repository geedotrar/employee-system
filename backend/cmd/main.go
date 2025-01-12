package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"main.go/config"
	"main.go/internal/handlers"
	"main.go/internal/repository"
	"main.go/internal/routes"
	"main.go/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	server()
}

func server() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	g := gin.Default()
	g.Use(gin.Recovery())

	usersGroup := g.Group("/users")
	gorm := config.NewGormPostgres()
	userRepo := repository.NewUserQuery(gorm)
	userSvc := service.NewUserService(userRepo)
	userHdl := handlers.NewUserHandler(userSvc)
	userRouter := routes.NewUserRouter(usersGroup, userHdl)
	userRouter.Mount()

	g.Run(":8080")
}
