package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"traineesheep/imageservice/pkg/errors"
	"traineesheep/imageservice/pkg/logger"
)

// mapDomainErrorToHTTPStatus сопоставляет доменные ошибки с HTTP статус кодами
func mapDomainErrorToHTTPStatus(err errors.DomainError) int {
	if err == nil {
		return http.StatusOK
	}

	switch err.Type() {
	case errors.ErrorTypeNotFound:
		return http.StatusNotFound // 404

	case errors.ErrorTypeUniqueConstraint:
		return http.StatusConflict // 409

	case errors.ErrorTypeInvalidParameters:
		return http.StatusBadRequest // 400

	case errors.ErrorTypeValidation:
		return http.StatusBadRequest // 400

	case errors.ErrorTypeUnauthorized:
		return http.StatusUnauthorized // 401

	case errors.ErrorTypeInternalService:
		return http.StatusInternalServerError // 500

	case errors.ErrorTypeAccessDenied:
		return http.StatusForbidden // 403

	case errors.ErrorTypeToken:
		return http.StatusUnauthorized // 401

	case errors.ErrorTypeSignature:
		return http.StatusUnauthorized // 401

	default:
		return http.StatusInternalServerError // 500 по умолчанию
	}
}

// makeErrorResponse создает ответ с ошибкой, автоматически определяя HTTP статус
func MakeErrorResponse(c *fiber.Ctx, logger *logger.ContextLogger, err error) {
	_ = logger

	if domainErr, ok := err.(errors.DomainError); ok {
		// Это доменная ошибка, используем правильный маппинг
		status := mapDomainErrorToHTTPStatus(domainErr)
		makeResponse(c, status, nil, false, err.Error())
	} else {
		// Обычная ошибка, возвращаем 500
		makeResponse(c, http.StatusInternalServerError, nil, false, errors.NewInternalServiceError("internal server error", nil).Error())
	}
}
