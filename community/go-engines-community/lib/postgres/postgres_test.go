package postgres

import (
	"context"
	"errors"
	mock_v4 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/github.com/jackc/pgx"
	mock_postgres "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/postgres"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"net"
	"testing"
	"time"
)

func TestPool_Exec_GivenContextDone_ShouldAbortRetries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 1
	minRetryTimeout := 2 * time.Second

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockPgxPool.EXPECT().Exec(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
		return nil, &net.OpError{Err: errors.New("test error")}
	}).AnyTimes()

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	start := time.Now()
	commandTag, err := pool.Exec(ctx, sql)

	if time.Since(start) > time.Second {
		t.Errorf("expected abort retry but method worked too long %s", time.Since(start))
	}

	if err == nil {
		t.Errorf("expected error but got nothing")
	}

	if commandTag != nil {
		t.Errorf("expected nil result but got %+v", commandTag)
	}
}

func TestPool_Exec_GivenConnectionError_ShouldRetryMaxTries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockPgxPool.EXPECT().Exec(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
		return nil, &net.OpError{Err: errors.New("test error")}
	}).Times(retryCount + 1)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	commandTag, err := pool.Exec(ctx, sql)
	if err == nil {
		t.Errorf("expected error but got nothing")
	}

	if commandTag != nil {
		t.Errorf("expected nil result but got %+v", commandTag)
	}
}

func TestPool_Exec_GivenNotConnectionError_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	execCount := -1
	mockPgxPool.EXPECT().Exec(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
		execCount++

		if execCount == 0 {
			return nil, &net.OpError{Err: errors.New("test error")}
		}

		return nil, &pgconn.PgError{Code: "42P09"}
	}).Times(2)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	commandTag, err := pool.Exec(ctx, sql)
	if err == nil {
		t.Errorf("expected error but got nothing")
	}

	if commandTag != nil {
		t.Errorf("expected nil result but got %+v", commandTag)
	}
}

func TestPool_Exec_GivenConnectionError_ShouldRetryUntilSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	execCount := -1
	mockPgxPool.EXPECT().Exec(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
		execCount++
		switch execCount {
		case 0:
			return nil, &net.OpError{Err: errors.New("test error")}
		case 1:
			return nil, &pgconn.PgError{Code: "57000"}
		case 2:
			return nil, &pgconn.PgError{Code: "57P01"}
		}

		return pgconn.CommandTag{}, nil
	}).Times(retryCount + 1)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	commandTag, err := pool.Exec(ctx, sql)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if commandTag == nil {
		t.Errorf("expected result but got nil")
	}
}

func TestPool_Query_GivenNotConnectionError_ShouldRetryMaxTries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockPgxPool.EXPECT().Query(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
		return nil, &net.OpError{Err: errors.New("test error")}
	}).Times(retryCount + 1)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	rows, err := pool.Query(ctx, sql)
	if err == nil {
		t.Errorf("expected error but got nothing")
	}

	if rows != nil {
		t.Errorf("expected nil result but got %+v", rows)
	}
}

func TestPool_Query_GivenConnectionError_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	execCount := -1
	mockPgxPool.EXPECT().Query(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
		execCount++

		if execCount == 0 {
			return nil, &net.OpError{Err: errors.New("test error")}
		}

		return nil, &pgconn.PgError{Code: "42P09"}
	}).Times(2)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	rows, err := pool.Query(ctx, sql)
	if err == nil {
		t.Errorf("expected error but got nothing")
	}

	if rows != nil {
		t.Errorf("expected nil result but got %+v", rows)
	}
}

func TestPool_Query_GivenConnectionError_ShouldRetryUntilSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	execCount := -1
	mockPgxPool.EXPECT().Query(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
		execCount++
		switch execCount {
		case 0:
			return nil, &net.OpError{Err: errors.New("test error")}
		case 1:
			return nil, &pgconn.PgError{Code: "57000"}
		case 2:
			return nil, &pgconn.PgError{Code: "57P01"}
		}

		return mock_v4.NewMockRows(ctrl), nil
	}).Times(retryCount + 1)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	rows, err := pool.Query(ctx, sql)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if rows == nil {
		t.Errorf("expected result but got nil")
	}
}

func TestPool_QueryRow_GivenNotConnectionError_ShouldRetryMaxTries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockPgxPool.EXPECT().Query(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
		return nil, &net.OpError{Err: errors.New("test error")}
	}).Times(retryCount + 1)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	row := pool.QueryRow(ctx, sql)
	err := row.Scan()
	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

