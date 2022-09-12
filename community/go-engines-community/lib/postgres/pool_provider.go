package postgres

import (
	"context"
	"sync"
	"time"
)

type PoolProvider interface {
	Get(ctx context.Context) (Pool, error)
	Close()
}

func NewPoolProvider(retryCount int, minRetryTimeout time.Duration) PoolProvider {
	return &poolProvider{
		retryCount:      retryCount,
		minRetryTimeout: minRetryTimeout,
	}
}

type poolProvider struct {
	retryCount      int
	minRetryTimeout time.Duration

	poolMx sync.Mutex
	pool   Pool
}

func (p *poolProvider) Get(ctx context.Context) (Pool, error) {
	p.poolMx.Lock()
	defer p.poolMx.Unlock()
	if p.pool == nil {
		var err error
		p.pool, err = NewPool(ctx, p.retryCount, p.minRetryTimeout)
		if err != nil {
			return nil, err
		}
	}

	return p.pool, nil
}

func (p *poolProvider) Close() {
	p.poolMx.Lock()
	defer p.poolMx.Unlock()

	if p.pool != nil {
		p.pool.Close()
		p.pool = nil
	}
}
