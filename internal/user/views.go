package user

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	User Repository
}

func (h *Handler) LoginPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	sessionID := c.Get("sessionID").(string)
	store := c.Get("sessionStore").(*SessionStore)

	if c.Request().Method == http.MethodGet {
		role, ok := store.Get(sessionID, "role")
		if ok || role == "admin" {
			return c.Redirect(http.StatusSeeOther, "/admin")
		}
		return c.Render(http.StatusOK, "admin/login.html", map[string]interface{}{
			"CSRFToken": csrfToken,
		})
	}

	requestUser := User{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	requestUser.CreatePasswordHash()

	u, err := h.User.GetOneByEmail(c.FormValue("email"))

	fmt.Println(u.Email)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User with this email does not exist")
	}

	if u.PasswordHash != requestUser.PasswordHash {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid password")
	}

	if !u.IsSuperuser {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	store.Set(sessionID, "authenticated", true)
	store.Set(sessionID, "role", "admin")

	return c.Redirect(http.StatusSeeOther, "/admin")
}

func (h *Handler) LogoutPage(c echo.Context) error {
	sessionID := c.Get("sessionID").(string)
	store := c.Get("sessionStore").(*SessionStore)

	store.Delete(sessionID)
	return c.Render(http.StatusOK, "admin/logout.html", nil)
}
