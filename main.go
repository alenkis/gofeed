package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	// Define flags for the time range
	startTime := flag.String("start", "", "Start time for timerange")
	endTime := flag.String("end", "", "End time for timerange")
	flag.Parse()

	// Validate time range
	if *startTime == "" || *endTime == "" {
		fmt.Println("Please specify both start and end time")
		return
	}

	// Construct the query
	query := fmt.Sprintf(`{"timerange": {"$gte": "%s", "$lte": "%s"}}`, *startTime, *endTime)

	// Execute mongoexport command
	mongoExportCmd := fmt.Sprintf("mongoexport --uri=yourMongoDBUri --collection=products --query='%s' --out=products.json", query)
	err := exec.Command("bash", "-c", mongoExportCmd).Run()
	if err != nil {
		fmt.Printf("Error executing mongoexport: %v\n", err)
		return
	}

	// Execute psql command
	psqlCmd := "psql yourPostgresDBUri -c \"\\COPY products FROM 'products.json' WITH (FORMAT json)\""
	err = exec.Command("bash", "-c", psqlCmd).Run()
	if err != nil {
		fmt.Printf("Error executing psql: %v\n", err)
		return
	}

	fmt.Println("Data exported and copied successfully")
}
