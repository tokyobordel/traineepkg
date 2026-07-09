// Package swagger подключает Swagger UI к Fiber-приложению.
package swagger

import (
	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "github.com/tokyobordel/traineepkg/adapters/api/v1/docs"
)

// SetupRouter подключает Swagger UI к Fiber-приложению.
// Документация доступна по адресу /auth/swagger/index.html.
func SetupRouter(app *fiber.App) {
	app.Get("/auth/swagger/*", swaggo.HandlerDefault)
}
