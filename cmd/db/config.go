package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Carga .env y retorna los datos de configuración
func LoadEnv() {
	_ = godotenv.Load()
}

// Obtiene el DSN y config para MySQL
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

// Intenta abrir la base de datos y si no existe, la crea y vuelve a intentar.
func MustOpenDB() *sql.DB {
	LoadEnv()
	cfg := GetDBConfigFromEnv()
	for {
		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Fatalf("No se pudo abrir DB: %v", err)
		}

		// Intenta el ping real, posible error de "no existe base"
		pingErr := db.Ping()
		if pingErr == nil {
			// ¡OK!
			return db
		}

		// Detecta error "unknown database"
		if merr, ok := pingErr.(*mysql.MySQLError); ok && merr.Number == 1049 {
			log.Printf("La base de datos %s no existe. Intentando crearla...", cfg.DBName)
			// Conecta sin DB para crear la database
			dbConfigSinDB := cfg
			dbConfigSinDB.DBName = ""
			dbAux, err := sql.Open("mysql", dbConfigSinDB.FormatDSN())
			if err != nil {
				log.Fatalf("No se pudo conectar sin DB para crear la database: %v", err)
			}
			defer dbAux.Close()
			_, err = dbAux.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", cfg.DBName))
			if err != nil {
				log.Fatalf("Error creando la base %s: %v", cfg.DBName, err)
			}
			log.Printf("Base de datos %s creada. Retentando conexión...", cfg.DBName)
			continue // vuelve a intentar de nuevo
		}
		log.Fatalf("No se pudo conectar a la base de datos: %v", pingErr)
	}
}
