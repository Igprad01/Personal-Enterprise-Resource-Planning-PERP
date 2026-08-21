package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	ServerPort string
	DBUser     string
	DBPass     string
	DBHost     string
	DBPort     string
	DBName     string
}

func Load() *Config {
	loadDotEnv()
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPass:     getEnv("DB_PASS", ""),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "personal_erp"),
	}
}

func (c *Config) DSN() string {
	return c.DBUser + ":" + c.DBPass + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" +
		c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func (c *Config) DSNWithoutDB() string {
	return c.DBUser + ":" + c.DBPass + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" +
		"?charset=utf8mb4&parseTime=True&loc=Local"
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func loadDotEnv() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		p := dir + string(os.PathSeparator) + ".env"
		if f, err := os.Open(p); err == nil {
			parseEnvFile(f)
			f.Close()
			return
		}
		parent := dir[:strings.LastIndex(dir, string(os.PathSeparator))]
		if parent == dir {
			return
		}
		dir = parent
	}
}

func parseEnvFile(f *os.File) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}
}