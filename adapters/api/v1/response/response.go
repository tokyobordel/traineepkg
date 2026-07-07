package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"traineepkg/context/trace"
)

func makeEnvelope(c *fiber.Ctx, data interface{}, success bool, errMessage string) fiber.Map {
	spreadId, ok := trace.SpreadFromContext(c.UserContext())
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

func makeResponse(c *fiber.Ctx, status int, data interface{}, success bool, errMessage string) {
	c.Status(status).JSON(makeEnvelope(c, data, success, errMessage))
}

func MakeSuccessResponse(c *fiber.Ctx, data interface{}) {
	makeResponse(c, http.StatusOK, data, true, "")
}

func MakeSuccessResponseWithStatus(c *fiber.Ctx, status int, data interface{}) {
	makeResponse(c, status, data, true, "")
}
