package auth

import "github.com/gofiber/fiber/v2"

func SetupRouter(app *fiber.App, handler *Handler) {
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.Refresh)
	authGroup.Post("/logout", handler.Logout)
	authGroup.Get("/me", handler.GetMe)
}
