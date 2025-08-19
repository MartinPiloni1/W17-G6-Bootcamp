package server_test

import (
	"os"
	"testing"

	"github.com/aaguero_meli/W17-G6-Bootcamp/cmd/server"
	"github.com/stretchr/testify/require"
)

func TestLoadServerConf_OK(t *testing.T) {
	os.Setenv("ADDRESS", "9001")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3333")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "1234")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
	}()

	conf, err := server.LoadServerConf()
	require.NoError(t, err)
	require.NotNil(t, conf)
	require.Equal(t, ":9001", conf.ServerAddr)
	require.Equal(t, "testdb", conf.DatabaseConfig.DBName)
	require.Equal(t, "localhost:3333", conf.DatabaseConfig.Addr)
	require.Equal(t, "root", conf.DatabaseConfig.User)
	require.Equal(t, "1234", conf.DatabaseConfig.Passwd)
}

func TestLoadServerConf_MissingDBHost(t *testing.T) {
	// Quita todas las env
	os.Unsetenv("ADDRESS")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")

	os.Setenv("DB_PORT", "3333")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "1234")

	conf, err := server.LoadServerConf()
	require.Error(t, err)
	require.Nil(t, conf)
	require.Contains(t, err.Error(), "DB conn settings not established")
}

func TestLoadServerConf_DefaultAddress(t *testing.T) {
	os.Setenv("DB_HOST", "H")
	os.Setenv("DB_PORT", "P")
	os.Setenv("DB_NAME", "N")
	os.Setenv("DB_USER", "U")
	os.Setenv("DB_PASS", "PW")
	os.Unsetenv("ADDRESS")
	defer func() {
		os.Unsetenv("ADDRESS")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
	}()

	conf, err := server.LoadServerConf()
	require.NoError(t, err)
	require.NotNil(t, conf)
	require.Equal(t, ":8080", conf.ServerAddr)
}
