package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Export ExportConfig `yaml:"export"`
	Import ImportConfig `yaml:"import"`
	Job    JobConfig    `yaml:"job"`
}

type ExportConfig struct {
	MongoUri        string `yaml:"mongoUri"`
	MongoCollection string `yaml:"mongoCollection"`
}

type ImportConfig struct {
	PostgresUri   string `yaml:"postgresUri"`
	PostgresTable string `yaml:"postgresTable"`
}

type JobConfig struct {
	Name          string `yaml:"name"`
	Start         string `yaml:"start"`
	End           string
	ScheduleValue int    `yaml:"schedule"`
	ScheduleUnit  string `yaml:"scheduleUnit"`
}

// StartTime returns the RFC3339 formatted start time
func (c *Config) StartTime() time.Time {
	startTime, err := time.Parse(time.RFC3339, c.Job.Start)
	if err != nil {
		log.Fatalf("Error parsing start time: %v", err)
	}

	return startTime
}

func (config *Config) calculateTimeRange(state *SchedulerState) (string, string) {
	state.Lock()
	defer state.Unlock()

	duration := config.Duration()

	// Calculate the current start and end times
	currentStartTime := state.initialTime.Add(duration * time.Duration(state.elapsedCycles))
	currentEndTime := currentStartTime.Add(duration)

	// Increment the number of elapsed cycles for the next calculation
	state.elapsedCycles++

	return currentStartTime.Format(time.RFC3339), currentEndTime.Format(time.RFC3339)
}

// Validate returns an error if the config is invalid
func (c *Config) Validate() error {
	startTime := c.StartTime()

	var endTime time.Time
	configDuration := c.Duration()

	endTime = startTime.Add(configDuration)

	// Update config with formatted end time
	c.Job.End = endTime.Format(time.RFC3339)

	return nil
}

// ParseConfig parses a YAML config file, and then
// validates it
func ParseConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Validate() != nil {
		log.Fatalf("Error validating config file: %v", err)
	}

	return &config, nil
}

func (c *Config) Duration() time.Duration {
	durationString := fmt.Sprintf("%d%s", c.Job.ScheduleValue, c.Job.ScheduleUnit)
	configDuration, err := time.ParseDuration(durationString)

	if err != nil {
		log.Fatalf("Error parsing schedule duration: %v", err)
	}

	return configDuration

}
