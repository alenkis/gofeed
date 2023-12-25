package main

import (
	"fmt"
	"os/exec"
)

func mongoExportQuery(start string, end string) string {
	mongoQuery := fmt.Sprintf(`{"updatedAt": {"$gte": {"$date": "%s"}, "$lte": {"$date": "%s"}}}`, start, end)
	mongoUri := "mongodb://mongoadmin:secret@localhost:27017/mydatabase?authSource=admin"
	mongoCollection := "products"
	exportedProducts := "products.out.json"

	return fmt.Sprintf("mongoexport --uri='%s' --collection='%s' --query='%s' --out='%s'", mongoUri, mongoCollection, mongoQuery, exportedProducts)
}

func MongoExport(config *Config, start string, end string) error {
	q := mongoExportQuery(start, end)

	fmt.Printf("Executing mongoexport query: %s\n", q)

	cmd := exec.Command("bash", "-c", q)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mongoexport failed: %v, output: %s", err, string(output))
	}

	return nil
}
