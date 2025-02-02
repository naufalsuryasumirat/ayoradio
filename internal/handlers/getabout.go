package handlers

import (
	"github.com/naufalsuryasumirat/ayoradio/internal/templates"
	"net/http"
)

type AboutHandler struct{}

func NewAboutHandler() *AboutHandler {
	return &AboutHandler{}
}

func (h *AboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c := templates.About()
	err := templates.Layout(c, "ayoradio").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
