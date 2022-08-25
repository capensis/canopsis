package postgres

import (
	"context"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/rs/zerolog"
)

type PoolProvider interface {
	GetPool() Pool
}

type techMetricsPoolProvider struct {
	configProvider      config.MetricsConfigProvider
	pool                Pool
	pMx                 sync.RWMutex
	logger              zerolog.Logger
	configCheckDuration time.Duration
	retries             int
	timeout             time.Duration
}

func NewTechMetricsPoolProvider(ctx context.Context, configProvider config.MetricsConfigProvider, retries int, timeout time.Duration, logger zerolog.Logger) PoolProvider {
	var err error

	p := &techMetricsPoolProvider{
		configProvider:      configProvider,
		pMx:                 sync.RWMutex{},
		logger:              logger,
		configCheckDuration: time.Second * 10,
		retries:             retries,
		timeout:             timeout,
	}

	if configProvider.Get().EnableTechMetrics {
		p.pool, err = NewTechMetricsPool(ctx, retries, timeout)
		if err != nil {
			p.logger.Err(err).Msg("can't create postgres pool for tech metrics")
		}
	}

	go p.poolWatcher(ctx)

	return p
}

func (p *techMetricsPoolProvider) GetPool() Pool {
	p.pMx.RLock()
	defer p.pMx.RUnlock()

	return p.pool
}

func (p *techMetricsPoolProvider) poolWatcher(ctx context.Context) {
	ticker := time.NewTicker(p.configCheckDuration)

	defer func() {
		if p.pool != nil {
			p.pool.Close()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.resolvePool(ctx)
		}
	}
}

func (p *techMetricsPoolProvider) resolvePool(ctx context.Context) {
	p.pMx.Lock()
	defer p.pMx.Unlock()

	var err error

	if p.configProvider.Get().EnableTechMetrics {
		if p.pool != nil {
			return
		}

		p.pool, err = NewTechMetricsPool(ctx, p.retries, p.timeout)
		if err != nil {
			p.logger.Err(err).Msg("can't create postgres pool for tech metrics")
		}
	} else {
		if p.pool == nil {
			return
		}

		p.pool.Close()
		p.pool = nil
	}
}
