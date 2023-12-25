package main

import (
	"flag"
	"log"
	"sync"
	"time"
)

type SchedulerState struct {
	sync.Mutex
	elapsedCycles int
	initialTime   time.Time
}

func NewSchedulerState(initialTime time.Time) *SchedulerState {
	return &SchedulerState{
		initialTime: initialTime,
	}
}

func main() {
	configPath := flag.String("config", "config.yml", "path to configuration file")
	flag.Parse()

	config, err := ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	initialStartTime := config.StartTime()

	// Create a new scheduler state
	state := NewSchedulerState(initialStartTime)
	duration := config.Duration()
	ticker := time.NewTicker(duration)

	log.Printf("Starting %s, with export-import cycle every %d %s", config.Job.Name, config.Job.ScheduleValue, config.Job.ScheduleUnit)

	// Scheduler loop
	for {
		select {
		case <-ticker.C:
			startTime, endTime := config.calculateTimeRange(state)
			go handleExportImport(config, startTime, endTime)
		}
	}
}

func handleExportImport(config *Config, startTime string, endTime string) {
	log.Println("Starting the job")
	// Implement retry logic if needed
	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		if err := MongoExport(config, startTime, endTime); err != nil {
			log.Printf("MongoExport failed: %v, retrying...", err)
			continue
		}

		if err := PostgresImport(config); err != nil {
			log.Printf("PostgresImport failed: %v, retrying...", err)
			continue
		}

		log.Println("Export-Import cycle completed successfully")
		return
	}

	log.Printf("Export-Import cycle failed after %d retries", maxRetries)
}
