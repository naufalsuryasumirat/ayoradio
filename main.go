package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"

	_ "github.com/naufalsuryasumirat/ayoradio/util"
	job "github.com/naufalsuryasumirat/ayoradio/jobs"
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

	j, err = s.NewJob(
		gocron.DurationJob(150 * time.Second),
		gocron.NewTask(job.TurnOnRadio),
	)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("JobRadio[ID]: %s\n", j.ID().String())

	j, err = s.NewJob(
        gocron.DurationJob(24 * time.Hour),
        gocron.NewTask(job.LoadBlacklistedDevices),
    )
	log.Printf("JobLoad[ID]: %s\n", j.ID().String())

	s.Start()

    done := make(chan os.Signal, 1)
    signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
    <-done

	s.Shutdown()
}
