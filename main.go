package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/misgorod/tucktuck-pull/handlers/health"
	"github.com/misgorod/tucktuck-pull/handlers/pull"
	"log"
	"net/http"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	pullHandler := pull.New(ctx)
	healthHandler := health.New()
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)
	r.Route("/api", func(r chi.Router) {
		r.Get("/healthcheck", healthHandler.Get)
		r.Get("/pull", pullHandler.Get)
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}
