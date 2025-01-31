package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"

	job "github.com/naufalsuryasumirat/ayoradio/jobs"
	_ "github.com/naufalsuryasumirat/ayoradio/util"
)

// run gocron to run the function for arp-scan
func main() {
	locals := job.ScanLocalDevices()
	for _, l := range locals {
		fmt.Printf("local{whitelisted}: %s", l)
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Panic(err.Error())
	}

	j, err := s.NewJob(
		gocron.DurationJob(5*time.Minute),
		gocron.NewTask(job.ScanLocalDevices),
	)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Printf("JobScan[ID]: %s\n", j.ID().String())

	j, err = s.NewJob(
        gocron.DurationJob(24*time.Hour),
        gocron.NewTask(job.LoadBlacklistedDevices),
    )
	log.Printf("JobLoad[ID]: %s\n", j.ID().String())

	s.Start()

    select {
    case <-time.After(time.Minute):
    }

	s.Shutdown()
}
