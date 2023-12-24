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

	// Construct the mongoexport mongoQuery
	mongoQuery := fmt.Sprintf(`{"updatedAt": {"$gte": {"$date": "%s"}, "$lte": {"$date": "%s"}}}`, *startTime, *endTime)
	mongoUri := "mongodb://mongoadmin:secret@localhost:27017/mydatabase?authSource=admin"
	mongoCollection := "products"
	exportedProducts := "products.out.json"

	// Execute mongoexport command
	mongoExportCmd := fmt.Sprintf("mongoexport --uri='%s' --collection='%s' --query='%s' --out='%s'", mongoUri, mongoCollection, mongoQuery, exportedProducts)
	cmd := exec.Command("bash", "-c", mongoExportCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing mongoexport: %v\nOutput: %s\n", err, string(output))
		return
	}

	postgresUri := "postgresql://postgres:secret@localhost:5432/products"
	psqlCmd := `psql ` + postgresUri + ` -c "\copy public.products_raw FROM 'products.out.json' WITH (FORMAT csv, QUOTE E'\x01', DELIMITER E'\x02')"`

	cmd = exec.Command("bash", "-c", psqlCmd)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing psql: %v\nOutput: %s\n", err, string(output))
		return
	}

	fmt.Println("Data exported and copied successfully")
}
