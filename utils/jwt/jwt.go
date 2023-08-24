package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	TokenCookieKey  = "session_token"
	ExpiresDuration = 24 * time.Hour

	AuthClaimKey = "AuthorizationClaim"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
}

type Service struct {
	SecretKey []byte
}

func NewService(SecretKey []byte) *Service {
	return &Service{
		SecretKey: SecretKey,
	}
}

func (s *Service) GenerateToken(userID uint, userName string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpiresDuration)),
		},
		UserID:   userID,
		UserName: userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) VerifyToken(tokenString string) (*Claims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(s.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *Service) ForgetToken(c echo.Context) error {
	cookie, err := c.Cookie(TokenCookieKey)
	if err != nil {
		return err
	}

	cookie.Expires = time.Now().Add(-1 * time.Hour)
	c.SetCookie(cookie)
	return nil
}

func (s *Service) GenerateTokenAndStore(c echo.Context, userID uint, userName string) error {
	token, err := s.GenerateToken(userID, userName)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     userName,
		Value:    token,
		Expires:  time.Now().Add(ExpiresDuration),
		HttpOnly: true,
	}

	c.SetCookie(cookie)
	return nil
}

func (s *Service) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(TokenCookieKey)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			claims, err := s.VerifyToken(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			c.Set(AuthClaimKey, claims)
			return next(c)
		}
	}
}
