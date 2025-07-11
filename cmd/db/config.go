package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Loads .env and returns the configuration data
func LoadEnv() {
	_ = godotenv.Load()
}

// Gets the DSN and config for MySQL
func GetDBConfigFromEnv() mysql.Config {
	return mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOST"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
		MultiStatements:      true,
	}
}

// Attempts to open the database and if it does not exist, creates it and tries again.
func MustOpenDB() *sql.DB {
	LoadEnv()
	cfg := GetDBConfigFromEnv()
	for {
		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Fatalf("Could not open DB: %v", err)
		}

		// Attempts the real ping, possible error of "database does not exist"
		pingErr := db.Ping()
		if pingErr == nil {
			// OK!
			return db
		}

		// Detects "unknown database" error
		if merr, ok := pingErr.(*mysql.MySQLError); ok && merr.Number == 1049 {
			log.Printf("Database %s does not exist. Trying to create it...", cfg.DBName)
			// Connects without DB to create the database
			dbConfigSinDB := cfg
			dbConfigSinDB.DBName = ""
			dbAux, err := sql.Open("mysql", dbConfigSinDB.FormatDSN())
			if err != nil {
				log.Fatalf("Could not connect without DB to create the database: %v", err)
			}
			defer dbAux.Close()
			_, err = dbAux.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.DBName))
			if err != nil {
				log.Fatalf("Error creating the database %s: %v", cfg.DBName, err)
			}
			log.Printf("Database %s created. Retrying connection...", cfg.DBName)
			continue // try again
		}
		log.Fatalf("Could not connect to the database: %v", pingErr)
	}
}
