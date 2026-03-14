package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	driver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouse struct {
	conn ch.Conn

	connAttempts int
	connTimeout  time.Duration
}

const (
	defaultAttempts = 10
	defaultTimeout  = time.Second
)

func NewMigrationDB(addr, db, user, pass string) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"clickhouse://%s:%s@%s/%s",
		user,
		pass,
		addr,
		db,
	)

	conn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

func New(addr, db, user, pass string) (*ClickHouse, error) {

	conn, err := ch.Open(&ch.Options{
		Addr: []string{addr},
		Auth: ch.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	for i := 0; i < defaultAttempts; i++ {

		if err := conn.Ping(ctx); err == nil {
			return &ClickHouse{conn: conn}, nil
		}

		time.Sleep(defaultTimeout)
	}

	return nil, fmt.Errorf("clickhouse connection failed")
}

func (c *ClickHouse) Conn() ch.Conn {
	return c.conn
}

func (c *ClickHouse) Close() {
	c.conn.Close()
}

func (c *ClickHouse) PrepareBatch(ctx context.Context, query string) (driver.Batch, error) {
	return c.conn.PrepareBatch(ctx, query)
}
