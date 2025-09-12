package cookie

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	CookieName = "my_game_list_id"
)

func GetCookie(ectx echo.Context) (*http.Cookie, error) {
	cookie, err := ectx.Cookie(CookieName)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func SetCookie(ectx echo.Context, value string) {
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}

	ectx.SetCookie(cookie)
}

func DeleteCookie(ectx echo.Context) {
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	ectx.SetCookie(cookie)
}
