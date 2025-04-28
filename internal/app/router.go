package app

import (
	"fullsteak/internal/admin"
	"fullsteak/internal/contact"
	"fullsteak/internal/public"
	"fullsteak/internal/user"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) Router() http.Handler {
	e := echo.New()
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:  "header:X-CSRF-Token,form:_csrf",
		CookieName:   "_csrf", // name of the cookie that stores it
		CookieSecure: true,    // set to true if using HTTPS
	}))
	e.Use(user.SessionMiddleware(s.store))
	e.Renderer = Renderer()
	e.Static("/static", "web/static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	publicHandler := public.Handler{User: s.services.User, Article: s.services.Article, Contact: s.services.Contact}
	adminHandler := admin.Handler{User: s.services.User, Article: s.services.Article, Contact: s.services.Contact}
	userHandler := user.Handler{User: s.services.User}
	contactHandler := contact.Handler{Contact: s.services.Contact}

	e.GET("/", publicHandler.HomePage)
	e.GET("/portfolio", publicHandler.PortfolioPage)
	e.GET("/blog", publicHandler.BlogPage)
	e.GET("/blog/:article_id", publicHandler.ArticleView)
	e.POST("/send-message", publicHandler.ContactUsPost)
	e.POST("/messages", contactHandler.MessageCreate)

	adminGroup := e.Group("/admin")
	adminGroup.GET("", user.AdminAuth(adminHandler.HomePage))
	adminGroup.GET("/messages", user.AdminAuth(adminHandler.MessagesPage))
	adminGroup.DELETE("/messages/:id", user.AdminAuth(contactHandler.MessageDelete))
	adminGroup.GET("/articles", user.AdminAuth(adminHandler.ArticlesPage))
	adminGroup.POST("/articles/create", user.AdminAuth(adminHandler.ArticleCreate))
	adminGroup.POST("/articles/:article_id/delete", user.AdminAuth(adminHandler.ArticleDelete))
	adminGroup.POST("/articles/:article_id/update", user.AdminAuth(adminHandler.ArticleUpdate))
	adminGroup.POST("/articles/:article_id/attach-image", user.AdminAuth(adminHandler.ArticleAttachImage))
	adminGroup.POST("/articles/delete-image", user.AdminAuth(adminHandler.ArticleDeleteImage))
	adminGroup.GET("/articles/:article_id/edit", user.AdminAuth(adminHandler.ArticleEditPage))
	adminGroup.Any("/login", userHandler.LoginPage)
	adminGroup.Any("/logout", userHandler.LogoutPage)

	return e
}
