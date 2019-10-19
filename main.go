package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/misgorod/tucktuck-pull/common"
	"github.com/misgorod/tucktuck-pull/handlers/health"
	"github.com/misgorod/tucktuck-pull/handlers/pull"
	"github.com/misgorod/tucktuck-pull/repository"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	client, err := repository.New()
	if err != nil {
		log.WithError(err).Fatal("couldn't connect to database")
	}
	pullHandler := pull.New(client)
	healthHandler := health.New()
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)
	r.Route("/api", func(r chi.Router) {
		r.Get("/healthcheck", healthHandler.Get)
		r.Get("/pull", pullHandler.Get)
	})
	port := common.GetEnv("PORT", "8080")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
