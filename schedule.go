package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func InitSchedule() {
	s := gocron.NewScheduler(time.UTC)

	// s.Every(5).Seconds().Do(func(){ ... })

	// s.Every("5m").Do(func(){ ... })

	// s.Every(5).Days().Do(func(){ ... })

	s.Cron("10 16 * * 1-5").Do(func() {

	}) // every minute

	s.StartAsync()
}

// func setTasks(s) {

// }
