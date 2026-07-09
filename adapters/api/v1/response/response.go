// Package response формирует унифицированные HTTP-ответы Fiber API.
package response

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/tokyobordel/traineepkg/context/trace"
)

func makeEnvelope(c fiber.Ctx, data interface{}, success bool, errMessage string) fiber.Map {
	spreadId, ok := trace.SpreadFromContext(c.Context())
	if !ok {
		spreadId = "unknown"
	}

	return fiber.Map{
		"data":        data,
		"success":     success,
		"err_message": errMessage,
		"spread_id":   spreadId,
	}
}

func makeResponse(c fiber.Ctx, status int, data interface{}, success bool, errMessage string) {
	c.Status(status).JSON(makeEnvelope(c, data, success, errMessage))
}

// MakeSuccessResponse отправляет успешный JSON-ответ со статусом 200.
func MakeSuccessResponse(c fiber.Ctx, data interface{}) {
	makeResponse(c, http.StatusOK, data, true, "")
}

// MakeSuccessResponseWithStatus отправляет успешный JSON-ответ с указанным HTTP-статусом.
func MakeSuccessResponseWithStatus(c fiber.Ctx, status int, data interface{}) {
	makeResponse(c, status, data, true, "")
}
