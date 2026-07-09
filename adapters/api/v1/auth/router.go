package auth

import (
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	authjwt "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"

	_ "github.com/tokyobordel/traineepkg/adapters/api/v1/docs"
)

// SetupRouter регистрирует публичные маршруты аутентификации под префиксом /api/auth.
func SetupRouter(app *fiber.App, handler *Handler) {
	defGroup := app.Group("/api/auth")
	defGroup.Post("/register", handler.Register)
	defGroup.Post("/login", handler.Login)
	defGroup.Post("/refresh", handler.Refresh)
	defGroup.Post("/logout", handler.Logout)

	authMiddleware := authjwt.NewMiddleware(handler.jwtService)
	defGroup.Get("/me", authMiddleware.RequireAccessToken(), handler.GetMe)
}

// SetupAuthRouter подключает Swagger UI документации auth API по адресу /auth/swagger.
func SetupAuthRouter(app *fiber.App) {
	app.Get("/auth/swagger/*", swaggo.HandlerDefault)
}
