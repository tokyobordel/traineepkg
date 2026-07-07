package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/tokyobordel/traineepkg/adapters/api/v1/response"
	authService "github.com/tokyobordel/traineepkg/auth/service"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"
	"github.com/tokyobordel/traineepkg/errors"
	"github.com/tokyobordel/traineepkg/logger"
)

type Handler struct {
	authService authService.IAuthService
	jwtService  *jwtAuth.Service
	logger      *logger.ContextLogger
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewHandler(authService authService.IAuthService, jwtService *jwtAuth.Service, logger *logger.ContextLogger, accessTTL time.Duration, refreshTTL time.Duration) *Handler {
	return &Handler{
		authService: authService,
		jwtService:  jwtService,
		logger:      logger,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req credentialsRequest
	if err := c.BodyParser(&req); err != nil || req.Login == "" || req.Pass == "" {
		response.MakeErrorResponse(c, h.logger, errors.NewInvalidParametersError("body", nil, "invalid request body"))
		return nil
	}

	user, err := h.authService.Register(req.Pass, req.Login)
	if err != nil {
		response.MakeErrorResponse(c, h.logger, err)
		return nil
	}

	response.MakeSuccessResponseWithStatus(c, fiber.StatusCreated, authResponse{User: user})
	return nil
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req credentialsRequest
	if err := c.BodyParser(&req); err != nil || req.Login == "" || req.Pass == "" {
		response.MakeErrorResponse(c, h.logger, errors.NewInvalidParametersError("body", nil, "invalid request body"))
		return nil
	}

	user, err := h.authService.Login(req.Pass, req.Login)
	if err != nil {
		response.MakeErrorResponse(c, h.logger, err)
		return nil
	}

	if !h.setTokenCookies(c, user.ID) {
		response.MakeErrorResponse(c, h.logger, errors.NewInternalServiceError("failed to generate jwt tokens", nil))
		return nil
	}

	response.MakeSuccessResponse(c, authResponse{User: user})
	return nil
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies(jwtAuth.RefreshTokenCookieName)
	if refreshToken == "" {
		response.MakeErrorResponse(c, h.logger, errors.NewAuthTokenError(errors.TokenNotFound))
		return nil
	}

	userID, err := h.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		response.MakeErrorResponse(c, h.logger, errors.NewAuthTokenError(errors.InvalidToken))
		return nil
	}

	user, getMeErr := h.authService.GetMe(userID)
	if getMeErr != nil {
		response.MakeErrorResponse(c, h.logger, getMeErr)
		return nil
	}

	if !h.setTokenCookies(c, userID) {
		response.MakeErrorResponse(c, h.logger, errors.NewInternalServiceError("failed to generate jwt tokens", nil))
		return nil
	}

	response.MakeSuccessResponse(c, authResponse{User: user})
	return nil
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     jwtAuth.AccessTokenCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     jwtAuth.RefreshTokenCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		HTTPOnly: true,
	})

	response.MakeSuccessResponse(c, fiber.Map{"message": "logout successful"})
	return nil
}

func (h *Handler) setTokenCookies(c *fiber.Ctx, userID int) bool {
	tokenPair, err := h.jwtService.GenerateTokenPair(userID)
	if err != nil {
		return false
	}

	c.Cookie(&fiber.Cookie{
		Name:     jwtAuth.AccessTokenCookieName,
		Value:    tokenPair.AccessToken,
		Path:     "/",
		MaxAge:   int(h.accessTTL.Seconds()),
		HTTPOnly: true,
		Secure:   false,
	})
	c.Cookie(&fiber.Cookie{
		Name:     jwtAuth.RefreshTokenCookieName,
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		MaxAge:   int(h.refreshTTL.Seconds()),
		HTTPOnly: true,
		Secure:   false,
	})
	return true
}
