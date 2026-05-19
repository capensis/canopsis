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

func TestChannel_waitReconnect_GivenSetChannel_ShouldReturnChannelImmediately(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		amqpCh := NewMockamqp091Channel(ctrl)
		c := newTestChannel(nil, amqpCh)

		var got amqp091Channel
		go func() {
			got = c.waitReconnect(t.Context(), nil)
		}()

		synctest.Wait()

		if got != amqpCh {
			t.Fatalf("expected %+v, got %+v", amqpCh, got)
		}
	})
}

func TestChannel_waitReconnect_GivenChannelDoneClosed_ShouldReturnNil(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestChannel(nil, nil)
		var got amqp091Channel
		go func() {
			got = c.waitReconnect(t.Context(), nil)
		}()

		synctest.Wait()
		close(c.done)
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})
}

func TestChannel_waitReconnect_GivenChannelClosed_ShouldReturnNil(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestChannel(nil, nil)
		var got amqp091Channel
		go func() {
			got = c.waitReconnect(t.Context(), nil)
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

func TestChannel_waitReconnect_GivenCallerCtxCanceled_ShouldReturnNil(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestChannel(nil, nil)
		ctx, cancel := context.WithCancel(t.Context())
		var got amqp091Channel
		go func() {
			got = c.waitReconnect(ctx, nil)
		}()

		synctest.Wait()
		cancel()
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}
	})
}

func TestChannel_waitReconnect_GivenReconnection_ShouldReturnNewChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		var got amqp091Channel
		go func() {
			got = c.waitReconnect(t.Context(), nil)
		}()

		synctest.Wait()
		expected := NewMockamqp091Channel(ctrl)
		simulateChannelReconnect(c, expected)
		synctest.Wait()

		if got != expected {
			t.Fatalf("expected %+v, got %+v", expected, got)
		}
	})
}

func TestChannel_waitReconnect_GivenStaleChannel_ShouldSkipIt(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stale := NewMockamqp091Channel(ctrl)
		c := newTestChannel(nil, stale)
		var got amqp091Channel
		go func() {
			got = c.waitReconnect(t.Context(), stale)
		}()

		synctest.Wait()
		expected := NewMockamqp091Channel(ctrl)
		simulateChannelReconnect(c, expected)
		synctest.Wait()

		if got != expected {
			t.Fatalf("expected %+v, got %+v", expected, got)
		}
	})
}

func TestChannel_waitReconnect_GivenMultipleWaiters_ShouldUnblockAllOfThem(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		size := 5
		got := make([]amqp091Channel, size)
		for i := 0; i < size; i++ {
			go func() {
				got[i] = c.waitReconnect(t.Context(), nil)
			}()
		}

		synctest.Wait()
		expected := NewMockamqp091Channel(ctrl)
		simulateChannelReconnect(c, expected)
		synctest.Wait()

		for i := 0; i < size; i++ {
			if got[i] != expected {
				t.Fatalf("[%d] expected %+v, got %+v", i, expected, got[i])
			}
		}
	})
}

func TestChannel_waitReconnect_GivenMultipleReconnections_ShouldReturnLastChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		var got1 amqp091Channel
		go func() {
			got1 = c.waitReconnect(t.Context(), nil)
		}()

		synctest.Wait()
		ch1 := NewMockamqp091Channel(ctrl)
		simulateChannelReconnect(c, ch1)
		synctest.Wait()

		if got1 != ch1 {
			t.Fatalf("expected %+v, got %+v", ch1, got1)
		}

		simulateChannelDisconnect(c)

		var got2 amqp091Channel
		go func() {
			got2 = c.waitReconnect(t.Context(), ch1)
		}()

		synctest.Wait()
		ch2 := NewMockamqp091Channel(ctrl)
		simulateChannelReconnect(c, ch2)
		synctest.Wait()

		if got2 != ch2 {
			t.Fatalf("expected %+v, got %+v", ch2, got2)
		}
	})
}

func TestChannel_handleReconnect_GivenExhaustRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		c.conn.reconnectCount = 3
		c.conn.reconnectTimeout = time.Second
		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().Channel().Times(c.conn.reconnectCount).Return(nil, errors.New("Channel error"))

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
			stopped = true
		}()

		time.Sleep(time.Duration(c.conn.reconnectCount) * c.conn.reconnectTimeout)
		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}
	})
}

