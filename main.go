package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
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

const mode = "DEVELOPMENT"

func init() {
	os.Setenv("AYORADIO_MODE", mode)
	exec.Command("tailwindcss", "-i ./static/css/input.css", "-o ./static/css/style.css", "--watch").Run()
}

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

	// j, err = s.NewJob(
	// 	gocron.DurationJob(150 * time.Second),
	// 	gocron.NewTask(job.TurnOnRadio),
	// )
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Printf("JobRadio[ID]: %s\n", j.ID().String())

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

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			m.TextHTMLMiddleware,
			m.CSPMiddleware,
		)

		r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

		r.Get("/", handlers.NewHomeHandler().ServeHTTP)
		r.Get("/about", handlers.NewAboutHandler().ServeHTTP)
	})

	srv := &http.Server{
		Addr:    ":3000",
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
