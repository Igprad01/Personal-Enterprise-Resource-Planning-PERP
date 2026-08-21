package database

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"personal-erp-backend/internal/config"
	"personal-erp-backend/internal/domain"
)

//go:embed schema.sql
var schemaSQL string

func Connect(cfg *config.Config) (*gorm.DB, error) {
	if err := ensureDatabase(cfg); err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

func ensureDatabase(cfg *config.Config) error {
	conn, err := gorm.Open(mysql.Open(cfg.DSNWithoutDB()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("open db (bootstrap): %w", err)
	}
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping mysql server: %w", err)
	}

	stmt := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		cfg.DBName,
	)
	if err := conn.Exec(stmt).Error; err != nil {
		return fmt.Errorf("create database: %w", err)
	}
	return nil
}

func Migrate(db *gorm.DB) error {
	for _, stmt := range splitStatements(schemaSQL) {
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("run migration: %w", err)
		}
	}
	return nil
}

func SeedCategories(db *gorm.DB) error {
	seeds := []domain.Category{
		{Name: "Work", Color: "#2563eb", Icon: ptr("briefcase")},
		{Name: "Personal", Color: "#16a34a", Icon: ptr("user")},
		{Name: "Health", Color: "#dc2626", Icon: ptr("heart")},
		{Name: "Learning", Color: "#9333ea", Icon: ptr("book")},
		{Name: "Errands", Color: "#d97706", Icon: ptr("shopping-cart")},
	}
	for _, s := range seeds {
		if err := db.Where("name = ?", s.Name).FirstOrCreate(&s).Error; err != nil {
			return fmt.Errorf("seed category %s: %w", s.Name, err)
		}
	}
	return nil
}

func ptr(s string) *string { return &s }

func splitStatements(script string) []string {
	var stmts []string
	var current strings.Builder
	inString := false

	for _, r := range script {
		switch {
		case r == '\'':
			inString = !inString
			current.WriteRune(r)
		case r == ';' && !inString:
			stmt := strings.TrimSpace(current.String())
			current.Reset()
			if stmt != "" {
				stmts = append(stmts, stmt)
			}
		default:
			current.WriteRune(r)
		}
	}
	if tail := strings.TrimSpace(current.String()); tail != "" {
		stmts = append(stmts, tail)
	}
	return stmts
}