func TestChannel_handleReconnect_GivenNoRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		c.conn.reconnectCount = 0
		c.conn.reconnectTimeout = 0
		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().Channel().Return(nil, errors.New("Channel error"))

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
			stopped = true
		}()

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}
	})
}

func TestChannel_handleReconnect_GivenClosedChannel_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)

		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().NotifyClose(gomock.Any())

		amqpCh := NewMockamqp091Channel(ctrl)
		amqpCh.EXPECT().NotifyClose(gomock.Any())
		amqpCh.EXPECT().Close()
		amqpConn.EXPECT().Channel().Return(amqpCh, nil)

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
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

func TestChannel_handleReconnect_GivenClosedChannelDuringRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		c.conn.reconnectCount = 3
		c.conn.reconnectTimeout = time.Second

		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().Channel().Times(2).Return(nil, errors.New("Channel error"))

		stopped := false
		var elapsed time.Duration
		go func() {
			start := time.Now()
			c.handleReconnect(amqpConn)
			elapsed = time.Since(start)
			stopped = true
		}()

		synctest.Wait()

		timeout := c.conn.reconnectTimeout + c.conn.reconnectTimeout/2
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

func TestChannel_handleReconnect_GivenCloseChannelNotif_ShouldReconnect(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().NotifyClose(gomock.Any()).Times(2)

		amqpCh1 := NewMockamqp091Channel(ctrl)
		var notifChClose1 chan *amqp.Error
		amqpCh1.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifChClose1 = ch
		})
		amqpCh1.EXPECT().Close().Times(0)

		amqpCh2 := NewMockamqp091Channel(ctrl)
		amqpCh2.EXPECT().NotifyClose(gomock.Any())
		amqpCh2.EXPECT().Close()

		amqpConn.EXPECT().Channel().Return(amqpCh1, nil)
		amqpConn.EXPECT().Channel().Return(amqpCh2, nil)

		go func() {
			c.handleReconnect(amqpConn)
		}()

		synctest.Wait()

		close(notifChClose1)

		synctest.Wait()

		if c.amqpChannel != amqpCh2 {
			t.Fatalf("expected %+v, got %+v", amqpCh2, c.amqpChannel)
		}

		cerr := c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()
	})
}

func TestChannel_handleReconnect_GivenCloseConnNotif_ShouldReconnectConnAndChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		amqpConn1 := NewMockamqp091Conn(ctrl)
		var notifConnClose1 chan *amqp.Error
		amqpConn1.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifConnClose1 = ch
		})

		amqpConn2 := NewMockamqp091Conn(ctrl)
		amqpConn2.EXPECT().NotifyClose(gomock.Any())

		amqpCh1 := NewMockamqp091Channel(ctrl)
		amqpCh1.EXPECT().NotifyClose(gomock.Any())
		amqpCh1.EXPECT().Close().Times(0)

		amqpCh2 := NewMockamqp091Channel(ctrl)
		amqpCh2.EXPECT().NotifyClose(gomock.Any())
		amqpCh2.EXPECT().Close()

		amqpConn1.EXPECT().Channel().Return(amqpCh1, nil)
		amqpConn2.EXPECT().Channel().Return(amqpCh2, nil)

		go func() {
			c.handleReconnect(amqpConn1)
		}()

		synctest.Wait()

		close(notifConnClose1)

		synctest.Wait()

		simulateReconnect(c.conn, amqpConn2)

		synctest.Wait()

		if c.amqpChannel != amqpCh2 {
			t.Fatalf("expected %+v, got %+v", amqpCh2, c.amqpChannel)
		}

		cerr := c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()
	})
}

func TestChannel_handleReconnect_GivenClosedConn_ShouldReconnectConnAndChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		amqpConn1 := NewMockamqp091Conn(ctrl)
		amqpConn1.EXPECT().Channel().Return(nil, amqp.ErrClosed)

		amqpConn2 := NewMockamqp091Conn(ctrl)
		amqpConn2.EXPECT().NotifyClose(gomock.Any())

		amqpCh2 := NewMockamqp091Channel(ctrl)
		amqpCh2.EXPECT().NotifyClose(gomock.Any())
		amqpCh2.EXPECT().Close()
		amqpConn2.EXPECT().Channel().Return(amqpCh2, nil)

		go func() {
			c.handleReconnect(amqpConn1)
		}()

		synctest.Wait()

		simulateReconnect(c.conn, amqpConn2)

		synctest.Wait()

		if c.amqpChannel != amqpCh2 {
			t.Fatalf("expected %+v, got %+v", amqpCh2, c.amqpChannel)
		}

		cerr := c.Close()
		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}

		synctest.Wait()
	})
}

