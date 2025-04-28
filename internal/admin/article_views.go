package admin

import (
	"fullsteak/internal/article"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ArticlesPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	articles, err := h.Article.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Render(http.StatusOK, "admin/articles.html", map[string]interface{}{
		"Articles":  articles,
		"CSRFToken": csrfToken,
	})
}

func (h *Handler) ArticleCreate(c echo.Context) error {
	a := &article.Article{
		Title:       c.FormValue("title"),
		Description: c.FormValue("desc"),
		Public:      false,
	}
	err := h.Article.CreateOne(a)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Redirect(http.StatusSeeOther, "/admin/articles")
}

func (h *Handler) ArticleDelete(c echo.Context) error {
	images, err1 := h.Article.GetArticleImages(c.Param("article_id"))
	if err1 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	for _, image := range images {
		err := image.Remove()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}
	err2 := h.Article.DeleteOneById(c.Param("article_id"))
	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Redirect(http.StatusSeeOther, "/admin/articles")
}

func (h *Handler) ArticleUpdate(c echo.Context) error {
	a := &article.Article{
		Title:       c.FormValue("title"),
		Description: c.FormValue("desc"),
		Content:     c.FormValue("content"),
	}
	if c.FormValue("public") == "on" {
		a.Public = true
	} else {
		a.Public = false
	}
	err := h.Article.UpdateOneById(a, c.Param("article_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// Check if there is a new cover image uploaded
	coverImage, _ := c.FormFile("cover_image")

	if coverImage != nil {
		src, err := coverImage.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
		}
		defer src.Close()
		newImageID, err := h.Article.AttachImage(c.Param("article_id"))
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error saving image metadata: "+err.Error())
		}
		i := &article.Image{Id: newImageID}
		err = i.Save(src)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
		}

		err = h.Article.UpdateArticleCoverImage(c.Param("article_id"), newImageID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error updating article's cover image: "+err.Error())
		}
	}

	return c.Redirect(http.StatusSeeOther, "/admin/articles")
}

func (h *Handler) ArticleAttachImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}
	defer src.Close()
	newImageID, err := h.Article.AttachImage(c.Param("article_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	i := &article.Image{Id: newImageID}
	err = i.Save(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}
	return c.Redirect(http.StatusSeeOther, "/admin/articles/"+c.Param("article_id")+"/edit")
}

func (h *Handler) ArticleDeleteImage(c echo.Context) error {
	i := &article.Image{Id: c.QueryParam("id")}
	err := i.Remove()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	err = h.Article.RemoveImage(i.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusOK, "deleted")
}

func (h *Handler) ArticleEditPage(c echo.Context) error {
	csrfToken := c.Get("csrf").(string)
	a, err := h.Article.GetOneById(c.Param("article_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	images, err := h.Article.GetArticleImages(c.Param("article_id"))
	if err != nil {
		log.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	a.Images = append(a.Images, images...)
	return c.Render(http.StatusOK, "admin/article.html", map[string]interface{}{
		"Article":   a,
		"CSRFToken": csrfToken,
	})
}
