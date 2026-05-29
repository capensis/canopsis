package amqp

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestConnection_waitReconnect_GivenSetConn_ShouldReturnConnImmediately(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		amqpConn := NewMockamqp091Conn(ctrl)
		c := newTestConn(amqpConn)

		var got amqp091Conn
		go func() {
			got = c.waitReconnect(nil, nil)
		}()

		synctest.Wait()

		if got != amqpConn {
			t.Fatalf("expected %+v, got %+v", amqpConn, got)
		}
	})
}

func TestConnection_waitReconnect_GivenConnDoneClosed_ShouldReturnNil(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestConn(nil)
		var got amqp091Conn
		go func() {
			got = c.waitReconnect(nil, nil)
		}()

		synctest.Wait()
		close(c.done)
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})
}

func TestConnection_waitReconnect_GivenConnClosed_ShouldReturnNil(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestConn(nil)
		var got amqp091Conn
		go func() {
			got = c.waitReconnect(nil, nil)
		}()

		synctest.Wait()
		cerr := c.Close()
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}
	})
}

func TestConnection_waitReconnect_GivenCallerDoneClosed_ShouldReturnNil(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestConn(nil)
		callerDone := make(chan struct{})
		var got amqp091Conn
		go func() {
			got = c.waitReconnect(callerDone, nil)
		}()

		synctest.Wait()
		close(callerDone)
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})
}

func TestConnection_waitReconnect_GivenReconnection_ShouldReturnNewConn(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		var got amqp091Conn
		go func() {
			got = c.waitReconnect(nil, nil)
		}()

		synctest.Wait()
		expected := NewMockamqp091Conn(ctrl)
		simulateReconnect(c, expected)
		synctest.Wait()

		if got != expected {
			t.Fatalf("expected %+v, got %+v", expected, got)
		}
	})
}

func TestConnection_waitReconnect_GivenStaleConn_ShouldSkipIt(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stale := NewMockamqp091Conn(ctrl)
		c := newTestConn(stale)
		var got amqp091Conn
		go func() {
			got = c.waitReconnect(nil, stale)
		}()

		synctest.Wait()
		expected := NewMockamqp091Conn(ctrl)
		simulateReconnect(c, expected)
		synctest.Wait()

		if got != expected {
			t.Fatalf("expected %+v, got %+v", expected, got)
		}
	})
}

func TestConnection_waitReconnect_GivenMultipleWaiters_ShouldUnblockAllOfThem(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		size := 5
		got := make([]amqp091Conn, size)
		for i := 0; i < size; i++ {
			go func() {
				got[i] = c.waitReconnect(nil, nil)
			}()
		}

		synctest.Wait()
		expected := NewMockamqp091Conn(ctrl)
		simulateReconnect(c, expected)
		synctest.Wait()

		for i := 0; i < size; i++ {
			if got[i] != expected {
				t.Fatalf("[%d] expected %+v, got %+v", i, expected, got[i])
			}
		}
	})
}

func TestConnection_waitReconnect_GivenMultipleReconnections_ShouldReturnLastConn(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		var got1 amqp091Conn
		go func() {
			got1 = c.waitReconnect(nil, nil)
		}()

		synctest.Wait()
		conn1 := NewMockamqp091Conn(ctrl)
		simulateReconnect(c, conn1)
		synctest.Wait()

		if got1 != conn1 {
			t.Fatalf("expected %+v, got %+v", conn1, got1)
		}

		simulateDisconnect(c)

		var got2 amqp091Conn
		go func() {
			got2 = c.waitReconnect(nil, conn1)
		}()

		synctest.Wait()
		conn2 := NewMockamqp091Conn(ctrl)
		simulateReconnect(c, conn2)
		synctest.Wait()

		if got2 != conn2 {
			t.Fatalf("expected %+v, got %+v", conn2, got2)
		}
	})
}

