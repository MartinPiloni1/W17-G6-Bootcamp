package storage_test

import (
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/internal/storage"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestNewMySQLConfig(t *testing.T) {
	cfg := storage.NewMySQLConfig("host", "3306", "user", "pass", "mydb")
	require.Equal(t, "user", cfg.User)
	require.Equal(t, "pass", cfg.Passwd)
	require.Equal(t, "tcp", cfg.Net)
	require.Equal(t, "host:3306", cfg.Addr)
	require.Equal(t, "mydb", cfg.DBName)
	require.True(t, cfg.ParseTime)
	require.True(t, cfg.AllowNativePasswords)
	require.True(t, cfg.MultiStatements)
	require.Equal(t, "utf8mb4", cfg.Params["charset"])
}

func TestInitMySQLConnection_InvalidDSN(t *testing.T) {
	cfg := mysql.Config{
		User:   "",
		Passwd: "",
		Net:    "invalidnet",
		Addr:   "",
		DBName: "",
	}
	db, err := storage.InitMySQLConnection(cfg)
	require.Error(t, err)
	require.Nil(t, db)
}
