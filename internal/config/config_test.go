package config_test

import (
	"book-store/internal/config"
	"os"
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/require"
)

func writeTempConfig(t *testing.T, content string) string {
	f, err := os.CreateTemp("", "cfg_*.json")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

func TestLoadConfig_Success(t *testing.T) {
	valid := `{
  "db": {
    "host": "DB_HOST",
    "port": "DB_PORT",
    "user": "DB_USER",
    "password": "DB_PASSWORD",
    "name": "DB_NAME"
    }
}`
	path := writeTempConfig(t, valid)
	cfg, err := config.LoadConfig(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, "DB_HOST", cfg.GetHost())
	require.Equal(t, "DB_PORT", cfg.GetPort())
	require.Equal(t, "DB_USER", cfg.GetUser())
	require.Equal(t, "DB_PASSWORD", cfg.GetPassword())
	require.Equal(t, "DB_NAME", cfg.GetName())

}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := config.LoadConfig(filepath.Join(os.TempDir(), "does_not_exist.json"))
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTempConfig(t, `{"FieldA": "ok",`)
	_, err := config.LoadConfig(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected EOF")
}


func TestLoadConfig_MissingField(t *testing.T) {
    missing := `{}`
    path := writeTempConfig(t, missing)
	_, err := config.LoadConfig(path)
	require.Contains(t,err.Error(),"Error:Field validation for 'User'")
}
