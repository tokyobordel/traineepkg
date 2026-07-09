package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    int    `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(secret string, accessTTL time.Duration, refreshTTL time.Duration) *Service {
	return &Service{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

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

func (s *Service) GenerateAccess(userID int) (string, error) {
	accessToken, err := s.generateToken(userID, AccessTokenType, s.accessTTL)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *Service) GenerateRefresh(userID int) (string, error) {
	accessToken, err := s.generateToken(userID, RefreshTokenType, s.refreshTTL)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s *Service) GetAccessTTL() time.Duration {
	return s.accessTTL
}

func (s *Service) GetRefreshTTL() time.Duration {
	return s.refreshTTL
}

func (s *Service) GetSecret() []byte {
	return s.secret
}

func (s *Service) ValidateAccessToken(token string) (int, error) {
	return s.validateByType(token, AccessTokenType)
}

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
