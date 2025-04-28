package admin

import (
	"fullsteak/internal/article"
	"fullsteak/internal/contact"
	"fullsteak/internal/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	User    user.Repository
	Article article.Repository
	Contact contact.Repository
}

func (h *Handler) HomePage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/home.html", nil)
}

func (h *Handler) MessagesPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	messages, err := h.Contact.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Render(http.StatusOK, "admin/messages.html", map[string]any{
		"Count":     len(messages),
		"Messages":  messages,
		"CSRFToken": csrfToken,
	})
}

func (h *Handler) MessagesDelete(c echo.Context) error {
	id := c.FormValue("message_id")
	err := h.Contact.DeleteOneById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Redirect(http.StatusSeeOther, "/admin/messages")
}
