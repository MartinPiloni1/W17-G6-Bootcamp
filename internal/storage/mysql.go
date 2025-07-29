package storage

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// formats the parameters of the connection with the standar structure
// mysql.Config struct, adding extra parameters
func NewMySQLConfig(host, port, user, pass, dbname string) mysql.Config {
	return mysql.Config{
		User:                 user,
		Passwd:               pass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		DBName:               dbname,
		ParseTime:            true,
		AllowNativePasswords: true,
		MultiStatements:      true,
		Params:               map[string]string{"charset": "utf8mb4"},
	}
}

// generates a connector to mysql, check if it can be used with Ping
// returns error if it cant ping it, or the connector pointer for future use in the app
func InitMySQLConnection(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("invalid args to mysql: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("could not connect to mysql: %w", err)
	}
	return db, nil
}
