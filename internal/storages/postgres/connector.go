package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PSQL struct {
	timeout time.Duration
	pool    *pgxpool.Pool
}

func New() *PSQL {
	return &PSQL{
		pool: nil,
	}
}

func (p *PSQL) Start(ctx context.Context, url string, timeout time.Duration, migrationsPath string) error {
	const op = "PSQL Start"

	p.timeout = timeout

	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	pool, err := pgxpool.Connect(ctxTimeout, url)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.pool = pool

	err = doMigrate(url, migrationsPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}

func doMigrate(dbURL, migrationsPath string) error {
	const op = "PSQL Migrate"
	if migrationsPath == "" {
		return fmt.Errorf("%s: %s", op, "migrations-path is required")
	}

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PSQL) Stop() {
	if p.pool != nil {
		p.pool.Close()
	}
}
