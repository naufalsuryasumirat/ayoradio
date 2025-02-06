package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"

	"github.com/naufalsuryasumirat/ayoradio/internal/handlers"
	m "github.com/naufalsuryasumirat/ayoradio/internal/middleware"
	job "github.com/naufalsuryasumirat/ayoradio/jobs"
	_ "github.com/naufalsuryasumirat/ayoradio/util"
)

// run gocron to run the function for arp-scan
func main() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Panic(err)
	}

	j, err := s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(7, 0, 0))),
		gocron.NewTask(job.ResetSkipDay),
	)
	log.Printf("JobReset[ID]: %s\n", j.ID().String())

	is_prod := os.Getenv("AYORADIO_MODE") == "PRODUCTION"

	if is_prod {
		go func() {
			time.Sleep(5 * time.Second)
			job.TurnOnRadio()
		}()
		j, err = s.NewJob(
			gocron.DurationJob(150*time.Second),
			gocron.NewTask(job.TurnOnRadio),
		)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("JobRadio[ID]: %s\n", j.ID().String())

        j, err = s.NewJob(
            gocron.DurationJob(300*time.Second),
            gocron.NewTask(job.TurnOnDevice),
        )
        if err != nil {
            log.Panic(err)
        }
        log.Printf("JobWake[ID]: %s\n", j.ID().String())
	}

	j, err = s.NewJob(
		gocron.DurationJob(24*time.Hour),
		gocron.NewTask(job.LoadBlacklistedDevices),
	)
	log.Printf("JobLoad[ID]: %s\n", j.ID().String())

	s.Start()
	defer s.Shutdown()

	r := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
    r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./favicon.ico")
    })

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			m.TextHTMLMiddleware,
			m.CSPMiddleware,
		)

		r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

		r.Get("/", handlers.NewHomeHandler().ServeHTTP)
		r.Get("/about", handlers.NewAboutHandler().ServeHTTP)
		r.Get("/register", handlers.NewGetRegisterHandler().ServeHTTP)
		r.Post("/register", handlers.NewPostRegisterHandler().ServeHTTP)
		r.Get("/volume", handlers.NewControlsHandler().Volume)
		r.Post("/volume-up", handlers.NewControlsHandler().VolumeUp)
		r.Post("/volume-down", handlers.NewControlsHandler().VolumeDown)
		r.Post("/play", handlers.NewControlsHandler().Play)
		r.Post("/playlist", handlers.NewControlsHandler().Playlist)
		r.Get("/playing", handlers.NewControlsHandler().CurrentPlaying)
	})

	srv := &http.Server{
		Addr: func() string {
			if is_prod {
				return ":3000"
			} else {
				return ":3003"
			}
		}(),
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server shutdown complete")
		} else if err != nil {
			log.Panic(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Println("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Panic(err, "Server shutdown failed")
	}

	log.Println("Server shutdown complete")
}
