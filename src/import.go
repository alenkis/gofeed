package main

import (
	"fmt"
	"os/exec"
)

func psqlImportQuery(c *Config) string {
	cmd := fmt.Sprintf(`"\copy public.products_raw FROM 'products.out.json' WITH (FORMAT csv, QUOTE E'\x01', DELIMITER E'\x02')"`)
	return fmt.Sprintf(`psql %s -c %s`, c.Import.PostgresUri, cmd)
}

func PostgresImport(c *Config) {
	q := psqlImportQuery(c)

	fmt.Printf("Executing psql: %s\n", q)

	cmd := exec.Command("bash", "-c", q)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing psql: %v\nOutput: %s\n", err, string(output))
		return
	}

	fmt.Println("Data exported and copied successfully")
}
