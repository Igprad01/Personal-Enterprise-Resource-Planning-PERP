package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBUser     string
	DBPass     string
	DBHost     string
	DBPort     string
	DBName     string
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// loadDotEnv memuat .env dari direktori kerja atau dari direktori
// induk terdekat (mis. .env global di root repo). Env dari OS tetap
// menang karena godotenv tidak menimpa variabel yang sudah terisi.
func loadDotEnv() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

func Load() (Config, error) {
	loadDotEnv()

	var missing []string
	for _, key := range []string{"DB_HOST", "DB_NAME"} {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return Config{
		ServerPort: getenv("SERVER_PORT", "8080"),
		DBUser:     getenv("DB_USER", "root"),
		DBPass:     os.Getenv("DB_PASS"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getenv("DB_PORT", "3306"),
		DBName:     os.Getenv("DB_NAME"),
	}, nil
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

// DSNWithoutDB menghubungkan ke server MySQL tanpa memilih database,
// dipakai untuk auto-create database saat server pertama kali dijalankan.
func (c Config) DSNWithoutDB() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&charset=utf8mb4&loc=Local",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort)
}
