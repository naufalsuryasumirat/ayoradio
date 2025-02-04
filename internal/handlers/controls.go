package handlers

import (
	"fmt"
	"net/http"

	"github.com/naufalsuryasumirat/ayoradio/internal/templates"
	"github.com/naufalsuryasumirat/ayoradio/jobs"
)

type ControlsHandler struct {}

func NewControlsHandler() *ControlsHandler {
	return &ControlsHandler{}
}

func (h *ControlsHandler) VolumeUp(w http.ResponseWriter, r *http.Request) {
    jobs.VolumeIncrease()
}

func (h *ControlsHandler) VolumeDown(w http.ResponseWriter, r *http.Request) {
    jobs.VolumeDecrease()
}

func (h *ControlsHandler) Volume(w http.ResponseWriter, r *http.Request) {
    vol := jobs.Volume()
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, "%.2f%%", vol)
}

func (h *ControlsHandler) Play(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("addr")
    playType := r.Form.Get("play")
    if playType == "queue" {
        jobs.QueueLink(addr)
    } else { // defaults replace
        jobs.PlayLink(addr)
    }
}

func (h *ControlsHandler) Playlist(w http.ResponseWriter, r *http.Request) {
    itype := r.Form.Get("playlist")
    if itype == "prev" {
        jobs.PlayPrev()
    } else { // defaults next
        jobs.PlayNext()
    }
}

// could also just change the src attr, but requires more effort (better perf?)
func (h * ControlsHandler) CurrentPlaying(w http.ResponseWriter, r *http.Request) {
    addr := jobs.CurrentPlaying()
	c := templates.CurPlaying(addr)
    err := c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}

