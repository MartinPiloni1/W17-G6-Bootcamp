package config_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/cmd/db"
	"github.com/stretchr/testify/require"
)

func TestLoadEnv(t *testing.T) {
	require.NotPanics(t, func() {
		config.LoadEnv()
	})
}

func TestGetDBConfigFromEnv(t *testing.T) {
	os.Setenv("DB_USER", "myuser")
	os.Setenv("DB_PASS", "mypassword")
	os.Setenv("DB_HOST", "localhost:9999")
	os.Setenv("DB_NAME", "mydatabase")
	defer func() {
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_NAME")
	}()

	cfg := config.GetDBConfigFromEnv()
	require.Equal(t, "myuser", cfg.User)
	require.Equal(t, "mypassword", cfg.Passwd)
	require.Equal(t, "localhost:9999", cfg.Addr)
	require.Equal(t, "mydatabase", cfg.DBName)
	require.True(t, cfg.AllowNativePasswords)
	require.True(t, cfg.ParseTime)
	require.True(t, cfg.MultiStatements)
}

func TestMustOpenDB_ErrorEnPing(t *testing.T) {
	os.Setenv("DB_USER", "u")
	os.Setenv("DB_PASS", "q")
	os.Setenv("DB_HOST", "localhost:1") // puerto inválido
	os.Setenv("DB_NAME", "dbtest")

	if os.Getenv("BE_CRASHER") == "1" {
		config.MustOpenDB()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMustOpenDB_ErrorEnPing")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	exitError, ok := err.(*exec.ExitError)
	require.True(t, ok, "MustOpenDB no llamó a os.Exit() (log.Fatalf no usado?)")
	require.Equal(t, 1, exitError.ExitCode())
}
