package token

import (
	"fmt"
	"time"

	"github.com/Bromolima/my-game-list/config"
	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/factory"
	"github.com/Bromolima/my-game-list/internal/http/cookie"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

//go:generate mockgen -source=token.go -destination=../../mocks/token_service.go -package=mocks
type JwtService interface {
	ValidateToken(echo.Context) error
	GenerateToken(*entities.User) (string, error)
	ExtractToken(ectx echo.Context) (*entities.User, error)
}

type jwtService struct {
	SecretKey string
}

func NewJwtService() JwtService {
	return &jwtService{
		SecretKey: config.Env.SecretKey,
	}
}

func (s *jwtService) GenerateToken(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"id":      user.ID,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *jwtService) ValidateToken(ectx echo.Context) error {
	cookie, err := ectx.Cookie(cookie.CookieName)
	if err != nil {
		return err
	}

	tokenValue := cookie.Value
	token, err := jwt.Parse(tokenValue, getIdentificationKey)
	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return err
	}

	return nil
}

func (s *jwtService) ExtractToken(ectx echo.Context) (*entities.User, error) {
	cookie, err := ectx.Cookie(cookie.CookieName)
	if err != nil {
		return nil, err
	}

	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, getIdentificationKey)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	userClaims := factory.NewUserClaims(
		claims["id"].(uuid.UUID),
		claims["role_id"].(uint),
	)

	return userClaims, nil
}

func getIdentificationKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("signing method unexpected! %v", token.Header["alg"])
	}

	return []byte(config.Env.SecretKey), nil
}
