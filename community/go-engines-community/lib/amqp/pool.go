package amqp

import (
	"context"
	"errors"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrPoolClosed = errors.New("amqp channel pool is closed")

// ChannelPool is a borrow/return pool of Channels backed by a Connection.
// Channels are opened lazily on Get and returned for reuse via Put.
// At most limit channels are ever open at once; when the limit is reached and
// no idle channel is available, Get blocks until a channel is returned via Put,
// the pool is closed, or ctx is cancelled.
type ChannelPool interface {
	Get(ctx context.Context) (Channel, error)
	Put(ch Channel)
	Close() error
}

func NewChannelPool(conn Connection, limit int) ChannelPool {
	if limit < 1 {
		return &channelPool{
			conn: conn,
		}
	}

	tokens := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		tokens <- struct{}{}
	}

	return &channelPoolWithLimit{
		conn:   conn,
		idle:   make(chan Channel, limit),
		tokens: tokens,
		done:   make(chan struct{}),
	}
}

type channelPoolWithLimit struct {
	conn Connection

	mx       sync.Mutex
	all      []Channel
	isClosed bool

	idle   chan Channel
	tokens chan struct{}
	done   chan struct{}
}

func (p *channelPoolWithLimit) Get(ctx context.Context) (Channel, error) {
	select {
	case <-p.done:
		return nil, ErrPoolClosed
	case ch := <-p.idle:
		return ch, nil
	default:
	}

	select {
	case <-p.done:
		return nil, ErrPoolClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	case ch := <-p.idle:
		return ch, nil
	case <-p.tokens:
		return p.openAndRegister(ctx)
	}
}

func (p *channelPoolWithLimit) Put(ch Channel) {
	if ch == nil {
		return
	}

	select {
	case <-p.done:
	case p.idle <- ch:
	}
}

func (p *channelPoolWithLimit) Close() error {
	p.mx.Lock()

	if p.isClosed {
		p.mx.Unlock()

		return nil
	}

	p.isClosed = true
	all := p.all
	p.all = nil
	close(p.done)

	p.mx.Unlock()

	errs := make([]error, len(all))
	for i, ch := range all {
		errs[i] = ch.Close()
	}

	return errors.Join(errs...)
}

func (p *channelPoolWithLimit) openAndRegister(ctx context.Context) (Channel, error) {
	ch, err := p.conn.Channel(ctx)
	if err != nil {
		select {
		case p.tokens <- struct{}{}:
		default:
		}

		return nil, err
	}

	p.mx.Lock()
	defer p.mx.Unlock()

	if p.isClosed {
		_ = ch.Close()

		return nil, ErrPoolClosed
	}

	p.all = append(p.all, ch)

	return ch, nil
}

type channelPool struct {
	conn Connection

	mx       sync.Mutex
	all      []Channel
	idle     []Channel
	isClosed bool
}

func (p *channelPool) Get(ctx context.Context) (Channel, error) {
	p.mx.Lock()

	if p.isClosed {
		p.mx.Unlock()

		return nil, ErrPoolClosed
	}

	if len(p.idle) > 0 {
		ch := p.idle[0]
		p.idle = p.idle[1:]

		p.mx.Unlock()

		return ch, nil
	}

	p.mx.Unlock()

	ch, err := p.conn.Channel(ctx)
	if err != nil {
		return nil, err
	}

	p.mx.Lock()
	defer p.mx.Unlock()

	if p.isClosed {
		_ = ch.Close()

		return nil, ErrPoolClosed
	}

	p.all = append(p.all, ch)

	return ch, nil
}

func (p *channelPool) Put(ch Channel) {
	if ch == nil {
		return
	}

	p.mx.Lock()
	defer p.mx.Unlock()

	p.idle = append(p.idle, ch)
}

func (p *channelPool) Close() error {
	p.mx.Lock()

	if p.isClosed {
		p.mx.Unlock()

		return nil
	}

	p.isClosed = true
	all := p.all
	p.all = nil

	p.mx.Unlock()

	errs := make([]error, len(all))
	for i, ch := range all {
		errs[i] = ch.Close()
	}

	return errors.Join(errs...)
}

// NewPooledPublisher returns a Publisher that borrows a channel from p for each publish and
// returns it when the call completes.
func NewPooledPublisher(p ChannelPool) Publisher {
	return &pooledPublisher{pool: p}
}

type pooledPublisher struct {
	pool ChannelPool
}

func (p *pooledPublisher) PublishWithContext(
	ctx context.Context,
	exchange, key string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) error {
	ch, err := p.pool.Get(ctx)
	if err != nil {
		return err
	}

	defer p.pool.Put(ch)

	return ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}