func TestChannel_handleReconnect_GivenCloseConnNotifAndFailedConnReconnect_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		amqpConn := NewMockamqp091Conn(ctrl)
		var notifConnClose chan *amqp.Error
		amqpConn.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifConnClose = ch
		})

		amqpCh := NewMockamqp091Channel(ctrl)
		amqpCh.EXPECT().NotifyClose(gomock.Any())
		amqpCh.EXPECT().Close().Times(0)

		amqpConn.EXPECT().Channel().Return(amqpCh, nil)

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
			stopped = true
		}()

		synctest.Wait()

		close(notifConnClose)

		synctest.Wait()

		close(c.conn.done)

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}
	})
}

func TestChannel_handleReconnect_GivenClosedConnAndFailedConnReconnect_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().Channel().Return(nil, amqp.ErrClosed)

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
			stopped = true
		}()

		synctest.Wait()

		close(c.conn.done)

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}
	})
}

func TestChannel_handleReconnect_GivenCloseConnNotifAndClosedChannelDuringReconnect_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		amqpConn := NewMockamqp091Conn(ctrl)
		var notifConnClose chan *amqp.Error
		amqpConn.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifConnClose = ch
		})

		amqpCh := NewMockamqp091Channel(ctrl)
		amqpCh.EXPECT().NotifyClose(gomock.Any())
		amqpCh.EXPECT().Close().Times(0)

		amqpConn.EXPECT().Channel().Return(amqpCh, nil)

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
			stopped = true
		}()

		synctest.Wait()

		close(notifConnClose)

		synctest.Wait()

		cerr := c.Close()

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}
	})
}

func TestChannel_handleReconnect_GivenClosedConnAndClosedChannelDuringReconnect_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().Channel().Return(nil, amqp.ErrClosed)

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
			stopped = true
		}()

		synctest.Wait()

		cerr := c.Close()

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}
	})
}

func TestChannel_handleReconnect_GivenCloseChannelNotifAndNoRetries_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		c.conn.reconnectCount = 0
		c.conn.reconnectTimeout = 0
		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().NotifyClose(gomock.Any())

		amqpCh := NewMockamqp091Channel(ctrl)
		var notifClose chan *amqp.Error
		amqpCh.EXPECT().NotifyClose(gomock.Any()).Do(func(ch chan *amqp.Error) {
			notifClose = ch
		})
		amqpCh.EXPECT().Close().Times(0)

		amqpConn.EXPECT().Channel().Return(amqpCh, nil)

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
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

func TestChannel_handleReconnect_GivenCloseChannelDuringSuccessfulReconnect_ShouldStop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)

		amqpCh := NewMockamqp091Channel(ctrl)
		amqpCh.EXPECT().Close()

		amqpConn := NewMockamqp091Conn(ctrl)
		channelCh := make(chan struct{})
		amqpConn.EXPECT().Channel().DoAndReturn(func() (amqp091Channel, error) {
			<-channelCh

			return amqpCh, nil
		})

		stopped := false
		go func() {
			c.handleReconnect(amqpConn)
			stopped = true
		}()

		synctest.Wait()

		cerr := c.Close()

		synctest.Wait()

		close(channelCh)

		synctest.Wait()

		if !stopped {
			t.Fatalf("expected true, got %+v", stopped)
		}

		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}
	})
}

func TestChannel_Close_ShouldBeIdempotent(t *testing.T) {
	c := newTestChannel(nil, nil)

	if err := c.Close(); err != nil {
		t.Fatalf("expected nil, got %+v", err)
	}

	if err := c.Close(); err != nil {
		t.Fatalf("expected nil, got %+v", err)
	}
}

