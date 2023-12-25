package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

func psqlImportQuery(c *Config) string {
	cmd := fmt.Sprintf(`"\copy public.products_raw FROM 'products.out.json' WITH (FORMAT csv, QUOTE E'\x01', DELIMITER E'\x02')"`)
	return fmt.Sprintf(`psql %s -c %s`, c.Import.PostgresUri, cmd)
}

func PostgresImport(c *Config) error {
	q := psqlImportQuery(c)

	log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "postgres-import",
	}).
		With("query", q).
		Info("Executing psql copy command\n")

	cmd := exec.Command("bash", "-c", q)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("psql import failed: %v, output: %s", err, string(output))
	}

	log.Info("Data exported and copied successfully")
	return nil
}
