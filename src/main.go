package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	configPath := flag.String("config", "config.yml", "path to configuration file")
	flag.Parse()

	config, err := ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	if config.Validate() != nil {
		log.Fatalf("Error validating config file: %v", err)
	}

	fmt.Printf("Start time: %s\nEnd time: %s\n", config.Export.Start, config.Export.End)

	MongoExport(config)
	PostgresImport(config)
}