func TestConnection_handleReconnect_GivenExhaustRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		c.reconnectCount = 12
		c.reconnectTimeout = 8 * time.Millisecond
		amqpDialer := NewMockamqpDialer(ctrl)
		amqpDialer.EXPECT().Dial().Times(c.reconnectCount).Return(nil, errors.New("Dial error"))
		c.dialer = amqpDialer

		var elapsed time.Duration
		stopped := false
		go func() {
			start := time.Now()
			c.handleReconnect()
			elapsed = time.Since(start)
			stopped = true
		}()

		reconnectTimeout := c.reconnectTimeout
		timeout := reconnectTimeout
		for i := 1; i < c.reconnectCount-1; i++ {
			reconnectTimeout = min(2*reconnectTimeout, maxReconnectTimeout)
			timeout += reconnectTimeout
		}

		time.Sleep(timeout)
		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if elapsed != timeout {
			t.Fatalf("expected %+v, got %+v", timeout, elapsed)
		}
	})
}

func TestConnection_handleReconnect_GivenReconnectAfterFailures_ShouldResetAttemptsAndTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		c.reconnectCount = 12
		c.reconnectTimeout = 8 * time.Millisecond
		amqpConn1 := NewMockamqp091Conn(ctrl)
		var notifClose1 chan *amqp.Error
		amqpConn1.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifClose1 = ch
		})
		amqpConn2 := NewMockamqp091Conn(ctrl)
		var notifClose2 chan *amqp.Error
		amqpConn2.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifClose2 = ch
		})
		amqpDialer := NewMockamqpDialer(ctrl)
		gomock.InOrder(
			amqpDialer.EXPECT().Dial().Return(amqpConn1, nil),
			amqpDialer.EXPECT().Dial().Times(c.reconnectCount-1).Return(nil, errors.New("Dial error")),
			amqpDialer.EXPECT().Dial().Return(amqpConn2, nil),
			amqpDialer.EXPECT().Dial().Times(c.reconnectCount).Return(nil, errors.New("Dial error")),
		)
		c.dialer = amqpDialer

		reconnectTimeout := c.reconnectTimeout
		timeout := reconnectTimeout
		for i := 1; i < c.reconnectCount-1; i++ {
			reconnectTimeout = min(2*reconnectTimeout, maxReconnectTimeout)
			timeout += reconnectTimeout
		}

		var elapsed time.Duration
		stopped := false
		go func() {
			start := time.Now()
			c.handleReconnect()
			elapsed = time.Since(start)
			stopped = true
		}()

		synctest.Wait()

		close(notifClose1)
		time.Sleep(timeout)

		synctest.Wait()

		close(notifClose2)
		time.Sleep(timeout)

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if elapsed != 2*timeout {
			t.Fatalf("expected %+v, got %+v", 2*timeout, elapsed)
		}
	})
}

func TestConnection_handleReconnect_GivenNoRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		c.reconnectCount = 0
		c.reconnectTimeout = 0
		amqpDialer := NewMockamqpDialer(ctrl)
		amqpDialer.EXPECT().Dial().Return(nil, errors.New("Dial error"))
		c.dialer = amqpDialer

		var elapsed time.Duration
		stopped := false
		go func() {
			start := time.Now()
			c.handleReconnect()
			elapsed = time.Since(start)
			stopped = true
		}()

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if elapsed > 0 {
			t.Fatalf("expected 0, got %+v", elapsed)
		}
	})
}

func TestConnection_handleReconnect_GivenClosedConn_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)

		amqpDialer := NewMockamqpDialer(ctrl)
		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().NotifyClose(gomock.Any())
		amqpConn.EXPECT().Close()
		amqpDialer.EXPECT().Dial().Return(amqpConn, nil)
		c.dialer = amqpDialer

		stopped := false
		go func() {
			c.handleReconnect()
			stopped = true
		}()

		synctest.Wait()

		cerr := c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}
	})
}

func TestConnection_handleReconnect_GivenClosedConnDuringRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		c.reconnectCount = 3
		c.reconnectTimeout = time.Second

		amqpDialer := NewMockamqpDialer(ctrl)
		amqpDialer.EXPECT().Dial().Times(2).Return(nil, errors.New("Dial error"))
		c.dialer = amqpDialer

		stopped := false
		var elapsed time.Duration
		go func() {
			start := time.Now()
			c.handleReconnect()
			elapsed = time.Since(start)
			stopped = true
		}()

		synctest.Wait()

		timeout := c.reconnectTimeout + c.reconnectTimeout/2
		time.Sleep(timeout)

		synctest.Wait()

		cerr := c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if elapsed != timeout {
			t.Fatalf("expected %v, got %v", timeout, elapsed)
		}
	})
}