func TestPool_QueryRow_GivenConnectionError_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	execCount := -1
	mockPgxPool.EXPECT().Query(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
		execCount++

		if execCount == 0 {
			return nil, &net.OpError{Err: errors.New("test error")}
		}

		return nil, &pgconn.PgError{Code: "42P09"}
	}).Times(2)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	row := pool.QueryRow(ctx, sql)
	err := row.Scan()
	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

func TestPool_QueryRow_GivenConnectionError_ShouldRetryUntilSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sql := "test sql"
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockRows := mock_v4.NewMockRows(ctrl)
	mockRows.EXPECT().Scan().Return(nil)
	mockRows.EXPECT().Err().Return(nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().Close()
	execCount := -1
	mockPgxPool.EXPECT().Query(gomock.Any(), gomock.Eq(sql)).DoAndReturn(func(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
		execCount++
		switch execCount {
		case 0:
			return nil, &net.OpError{Err: errors.New("test error")}
		case 1:
			return nil, &pgconn.PgError{Code: "57000"}
		case 2:
			return nil, &pgconn.PgError{Code: "57P01"}
		}

		return mockRows, nil
	}).Times(retryCount + 1)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	row := pool.QueryRow(ctx, sql)
	err := row.Scan()
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestPool_SendBatch_GivenNotConnectionError_ShouldRetryMaxTries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b := &pgx.Batch{}
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockTx := mock_v4.NewMockTx(ctrl)
	beginExecCount := -1
	mockPgxPool.EXPECT().Begin(gomock.Any()).DoAndReturn(func(_ context.Context) (pgx.Tx, error) {
		beginExecCount++

		if beginExecCount == 0 {
			return mockTx, nil
		}

		return nil, &net.OpError{Err: errors.New("test error")}
	}).Times(retryCount + 1)
	mockTx.EXPECT().SendBatch(gomock.Any(), gomock.Eq(b)).DoAndReturn(func(_ context.Context, _ *pgx.Batch) pgx.BatchResults {
		mockBr := mock_v4.NewMockBatchResults(ctrl)
		mockBr.EXPECT().Exec().Return(nil, &net.OpError{Err: errors.New("test error")})

		return mockBr
	})
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	err := pool.SendBatch(ctx, b)
	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

func TestPool_SendBatch_GivenConnectionError_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b := &pgx.Batch{}
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockTx := mock_v4.NewMockTx(ctrl)
	mockPgxPool.EXPECT().Begin(gomock.Any()).DoAndReturn(func(_ context.Context) (pgx.Tx, error) {
		return mockTx, nil
	})
	mockTx.EXPECT().SendBatch(gomock.Any(), gomock.Eq(b)).DoAndReturn(func(_ context.Context, _ *pgx.Batch) pgx.BatchResults {
		mockBr := mock_v4.NewMockBatchResults(ctrl)
		mockBr.EXPECT().Exec().Return(nil, errors.New("test error"))

		return mockBr
	})
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	err := pool.SendBatch(ctx, b)
	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

func TestPool_SendBatch_GivenConnectionError_ShouldRetryUntilSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b := &pgx.Batch{}
	retryCount := 3
	minRetryTimeout := time.Millisecond

	mockPgxPool := mock_postgres.NewMockBasePool(ctrl)
	mockTx := mock_v4.NewMockTx(ctrl)
	beginExecCount := -1
	mockPgxPool.EXPECT().Begin(gomock.Any()).DoAndReturn(func(_ context.Context) (pgx.Tx, error) {
		beginExecCount++

		switch beginExecCount {
		case 0:
			return nil, &net.OpError{Err: errors.New("test error")}
		case 1:
			return nil, &pgconn.PgError{Code: "57000"}
		case 2:
			return nil, &pgconn.PgError{Code: "57P01"}
		}

		return mockTx, nil
	}).Times(retryCount + 1)
	mockTx.EXPECT().SendBatch(gomock.Any(), gomock.Eq(b)).DoAndReturn(func(_ context.Context, _ *pgx.Batch) pgx.BatchResults {
		mockBr := mock_v4.NewMockBatchResults(ctrl)
		mockBr.EXPECT().Exec().Return(nil, nil)
		mockBr.EXPECT().Exec().Return(nil, errors.New("no result"))
		mockBr.EXPECT().Close()

		return mockBr
	})
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	pool := poolWithRetries{
		pgxPool:         mockPgxPool,
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}

	err := pool.SendBatch(ctx, b)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}
