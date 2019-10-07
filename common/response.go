package common

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RespondJSON(ctx context.Context, w http.ResponseWriter, status int, payload interface{}) {
	logger := log.WithFields(
		log.Fields{
			"request_id": middleware.GetReqID(ctx),
			"status":     status,
			"payload":    payload,
		},
	)
	if status < 500 {
		logger.Info()
	} else {
		logger.Error()
	}
	response, err := json.Marshal(payload)
	if err != nil {
		logger.WithError(err).Error("error while marshalling response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write([]byte(response))
	if err != nil {
		logger.WithError(err).Error("error while writing response")
		return
	}
}

func RespondError(ctx context.Context, w http.ResponseWriter, status int, err error) {
	RespondJSON(ctx, w, status, map[string]string{"error": err.Error()})
}