func TestConnection_handleReconnect_GivenCloseNotif_ShouldReconnect(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		amqpDialer := NewMockamqpDialer(ctrl)

		amqpConn1 := NewMockamqp091Conn(ctrl)
		var notifClose1 chan *amqp.Error
		amqpConn1.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifClose1 = ch
		})
		amqpConn1.EXPECT().Close().Times(0)

		amqpConn2 := NewMockamqp091Conn(ctrl)
		amqpConn2.EXPECT().NotifyClose(gomock.Any())
		amqpConn2.EXPECT().Close()

		amqpDialer.EXPECT().Dial().Return(amqpConn1, nil)
		amqpDialer.EXPECT().Dial().Return(amqpConn2, nil)
		c.dialer = amqpDialer

		go func() {
			c.handleReconnect()
		}()

		synctest.Wait()

		close(notifClose1)

		synctest.Wait()

		if c.amqpConn != amqpConn2 {
			t.Fatalf("expected %+v, got %+v", amqpConn2, c.amqpConn)
		}

		cerr := c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()
	})
}

func TestConnection_handleReconnect_GivenCloseNotifAndNoRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		c.reconnectCount = 0
		c.reconnectTimeout = 0
		amqpDialer := NewMockamqpDialer(ctrl)

		amqpConn := NewMockamqp091Conn(ctrl)
		var notifClose chan *amqp.Error
		amqpConn.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifClose = ch
		})
		amqpConn.EXPECT().Close().Times(0)

		amqpDialer.EXPECT().Dial().Return(amqpConn, nil)
		c.dialer = amqpDialer

		stopped := false
		go func() {
			c.handleReconnect()
			stopped = true
		}()

		synctest.Wait()

		close(notifClose)

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}
	})
}

func TestConnection_handleReconnect_GivenCloseDuringSuccessfulReconnect_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)

		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().Close()

		amqpDialer := NewMockamqpDialer(ctrl)
		dialCh := make(chan struct{})
		amqpDialer.EXPECT().Dial().DoAndReturn(func() (amqp091Conn, error) {
			<-dialCh

			return amqpConn, nil
		})
		c.dialer = amqpDialer

		stopped := false
		go func() {
			c.handleReconnect()
			stopped = true
		}()

		synctest.Wait()

		cerr := c.Close()

		synctest.Wait()

		close(dialCh)

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}
	})
}

func TestConnection_Close_ShouldBeIdempotent(t *testing.T) {
	c := newTestConn(nil)

	if err := c.Close(); err != nil {
		t.Fatalf("expected nil, got %+v", err)
	}

	if err := c.Close(); !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("expected %+v, got %+v", ErrConnectionClosed, err)
	}
}

func TestConnection_Channel_GivenSetConn_ShouldReturnChannelImmediately(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().NotifyClose(gomock.Any())
		amqpConn.EXPECT().Close()

		amqpChannel := NewMockamqp091Channel(ctrl)
		amqpConn.EXPECT().Channel().Return(amqpChannel, nil)
		amqpChannel.EXPECT().NotifyClose(gomock.Any())
		amqpChannel.EXPECT().Close()

		c := newTestConn(amqpConn)

		var got Channel
		var gerr error
		go func() {
			got, gerr = c.Channel(t.Context())
		}()

		synctest.Wait()

		if got == nil {
			t.Fatal("expected Channel, got nil")
		}

		if gerr != nil {
			t.Fatalf("expected nil, got %+v", gerr)
		}

		cerr := got.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		cerr = c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()
	})
}

func TestConnection_Channel_GivenConnClosed_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestConn(nil)

		var got Channel
		var gerr error
		go func() {
			got, gerr = c.Channel(t.Context())
		}()

		synctest.Wait()
		cerr := c.Close()
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if !errors.Is(gerr, ErrConnectionClosed) {
			t.Fatalf("expected %+v, got %+v", ErrConnectionClosed, gerr)
		}

		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}
	})
}

