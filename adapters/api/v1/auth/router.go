package auth

import (
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "github.com/tokyobordel/traineepkg/adapters/api/v1/docs"
)

// SetupRouter регистрирует публичные маршруты аутентификации под префиксом /api/auth.
func SetupRouter(app *fiber.App, handler *Handler) {
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.Refresh)
	authGroup.Post("/logout", handler.Logout)
	authGroup.Get("/me", handler.GetMe)
}

// SetupAuthRouter подключает Swagger UI документации auth API по адресу /auth/swagger.
func SetupAuthRouter(app *fiber.App) {
	app.Get("/auth/swagger/*", swaggo.HandlerDefault)
}
