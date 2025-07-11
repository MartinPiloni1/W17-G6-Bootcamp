package main

import (
	"fmt"
	config "github.com/aaguero_meli/W17-G6-Bootcamp/cmd/db"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	db := config.MustOpenDB()
	defer db.Close()

	migrationsDir := "docs/db/migrations"

	// -- Ensures the migrations table exists --
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		filename VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Error creating schema_migrations table: %v", err)
	}

	// -- Lists all migrations already executed --
	applied := map[string]struct{}{}
	rows, err := db.Query(`SELECT filename FROM schema_migrations`)
	if err != nil {
		log.Fatalf("Error reading schema_migrations: %v", err)
	}
	for rows.Next() {
		var fname string
		if err := rows.Scan(&fname); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		applied[fname] = struct{}{}
	}
	rows.Close()

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Could not read the migrations directory: %v", err)
	}

	// Filter and sort .sql files by name
	var migrations []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}
	sort.Strings(migrations)

	// Execute only those that are not registered
	var ran int
	for _, fname := range migrations {
		if _, ya := applied[fname]; ya {
			fmt.Printf("Skip %s (already applied)\n", fname)
			continue
		}
		fullPath := filepath.Join(migrationsDir, fname)
		fmt.Printf("Applying migration: %s\n", fname)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("Could not read %s: %v", fname, err)
		}
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Error executing %s: %v", fname, err)
		}
		_, err = db.Exec("INSERT INTO schema_migrations (filename) VALUES (?)", fname)
		if err != nil {
			log.Fatalf("Could not register migration %s: %v", fname, err)
		}
		ran++
	}

	if ran == 0 {
		fmt.Println("No new migrations to apply.")
	} else {
		fmt.Printf("Migrations completed! %d new applied.\n", ran)
	}
}