func TestConnection_Channel_GivenCallerCtxCanceled_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestConn(nil)
		ctx, cancel := context.WithCancel(t.Context())

		var got Channel
		var gerr error
		go func() {
			got, gerr = c.Channel(ctx)
		}()

		synctest.Wait()
		cancel()
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if !errors.Is(gerr, context.Canceled) {
			t.Fatalf("expected %+v, got %+v", context.Canceled, gerr)
		}
	})
}

func TestConnection_Channel_GivenReconnection_ShouldReturnNewChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)

		var got Channel
		var gerr error
		go func() {
			got, gerr = c.Channel(t.Context())
		}()

		synctest.Wait()

		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().NotifyClose(gomock.Any())
		amqpConn.EXPECT().Close()

		amqpChannel := NewMockamqp091Channel(ctrl)
		amqpConn.EXPECT().Channel().Return(amqpChannel, nil)
		amqpChannel.EXPECT().NotifyClose(gomock.Any())
		amqpChannel.EXPECT().Close()

		simulateReconnect(c, amqpConn)

		synctest.Wait()

		if got == nil {
			t.Fatal("expected Channel, got nil")
		}

		if gerr != nil {
			t.Fatalf("expected nil, got %+v", gerr)
		}

		cerr := got.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		cerr = c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()
	})
}

func TestConnection_Channel_GivenMultipleCallers_ShouldUnblockAllOfThem(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestConn(nil)
		size := 5
		got := make([]Channel, size)
		gerrs := make([]error, size)
		for i := 0; i < size; i++ {
			go func() {
				got[i], gerrs[i] = c.Channel(t.Context())
			}()
		}

		synctest.Wait()

		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().NotifyClose(gomock.Any()).Times(size)
		amqpConn.EXPECT().Close()

		for i := 0; i < size; i++ {
			amqpChannel := NewMockamqp091Channel(ctrl)
			amqpConn.EXPECT().Channel().Return(amqpChannel, nil)
			amqpChannel.EXPECT().NotifyClose(gomock.Any())
			amqpChannel.EXPECT().Close()
		}

		simulateReconnect(c, amqpConn)

		synctest.Wait()

		for i := 0; i < size; i++ {
			if got[i] == nil {
				t.Fatalf("[%d] expected Channel, got nil", i)
			}

			if gerrs[i] != nil {
				t.Fatalf("[%d] expected nil, got %+v", i, gerrs[i])
			}
		}

		for i := range got {
			cerr := got[i].Close()
			if cerr != nil {
				t.Fatalf("[%d] expected nil, got %+v", i, cerr)
			}
		}

		cerr := c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()
	})
}

func TestConnection_Channel_GivenExhaustRetries_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reconnects := 3
		reconnTimeout := time.Second
		amqpDialer := NewMockamqpDialer(ctrl)
		amqpDialer.EXPECT().Dial().Times(reconnects).Return(nil, errors.New("Dial error"))
		c := newConn(amqpDialer, reconnects, reconnTimeout, zerolog.Nop())

		var got Channel
		var gerr error
		go func() {
			got, gerr = c.Channel(t.Context())
		}()

		time.Sleep(time.Duration(reconnects) * reconnTimeout)
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if !errors.Is(gerr, ErrConnectionClosed) {
			t.Fatalf("expected %+v, got %+v", ErrConnectionClosed, gerr)
		}
	})
}

// newTestConn builds a connection without starting handleReconnect so tests
// can drive amqpConn / reconnNotify / done directly and simulate dial events.
func newTestConn(amqpConn amqp091Conn) *connection {
	return &connection{
		amqpConn:         amqpConn,
		done:             make(chan struct{}),
		reconnectCount:   10,
		reconnectTimeout: time.Second,
		reconnNotify:     make(chan struct{}),
		logger:           zerolog.Nop(),
	}
}

// simulateReconnect mimics what handleReconnect does on a successful dial:
// store the new conn under the lock, swap reconnNotify, then close the old one.
func simulateReconnect(c *connection, newConn amqp091Conn) {
	c.mx.Lock()
	c.amqpConn = newConn
	prevReconnNotify := c.reconnNotify
	c.reconnNotify = make(chan struct{})
	c.mx.Unlock()
	close(prevReconnNotify)
}

// simulateDisconnect mimics what handleReconnect does after notifyClose fires.
func simulateDisconnect(c *connection) {
	c.mx.Lock()
	c.amqpConn = nil
	c.mx.Unlock()
}
