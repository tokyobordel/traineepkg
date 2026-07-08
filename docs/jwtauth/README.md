# Auth API

HTTP-адаптер аутентификации для Fiber: регистрация, логин, refresh, logout, `/me`.

Реализуйте `auth/service.IAuthService` в своём сервисе — пакет не знает, где хранятся пользователи, и работает только через этот интерфейс (`Login`, `Register`, `GetMe`).

```go
import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/tokyobordel/traineepkg/adapters/api/v1/auth"
	authSwagger "github.com/tokyobordel/traineepkg/adapters/api/v1/swagger"
	authService "github.com/tokyobordel/traineepkg/auth/service"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"
	"github.com/tokyobordel/traineepkg/logger"
)

app := fiber.New()

var myAuthService authService.IAuthService // ваша реализация
log, _ := logger.NewContextLogger("", "", false)
jwt := jwtAuth.NewService("secret", 15*time.Minute, 7*24*time.Hour)

handler := auth.NewHandler(myAuthService, jwt, log, 15*time.Minute, 7*24*time.Hour)
auth.SetupRouter(app, handler)
authSwagger.SetupRouter(app) // /auth/swagger/index.html
```


## - Добавление защищенного end-point

```go
import (

	authMiddleware "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"

	"github.com/gofiber/fiber/v2"
)

var app *fiber.App
//...
api := app.Group("/api")
protected := api.Group("/", h.authMiddleware.RequireAccessToken())
protected.Get("/admin/image/:id", h.GetImageByIdAdmin)

```
