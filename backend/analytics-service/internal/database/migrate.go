package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
)

func RunMigrations(ctx context.Context, db *sql.DB) error {

	if err := goose.SetDialect("clickhouse"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(db, "/app/database/migrations"); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("get version: %w", err)
	}

	log.Printf("clickhouse migrations applied, version=%d", version)

	return nil
}
