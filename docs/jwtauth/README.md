# JWT Auth

Пакет даёт готовые HTTP-эндпоинты аутентификации и middleware для защиты маршрутов. Токены хранятся в HTTP-only cookies (`access_token`, `refresh_token`).

Бизнес-логику пользователей вы реализуете сами через `auth/service.IAuthService` — пакет не привязан к конкретной БД.

## Подключение

Один экземпляр `jwt.Service` нужен и для auth-хендлера, и для middleware: оба читают и проверяют одни и те же токены.

```go
package main

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/tokyobordel/traineepkg/adapters/api/v1/auth"
	authMiddleware "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware/authjwt"
	apiMiddleware "github.com/tokyobordel/traineepkg/adapters/api/v1/middleware"
	authSwagger "github.com/tokyobordel/traineepkg/adapters/api/v1/swagger"
	"github.com/tokyobordel/traineepkg/adapters/api/v1/response"
	authService "github.com/tokyobordel/traineepkg/auth/service"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"
	"github.com/tokyobordel/traineepkg/logger"
)

func main() {
	app := fiber.New()

	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	var myAuthService authService.IAuthService // ваша реализация Login/Register/GetMe
	jwtService := jwtAuth.NewService("your-secret", accessTTL, refreshTTL)

	authHandler := auth.NewHandler(myAuthService, jwtService, log, accessTTL, refreshTTL)
	authMW := authMiddleware.NewMiddleware(jwtService)

	auth.SetupRouter(app, authHandler)       // /api/auth/*
	authSwagger.SetupRouter(app)             // /auth/swagger/index.html

	api := app.Group("/api")
	protected := api.Group("/", authMW.RequireAccessToken())
	protected.Get("/admin/images/:id", getImageByIDAdmin)

	app.Listen(":8080")
}

func getImageByIDAdmin(c fiber.Ctx) error {
	userID, ok := authMiddleware.UserIDFromContext(c.Context())
	if !ok {
		return nil // middleware уже ответил 401
	}

	_ = userID
	// ... ваша логика

	response.MakeSuccessResponse(c, fiber.Map{"id": c.Params("id")})
	return nil
}
```

## Что подключается

| Компонент | Назначение |
|-----------|------------|
| `jwtAuth.NewService` | Генерация и валидация access/refresh токенов |
| `auth.NewHandler` + `auth.SetupRouter` | Публичные маршруты: register, login, refresh, logout, me |
| `authjwt.NewMiddleware` | Проверка `access_token` cookie на защищённых маршрутах |
## Защищённые маршруты

`RequireAccessToken()` читает cookie `access_token`, валидирует JWT и кладёт `userID` в контекст запроса. При ошибке клиент получает стандартный JSON-ответ с `success: false`.

В хендлере ID пользователя берётся так:

```go
userID, ok := authMiddleware.UserIDFromContext(c.Context())
```

Группу с middleware можно вешать на любой префикс — например, только `/api/admin`:

```go
admin := app.Group("/api/admin", authMW.RequireAccessToken())
admin.Delete("/users/:id", deleteUser)
```

## Публичные эндпоинты

После `POST /api/auth/login` браузер получает cookies автоматически. Для защищённых маршрутов отдельно передавать заголовок не нужно — middleware читает cookie.

| Метод | Путь | Доступ |
|-------|------|--------|
| POST | `/api/auth/register` | публичный |
| POST | `/api/auth/login` | публичный, выдаёт cookies |
| POST | `/api/auth/refresh` | публичный, обновляет cookies |
| POST | `/api/auth/logout` | публичный, очищает cookies |
| GET | `/api/auth/me` | проверяет cookie внутри хендлера |

Эндпоинт `/api/auth/me` защищён логикой самого хендлера. Остальные защищённые маршруты приложения — через `authjwt` middleware.




## Установка документации Swagger

```go
package swagger

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"

	_ "github.com/tokyobordel/traineepkg/adapters/api/v1/docs"
)

func SetupRouter(app *fiber.App) {
	app.Get("/auth/swagger/*", fiberSwagger.New(fiberSwagger.Config{
		InstanceName: "swagger",
	}))
}

```