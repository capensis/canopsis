package postgres

//go:generate mockgen -destination=../../mocks/lib/postgres/postgres.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres BasePool,Pool
//go:generate mockgen -destination=../../mocks/github.com/jackc/pgx/pgx.go github.com/jackc/pgx/v4 Rows,Tx,BatchResults

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const EnvURL = "CPS_POSTGRES_URL"

const (
	MetricsCriteria = "metrics_criteria"
	Entities        = "entities"
)

// See error codes table for "Class 57 - Operator Intervention" https://www.postgresql.org/docs/9.6/errcodes-appendix.html#ERRCODES-TABLE
const pgErrOperatorInterventionPrefix = "57"

func IsConnectionError(err error) bool {
	netError := &net.OpError{}
	if errors.As(err, &netError) {
		return true
	}

	pgError := &pgconn.PgError{}
	if errors.As(err, &pgError) && strings.HasPrefix(pgError.Code, pgErrOperatorInterventionPrefix) {
		return true
	}

	if err.Error() == "conn closed" {
		return true
	}

	return false
}

func GetConnStr() (string, error) {
	connStr := os.Getenv(EnvURL)
	if connStr == "" {
		return "", fmt.Errorf("environment variable %s empty", EnvURL)
	}

	return connStr, nil
}

func NewPool(ctx context.Context, retryCount int, minRetryTimeout time.Duration) (Pool, error) {
	connStr, err := GetConnStr()
	if err != nil {
		return nil, err
	}

	pgxPool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &poolWithRetries{
		pgxPool: pgxPool,

		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}, nil
}

type BasePool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
	Stat() *pgxpool.Stat
}

type Pool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) error
	Close()
	WithTransaction(ctx context.Context, f func(context.Context, pgx.Tx) error) error
	SetRetry(count int, timeout time.Duration)
	Stat() *pgxpool.Stat
}

type poolWithRetries struct {
	pgxPool BasePool

	retryCount      int
	minRetryTimeout time.Duration
}

func (p *poolWithRetries) SetRetry(count int, timeout time.Duration) {
	p.retryCount = count
	p.minRetryTimeout = timeout
}

func (p *poolWithRetries) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	var commandTag pgconn.CommandTag
	var err error
	p.retry(ctx, func() error {
		commandTag, err = p.pgxPool.Exec(ctx, sql, args...)
		return err
	})

	return commandTag, err
}

func (p *poolWithRetries) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	var rows pgx.Rows
	var err error
	p.retry(ctx, func() error {
		rows, err = p.pgxPool.Query(ctx, sql, args...)
		return err
	})

	return rows, err
}

func (p *poolWithRetries) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	rows, err := p.Query(ctx, sql, args...)

	return &row{
		err:  err,
		rows: rows,
	}
}

func (p *poolWithRetries) SendBatch(ctx context.Context, b *pgx.Batch) error {
	return p.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		br := tx.SendBatch(ctx, b)

		for {
			_, err := br.Exec()
			if err != nil {
				if err.Error() == "no result" {
					break
				}

				return err
			}
		}

		return br.Close()
	})
}

func (p *poolWithRetries) Close() {
	p.pgxPool.Close()
}

func (p *poolWithRetries) Stat() *pgxpool.Stat {
	return p.pgxPool.Stat()
}

func (p *poolWithRetries) WithTransaction(ctx context.Context, f func(context.Context, pgx.Tx) error) error {
	var err error
	p.retry(ctx, func() error {
		var tx pgx.Tx
		tx, err = p.pgxPool.Begin(ctx)
		if err != nil {
			return err
		}

		defer func() {
			if rbErr := tx.Rollback(ctx); rbErr != nil && err == nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				err = rbErr
			}
		}()

		err = f(ctx, tx)
		if err != nil {
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (p *poolWithRetries) retry(ctx context.Context, f func() error) {
	timeout := p.minRetryTimeout

	for i := 0; i <= p.retryCount; i++ {
		err := f()
		if err == nil {
			return
		}

		if p.retryCount == i || timeout == 0 {
			return
		}

		if !IsConnectionError(err) {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(timeout):
			timeout *= 2
		}
	}
}

type row struct {
	err  error
	rows pgx.Rows
}

func (r *row) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}

	defer r.rows.Close()

	if r.rows.Err() != nil {
		return r.rows.Err()
	}

	if !r.rows.Next() {
		return pgx.ErrNoRows
	}

	return r.rows.Scan(dest...)
}
