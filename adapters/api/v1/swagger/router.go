package swagger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/tokyobordel/traineepkg/adapters/api/v1/docs"
)

// SetupRouter подключает Swagger UI к Fiber-приложению.
// Документация доступна по адресу /auth/swagger/index.html.
func SetupRouter(app *fiber.App) {
	app.Get("/auth/swagger/*", swagger.HandlerDefault)
}
