package auth

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/tokyobordel/traineepkg/adapters/api/v1/response"
	authService "github.com/tokyobordel/traineepkg/auth/service"
	jwtAuth "github.com/tokyobordel/traineepkg/authorization/jwt"
	"github.com/tokyobordel/traineepkg/errors"
)

type Handler struct {
	authService authService.IAuthService
	jwtService  *jwtAuth.Service
	logger      *log.Logger
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func NewHandler(authService authService.IAuthService, jwtService *jwtAuth.Service, accessTTL time.Duration, refreshTTL time.Duration) *Handler {
	return &Handler{
		authService: authService,
		jwtService:  jwtService,
		logger:      log.Default(),
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
	}
}

// Register godoc
//
//	@Summary		Регистрация пользователя
//	@Description	Создаёт нового пользователя по логину и паролю
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CredentialsRequest	true	"Учётные данные"
//	@Success		201		{object}	response.Envelope{data=AuthResponse}
//	@Failure		400		{object}	response.ErrorEnvelope
//	@Failure		409		{object}	response.ErrorEnvelope
//	@Failure		500		{object}	response.ErrorEnvelope
//	@Router			/auth/register [post]
func (h *Handler) Register(c fiber.Ctx) error {
	var req CredentialsRequest
	if err := c.Bind().Body(&req); err != nil || req.Login == "" || req.Pass == "" {
		response.MakeErrorResponse(c, h.logger, errors.NewInvalidParametersError("body", nil, "invalid request body"))
		return nil
	}

	user, err := h.authService.Register(req.Pass, req.Login)
	if err != nil {
		response.MakeErrorResponse(c, h.logger, err)
		return nil
	}

	response.MakeSuccessResponseWithStatus(c, fiber.StatusCreated, AuthResponse{User: user})
	return nil
}

// Login godoc
//
//	@Summary		Вход в систему
//	@Description	Аутентифицирует пользователя и устанавливает access/refresh токены в HTTP-only cookies
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CredentialsRequest	true	"Учётные данные"
//	@Success		200		{object}	response.Envelope{data=AuthResponse}
//	@Failure		400		{object}	response.ErrorEnvelope
//	@Failure		401		{object}	response.ErrorEnvelope
//	@Failure		500		{object}	response.ErrorEnvelope
//	@Router			/auth/login [post]
func (h *Handler) Login(c fiber.Ctx) error {
	var req CredentialsRequest
	if err := c.Bind().Body(&req); err != nil || req.Login == "" || req.Pass == "" {
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

	response.MakeSuccessResponse(c, AuthResponse{User: user})
	return nil
}

// Refresh godoc
//
//	@Summary		Обновление токенов
//	@Description	Обновляет пару access/refresh токенов по refresh_token cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Envelope{data=AuthResponse}
//	@Failure		401	{object}	response.ErrorEnvelope
//	@Failure		404	{object}	response.ErrorEnvelope
//	@Failure		500	{object}	response.ErrorEnvelope
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c fiber.Ctx) error {
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

	response.MakeSuccessResponse(c, AuthResponse{User: user})
	return nil
}

// GetMe godoc
//
//	@Summary		Текущий пользователь
//	@Description	Возвращает данные текущего пользователя по access_token cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	response.Envelope{data=AuthResponse}
//	@Failure		401	{object}	response.ErrorEnvelope
//	@Failure		404	{object}	response.ErrorEnvelope
//	@Router			/auth/me [get]
func (h *Handler) GetMe(c fiber.Ctx) error {
	accessToken := c.Cookies(jwtAuth.AccessTokenCookieName)
	if accessToken == "" {
		response.MakeErrorResponse(c, h.logger, errors.NewAuthTokenError(errors.TokenNotFound))
		return nil
	}

	userID, err := h.jwtService.ValidateAccessToken(accessToken)
	if err != nil {
		response.MakeErrorResponse(c, h.logger, errors.NewAuthTokenError(errors.InvalidToken))
		return nil
	}

	user, getMeErr := h.authService.GetMe(userID)
	if getMeErr != nil {
		response.MakeErrorResponse(c, h.logger, getMeErr)
		return nil
	}

	response.MakeSuccessResponse(c, AuthResponse{User: user})

	return nil
}

// Logout godoc
//
//	@Summary		Выход из системы
//	@Description	Удаляет access и refresh токены из cookies
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Envelope{data=response.MessageData}
//	@Router			/auth/logout [post]
func (h *Handler) Logout(c fiber.Ctx) error {
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

func (h *Handler) setTokenCookies(c fiber.Ctx, userID int) bool {
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
