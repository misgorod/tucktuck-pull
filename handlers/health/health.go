package health

import (
	"net/http"
)

type HealthHandler struct{}

func New() HealthHandler {
	return HealthHandler{}
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
