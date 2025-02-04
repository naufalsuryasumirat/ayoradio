package handlers

import (
	"net/http"

	"github.com/naufalsuryasumirat/ayoradio/internal/templates"
	"github.com/naufalsuryasumirat/ayoradio/util"
)

type PostRegisterHandler struct {}

func NewPostRegisterHandler() *PostRegisterHandler {
	return &PostRegisterHandler{}
}

// FIXME: borked
func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mac := r.FormValue("mac")

    rtype := r.Form.Get("type")
    whitelist := rtype == "whitelist"
    err := util.TryAddDevice(mac, !whitelist)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	c := templates.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
