package public

import (
	"net/http"
	"strconv"

	"fullsteak/internal/article"
	"fullsteak/internal/contact"
	"fullsteak/internal/user"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	User    user.Repository
	Article article.Repository
	Contact contact.Repository
}

func (h *Handler) HomePage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"PageTitle": "Home",
		"CSRFToken": csrfToken,
	})
}

func (h *Handler) PortfolioPage(c echo.Context) error {
	return c.Render(http.StatusOK, "portfolio.html", nil)
}

func (h *Handler) BlogPage(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit := 8
	offset := (page - 1) * limit
	articles, err := h.Article.GetAllPublicBetween(limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	count, err := h.Article.GetCount()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	pagesCount := count/limit + 1
	return c.Render(http.StatusOK, "blog.html", map[string]interface{}{
		"PageTitle":     "Blog",
		"Articles":      articles,
		"articlesCount": count,
		"pagesCount":    pagesCount,
		"page":          page,
		"previous":      page - 1,
		"next":          page + 1,
	})
}

func (h *Handler) ArticleView(c echo.Context) error {
	id := c.Param("article_id")
	a, err := h.Article.GetOneById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.Render(http.StatusOK, "article.html", map[string]interface{}{
		"PageTitle": a.Title,
		"Article":   a,
	})
}

func (h *Handler) ContactUsPost(c echo.Context) error {
	msg := &contact.MessageForm{
		Name:    c.FormValue("name"),
		Email:   c.FormValue("email"),
		Message: c.FormValue("message"),
	}
	err := h.Contact.CreateOne(msg)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.Redirect(http.StatusSeeOther, "/#contact")
}

func (h *Handler) SimulationPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	return c.Render(http.StatusOK, "simulation.html", map[string]interface{}{
		"PageTitle": "Simulation",
		"CSRFToken": csrfToken,
	})
}
