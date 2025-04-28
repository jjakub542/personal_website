package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"fullsteak/internal/article"
	"fullsteak/internal/contact"
	"fullsteak/internal/user"
)

type Services struct {
	User    user.Repository
	Article article.Repository
	Contact contact.Repository
}

type Server struct {
	port     int
	store    *user.SessionStore
	services *Services
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := PostgresClient()
	store := user.NewSessionStore()
	NewServer := &Server{
		port:  port,
		store: store,
		services: &Services{
			User:    user.NewRepository(db),
			Article: article.NewRepository(db),
			Contact: contact.NewRepository(db),
		},
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.Router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
