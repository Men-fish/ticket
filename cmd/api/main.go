package main

import (
	"fmt"

	"github.com/Men-fish/ticket-v1/config"
	"github.com/Men-fish/ticket-v1/db"
	"github.com/Men-fish/ticket-v1/handlers"
	"github.com/Men-fish/ticket-v1/middlewares"
	"github.com/Men-fish/ticket-v1/repositories"
	"github.com/Men-fish/ticket-v1/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)

	app := fiber.New(fiber.Config{
		AppName:      "Ticket-Booking",
		ServerHeader: "Fiber",
	})

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8081",                       // Разрешить доступ только с вашего фронтенда
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE",                    // Разрешённые методы
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization", // Добавить 'Authorization' в разрешённые заголовки
		AllowCredentials: true,                                          // Разрешить использование куки и авторизационных заголовков
	}))

	// Repositories
	eventRepository := repositories.NewEventRepository(db)
	ticketRepository := repositories.NewTicketRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	// Service
	authService := services.NewAuthService(authRepository)

	// Routing
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	// Private routes protected by authentication middleware
	privateRoutes := server.Use(middlewares.AuthProtected(db))

	handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepository)
	handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepository)

	// Start the server
	app.Listen(fmt.Sprintf(":" + envConfig.ServerPort))
}
