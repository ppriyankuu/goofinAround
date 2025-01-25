package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

func scheduleTaskProduction() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(10).Second().Do(func() {
		log.Println("Scheduler triggered: producing tasks...")
		produceTasks()
	})
	s.StartBlocking()
}
