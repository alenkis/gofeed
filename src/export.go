package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

func mongoExportQuery(config *Config, start string, end string) string {
	mongoQuery := fmt.Sprintf(`{"updatedAt": {"$gte": {"$date": "%s"}, "$lte": {"$date": "%s"}}}`, start, end)
	exportedProducts := "products.out.json"

	return fmt.Sprintf("mongoexport --uri='%s' --collection='%s' --query='%s' --out='%s'", config.Export.MongoUri, config.Export.MongoCollection, mongoQuery, exportedProducts)
}

func MongoExport(config *Config, start string, end string) error {
	q := mongoExportQuery(config, start, end)

	log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "mongo-export",
	}).
		With("query", q).
		Info("Executing mongoexport query\n")

	cmd := exec.Command("bash", "-c", q)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mongoexport failed: %v, output: %s", err, string(output))
	}

	// Get the size of the exported file
	fileInfo, err := os.Stat("products.out.json")

	// Print file size in bytes
	log.Infof("Exported file size: %d bytes\n", fileInfo.Size())

	return nil
}
