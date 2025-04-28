package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID := c.Get("sessionID").(string)
		store := c.Get("sessionStore").(*SessionStore)
		role, ok := store.Get(sessionID, "role")
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden: Admin only"})
		}
		return next(c)
	}
}
