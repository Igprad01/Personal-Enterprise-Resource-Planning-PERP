package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"personal_erp/backend/internal/config"
)

//go:embed schema.sql
var schemaSQL string

func Connect(cfg config.Config) (*sql.DB, error) {
	if err := ensureDatabase(cfg); err != nil {
		return nil, err
	}

	pool, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}

// ensureDatabase membuat database jika belum ada sehingga server bisa
// dijalankan pertama kali tanpa setup manual di MySQL.
func ensureDatabase(cfg config.Config) error {
	conn, err := sql.Open("mysql", cfg.DSNWithoutDB())
	if err != nil {
		return fmt.Errorf("open db (bootstrap): %w", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("ping mysql server: %w", err)
	}

	stmt := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		cfg.DBName,
	)
	if _, err := conn.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("create database: %w", err)
	}
	return nil
}

func Migrate(pool *sql.DB) error {
	for _, stmt := range splitStatements(schemaSQL) {
		if _, err := pool.Exec(stmt); err != nil {
			return fmt.Errorf("run migration: %w", err)
		}
	}
	return nil
}

// splitStatements splits a SQL script into individual statements, ignoring
// semicolons inside single-quoted string literals.
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
