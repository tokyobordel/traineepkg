// Package docs содержит сгенерированную Swagger/OpenAPI документацию API.
//
// @title           TraineePkg Auth API
// @version         1.0
// @description     HTTP API аутентификации пакета traineepkg.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
//
// @license.name  MIT
//
// @host      localhost:8080
// @BasePath  /api
//
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token
package docs

//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g spec.go -d .,../auth,../response -o . --parseDependency --parseInternal
