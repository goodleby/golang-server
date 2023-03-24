package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/pure-go-server/client/database"
	"github.com/goodleby/pure-go-server/config"
	"github.com/goodleby/pure-go-server/server/handler"
)

const v1API string = "/api/v1"

type Server struct {
	Config *config.Config
	Router *chi.Mux
	HTTP   *http.Server
	DB     *database.Client
}

func New(ctx context.Context, config *config.Config) (*Server, error) {
	var s Server

	s.Config = config
	s.Router = chi.NewRouter()
	s.HTTP = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Config.Port),
		Handler: s.Router,
	}

	dbClient, err := database.New()
	if err != nil {
		return &s, fmt.Errorf("error creating database client: %w", err)
	}
	s.DB = dbClient

	s.setupRoutes()

	return &s, nil
}

func (s *Server) Start() error {
	log.Printf("Server is listening on port :%d", s.Config.Port)
	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("unexpected server error: %w", err)
	}

	return nil
}

func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Healthz)

	s.Router.Route(v1API, func(r chi.Router) {
		r.Get("/articles", handler.GetArticles(s.DB))
		r.Post("/articles", handler.CreateArticle(s.DB))

		r.Get("/articles/{id}", handler.GetArticle(s.DB))
		r.Delete("/articles/{id}", handler.DeleteArticle(s.DB))
		r.Patch("/articles/{id}", handler.UpdateArticle(s.DB))
	})
}