func TestChannel_ConsumeWithContext_GivenSetConn_ShouldReturnChannelImmediately(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		amqpCh := NewMockamqp091Channel(ctrl)
		amqpCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		amqpCh.EXPECT().NotifyClose(gomock.Any())

		c := newTestChannel(nil, amqpCh)

		var got <-chan amqp.Delivery
		var gotNotifyClose <-chan *amqp.Error
		var gerr error
		go func() {
			got, gotNotifyClose, gerr = c.ConsumeWithContext(t.Context(), "test", "test", false, false, false, false, nil)
		}()

		synctest.Wait()

		if got == nil {
			t.Fatal("expected consume channel, got nil")
		}

		if gotNotifyClose == nil {
			t.Fatal("expected notify close channel, got nil")
		}

		if gerr != nil {
			t.Fatalf("expected nil, got %+v", gerr)
		}
	})
}

func TestChannel_ConsumeWithContext_GivenChannelClosed_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestChannel(nil, nil)

		var got <-chan amqp.Delivery
		var gotNotifyClose <-chan *amqp.Error
		var gerr error
		go func() {
			got, gotNotifyClose, gerr = c.ConsumeWithContext(t.Context(), "test", "test", false, false, false, false, nil)
		}()

		synctest.Wait()
		cerr := c.Close()
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if gotNotifyClose != nil {
			t.Fatalf("expected nil, got %+v", gotNotifyClose)
		}

		if !errors.Is(gerr, ErrChannelClosed) {
			t.Fatalf("expected %+v, got %+v", ErrChannelClosed, gerr)
		}

		if cerr != nil {
			t.Fatalf("expected nil, got %+v", cerr)
		}
	})
}

func TestChannel_ConsumeWithContext_GivenCallerCtxCanceled_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := newTestChannel(nil, nil)
		ctx, cancel := context.WithCancel(t.Context())

		var got <-chan amqp.Delivery
		var gotNotifyClose <-chan *amqp.Error
		var gerr error
		go func() {
			got, gotNotifyClose, gerr = c.ConsumeWithContext(ctx, "test", "test", false, false, false, false, nil)
		}()

		synctest.Wait()
		cancel()
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if gotNotifyClose != nil {
			t.Fatalf("expected nil, got %+v", gotNotifyClose)
		}

		if !errors.Is(gerr, context.Canceled) {
			t.Fatalf("expected %+v, got %+v", context.Canceled, gerr)
		}
	})
}

func TestChannel_ConsumeWithContext_GivenReconnection_ShouldReturnNewChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)

		var got <-chan amqp.Delivery
		var gotNotifyClose <-chan *amqp.Error
		var gerr error
		go func() {
			got, gotNotifyClose, gerr = c.ConsumeWithContext(t.Context(), "test", "test", false, false, false, false, nil)
		}()

		synctest.Wait()

		amqpCh := NewMockamqp091Channel(ctrl)
		amqpCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		amqpCh.EXPECT().NotifyClose(gomock.Any())

		simulateChannelReconnect(c, amqpCh)

		synctest.Wait()

		if got == nil {
			t.Fatal("expected consume channel, got nil")
		}

		if gotNotifyClose == nil {
			t.Fatal("expected notify close channel, got nil")
		}

		if gerr != nil {
			t.Fatalf("expected nil, got %+v", gerr)
		}
	})
}

func TestChannel_ConsumeWithContext_GivenMultipleCallers_ShouldUnblockAllOfThem(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := newTestChannel(nil, nil)
		size := 5
		got := make([]<-chan amqp.Delivery, size)
		gotNotifyClose := make([]<-chan *amqp.Error, size)
		gerrs := make([]error, size)
		for i := 0; i < size; i++ {
			go func() {
				got[i], gotNotifyClose[i], gerrs[i] = c.ConsumeWithContext(t.Context(), "test", "test", false, false, false, false, nil)
			}()
		}

		synctest.Wait()

		amqpCh := NewMockamqp091Channel(ctrl)

		for i := 0; i < size; i++ {
			amqpCh.EXPECT().
				ConsumeWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(make(chan amqp.Delivery), nil)
			amqpCh.EXPECT().NotifyClose(gomock.Any())
		}

		simulateChannelReconnect(c, amqpCh)

		synctest.Wait()

		for i := 0; i < size; i++ {
			if got[i] == nil {
				t.Fatalf("[%d] expected consume channel, got nil", i)
			}

			if gotNotifyClose[i] == nil {
				t.Fatalf("[%d] expected notify close channel, got nil", i)
			}

			if gerrs[i] != nil {
				t.Fatalf("[%d] expected nil, got %+v", i, gerrs[i])
			}
		}
	})
}

