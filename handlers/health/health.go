package health

import (
	"github.com/misgorod/tucktuck-pull/common"
	"net/http"
)

type Handler struct{}

func New() Handler {
	return Handler{}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	common.RespondJSON(r.Context(), w, http.StatusOK, "Ok")
}
