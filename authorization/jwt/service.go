package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims описывает полезную нагрузку JWT-токена пользователя.
type Claims struct {
	// UserID — идентификатор пользователя.
	UserID int `json:"user_id"`
	// TokenType — тип токена: access или refresh.
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenPair содержит пару access- и refresh-токенов.
type TokenPair struct {
	// AccessToken — короткоживущий токен доступа.
	AccessToken string
	// RefreshToken — долгоживущий токен обновления.
	RefreshToken string
}

// Service инкапсулирует создание и проверку JWT-токенов.
type Service struct {
	secret                 []byte
	accessTTL              time.Duration
	refreshTTL             time.Duration
	accessTokenCookieName  string
	refreshTokenCookieName string
}

// NewService создаёт сервис JWT с заданным секретом и временем жизни токенов.
func NewService(secret string, accessTTL time.Duration, refreshTTL time.Duration) *Service {
	return &Service{
		secret:                 []byte(secret),
		accessTTL:              accessTTL,
		refreshTTL:             refreshTTL,
		accessTokenCookieName:  AccessTokenCookieNameDefault,
		refreshTokenCookieName: RefreshTokenCookieNameDefault,
	}
}

func NewServiceWithCookieName(secret string, accessTTL time.Duration, refreshTTL time.Duration,
	accessTokenCookieName string, refreshTokenCookieName string) *Service {
	return &Service{
		secret:                 []byte(secret),
		accessTTL:              accessTTL,
		refreshTTL:             refreshTTL,
		accessTokenCookieName:  accessTokenCookieName,
		refreshTokenCookieName: refreshTokenCookieName,
	}
}

// GenerateTokenPair создаёт access- и refresh-токены для пользователя userID.
func (s *Service) GenerateTokenPair(userID int) (TokenPair, error) {
	accessToken, err := s.GenerateAccess(userID)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := s.GenerateRefresh(userID)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateAccess создаёт access-токен для пользователя userID.
func (s *Service) GenerateAccess(userID int) (string, error) {
	accessToken, err := s.generateToken(userID, AccessTokenType, s.accessTTL)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// GenerateRefresh создаёт refresh-токен для пользователя userID.
func (s *Service) GenerateRefresh(userID int) (string, error) {
	accessToken, err := s.generateToken(userID, RefreshTokenType, s.refreshTTL)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// GetAccessTTL возвращает время жизни access-токена.
func (s *Service) GetAccessTTL() time.Duration {
	return s.accessTTL
}

// GetRefreshTTL возвращает время жизни refresh-токена.
func (s *Service) GetRefreshTTL() time.Duration {
	return s.refreshTTL
}

func (s *Service) GetAccessTokenCookieName() string {
	return s.accessTokenCookieName
}

func (s *Service) GetRefreshTokenCookieName() string {
	return s.refreshTokenCookieName
}

// GetSecret возвращает секрет подписи JWT.
func (s *Service) GetSecret() []byte {
	return s.secret
}

// ValidateAccessToken проверяет access-токен и возвращает идентификатор пользователя.
func (s *Service) ValidateAccessToken(token string) (int, error) {
	return s.validateByType(token, AccessTokenType)
}

// ValidateRefreshToken проверяет refresh-токен и возвращает идентификатор пользователя.
func (s *Service) ValidateRefreshToken(token string) (int, error) {
	return s.validateByType(token, RefreshTokenType)
}

func (s *Service) generateToken(userID int, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) validateByType(token string, expectedType string) (int, error) {
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return 0, err
	}
	if !parsedToken.Valid {
		return 0, errors.New("invalid jwt token")
	}
	if claims.TokenType != expectedType {
		return 0, errors.New("unexpected token type")
	}
	return claims.UserID, nil
}