func TestChannel_ConsumeWithContext_GivenExhaustRetries_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reconnects := 3
		reconnTimeout := time.Second

		amqpConn := NewMockamqp091Conn(ctrl)
		amqpConn.EXPECT().Channel().Times(reconnects).Return(nil, errors.New("Channel error"))

		c := newChannel(&connection{
			reconnectCount:   reconnects,
			reconnectTimeout: reconnTimeout,
			logger:           zerolog.Nop(),
		}, amqpConn)

		var got <-chan amqp.Delivery
		var gotNotifyClose <-chan *amqp.Error
		var gerr error
		go func() {
			got, gotNotifyClose, gerr = c.ConsumeWithContext(t.Context(), "test", "test", false, false, false, false, nil)
		}()

		time.Sleep(time.Duration(reconnects) * reconnTimeout)
		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if gotNotifyClose != nil {
			t.Fatalf("expected nil, got %+v", gotNotifyClose)
		}

		if !errors.Is(gerr, ErrChannelClosed) {
			t.Fatalf("expected %+v, got %+v", ErrChannelClosed, gerr)
		}
	})
}

func TestChannel_ConsumeWithContext_GivenErrClosed_ShouldRetryWithNewChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		amqpCh1 := NewMockamqp091Channel(ctrl)
		amqpCh1.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, amqp.ErrClosed)

		c := newTestChannel(nil, amqpCh1)

		var got <-chan amqp.Delivery
		var gotNotifyClose <-chan *amqp.Error
		var gerr error
		go func() {
			got, gotNotifyClose, gerr = c.ConsumeWithContext(t.Context(), "test", "test", false, false, false, false, nil)
		}()

		synctest.Wait()

		amqpCh2 := NewMockamqp091Channel(ctrl)
		amqpCh2.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		amqpCh2.EXPECT().NotifyClose(gomock.Any())

		simulateChannelReconnect(c, amqpCh2)

		synctest.Wait()

		if got == nil {
			t.Fatal("expected consume channel, got nil")
		}

		if gotNotifyClose == nil {
			t.Fatal("expected notify close channel, got nil")
		}

		if gerr != nil {
			t.Fatalf("expected nil, got %+v", gerr)
		}
	})
}

func TestChannel_ConsumeWithContext_GivenUnknownErr_ShouldNotRetry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		amqpCh := NewMockamqp091Channel(ctrl)
		unknownErr := errors.New("unknown error")
		amqpCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, unknownErr)

		c := newTestChannel(nil, amqpCh)

		var got <-chan amqp.Delivery
		var gotNotifyClose <-chan *amqp.Error
		var gerr error
		go func() {
			got, gotNotifyClose, gerr = c.ConsumeWithContext(t.Context(), "test", "test", false, false, false, false, nil)
		}()

		synctest.Wait()

		if got != nil {
			t.Fatalf("expected nil, got %+v", got)
		}

		if gotNotifyClose != nil {
			t.Fatalf("expected nil, got %+v", gotNotifyClose)
		}

		if !errors.Is(gerr, unknownErr) {
			t.Fatalf("expected %+v, got %+v", unknownErr, gerr)
		}
	})
}

// newTestChannel builds a channel without starting handleReconnect so tests
// can drive amqpConn / reconnNotify / done directly and simulate dial events.
func newTestChannel(amqpConn amqp091Conn, amqpCh amqp091Channel) *channel {
	return &channel{
		conn:         newTestConn(amqpConn),
		amqpChannel:  amqpCh,
		done:         make(chan struct{}),
		reconnNotify: make(chan struct{}),
		logger:       zerolog.Nop(),
	}
}

// simulateChannelReconnect mimics what handleReconnect does on a successful dial:
// store the new channel under the lock, swap reconnNotify, then close the old one.
func simulateChannelReconnect(c *channel, newChannel amqp091Channel) {
	c.mx.Lock()
	c.amqpChannel = newChannel
	prevReconnNotify := c.reconnNotify
	c.reconnNotify = make(chan struct{})
	c.mx.Unlock()
	close(prevReconnNotify)
}

// simulateChannelDisconnect mimics what handleReconnect does after notifyClose fires.
func simulateChannelDisconnect(c *channel) {
	c.mx.Lock()
	c.amqpChannel = nil
	c.mx.Unlock()
}
