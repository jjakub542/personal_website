package contact

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Contact Repository
}

func (h *Handler) MessageDelete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing message ID")
	}
	err := h.Contact.DeleteOneById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete message")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Message deleted",
	})
}

func (h *Handler) MessageCreate(c echo.Context) error {
	var m MessageForm

	if err := c.Bind(&m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	err := h.Contact.CreateOne(&m)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create message")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "success",
	})
}
