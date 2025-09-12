package token

import (
	"fmt"
	"time"

	"github.com/Bromolima/my-game-list/config"
	"github.com/Bromolima/my-game-list/internal/http/cookie"
	resterr "github.com/Bromolima/my-game-list/internal/http/rest_err"
	"github.com/Bromolima/my-game-list/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtService interface {
	GenerateToken(*models.User) (string, error)
	VerifyToken(echo.Context) error
	ExtractToken(echo.Context) (*models.User, error)
}

type jwtService struct {
	SecretKey string
}

func NewJwtService() JwtService {
	return &jwtService{
		SecretKey: config.Env.SecretKey,
	}
}

func (s *jwtService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *jwtService) VerifyToken(ectx echo.Context) error {
	cookie, err := ectx.Cookie(cookie.CookieName)
	if err != nil {
		return err
	}

	tokenValue := cookie.Value

	token, err := jwt.Parse(tokenValue, getIdentificationKey)
	if err != nil {
		return nil
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return resterr.ErrInvalidToken
	}

	return nil
}

func (s *jwtService) ExtractToken(ctx echo.Context) (*models.User, error) {
	cookie, err := ctx.Cookie(cookie.CookieName)
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

	return &models.User{
		ID:       claims["id"].(string),
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
	}, nil
}

func getIdentificationKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("signing method unexpected! %v", token.Header["alg"])
	}

	return []byte(config.Env.SecretKey), nil
}
