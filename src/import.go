package main

import (
	"fmt"
	"os/exec"
)

func psqlImportQuery(c *Config) string {
	cmd := fmt.Sprintf(`"\copy public.products_raw FROM 'products.out.json' WITH (FORMAT csv, QUOTE E'\x01', DELIMITER E'\x02')"`)
	return fmt.Sprintf(`psql %s -c %s`, c.Import.PostgresUri, cmd)
}

func PostgresImport(c *Config) error {
	q := psqlImportQuery(c)

	fmt.Printf("Executing psql: %s\n", q)

	cmd := exec.Command("bash", "-c", q)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("psql import failed: %v, output: %s", err, string(output))
	}

	fmt.Println("Data exported and copied successfully")
	return nil
}
