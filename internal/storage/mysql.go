package storage

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConfig(host, port, user, pass, dbname string) mysql.Config {
	return mysql.Config{
		User:      user,
		Passwd:    pass,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%s", host, port),
		DBName:    dbname,
		ParseTime: true,
		// Params:               map[string]string{"charset": "utf8mb4,utf8"},
		AllowNativePasswords: true,
		MultiStatements:      true,
	}
}

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
