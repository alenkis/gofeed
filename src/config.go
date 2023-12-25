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
	Start           string `yaml:"start"`
	End             string `yaml:"end"`
	MongoUri        string `yaml:"mongoUri"`
	MongoCollection string `yaml:"mongoCollection"`
}

type ImportConfig struct {
	PostgresUri   string `yaml:"postgresUri"`
	PostgresTable string `yaml:"postgresTable"`
}

type JobConfig struct {
	Name         string `yaml:"name"`
	Schedule     int    `yaml:"schedule"`
	ScheduleUnit string `yaml:"scheduleUnit"`
}

// Validate returns an error if the config is invalid
func (c *Config) Validate() error {
	startTime, err := time.Parse(time.RFC3339, c.Export.Start)
	_ = startTime

	if err != nil {
		return fmt.Errorf("Error parsing start time: %v", err)
	}

	var endTime time.Time
	if c.Export.End == "" {
		scheduleDuration := time.Duration(c.Job.Schedule)
		var t time.Duration

		switch c.Job.ScheduleUnit {
		case "minute":
			t = time.Minute
		case "hour":
			t = time.Hour
		case "day":
			t = time.Hour * 24
		case "week":
			t = time.Hour * 24 * 7
		case "month":
			t = time.Hour * 24 * 30
		case "year":
			t = time.Hour * 24 * 365
		default:
			log.Fatalf("Invalid schedule unit: %s", c.Job.ScheduleUnit)
		}

		endTime = startTime.Add(scheduleDuration * t)

		// Update config with formatted end time
		c.Export.End = endTime.Format(time.RFC3339)
	} else {
		// If end time is specified, parse it
		endTime, err = time.Parse(time.RFC3339, c.Export.End)
		if err != nil {
			return fmt.Errorf("Error parsing end time: %v", err)
		}
	}

	_ = endTime

	return nil
}

// ParseConfig parses a YAML config file
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

	return &config, nil
}
