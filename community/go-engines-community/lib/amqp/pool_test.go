package amqp_test

import (
	"testing"
	"testing/synctest"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	"go.uber.org/mock/gomock"
)

func TestChannelPool_Get_GivenLimitN_ShouldCreateNChannels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	limit := 10
	conn := mock_amqp.NewMockConnection(ctrl)
	conn.EXPECT().
		Channel(gomock.Any()).
		Return(mock_amqp.NewMockChannel(ctrl), nil).
		Times(limit)
	pool := amqp.NewChannelPool(conn, limit, nil)
	for i := 0; i < limit; i++ {
		ch, err := pool.Get(t.Context())
		if err != nil {
			t.Fatalf("expected nil, got %+v", err)
		}

		if ch == nil {
			t.Fatalf("expected channel, got nil")
		}
	}
}

func TestChannelPool_Get_GivenLimit_ShouldWaitForIdleChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		limit := 2
		conn := mock_amqp.NewMockConnection(ctrl)
		conn.EXPECT().
			Channel(gomock.Any()).
			Return(mock_amqp.NewMockChannel(ctrl), nil).
			Times(limit)
		pool := amqp.NewChannelPool(conn, limit, nil)
		for i := 0; i < limit; i++ {
			ch, err := pool.Get(t.Context())
			if err != nil {
				t.Fatalf("expected nil, got %+v", err)
			}

			if ch == nil {
				t.Fatalf("expected channel, got nil")
			}
		}

		var got amqp.Channel
		var gerr error
		go func() {
			got, gerr = pool.Get(t.Context())
		}()

		synctest.Wait()
		pool.Put(mock_amqp.NewMockChannel(ctrl))
		synctest.Wait()

		if gerr != nil {
			t.Fatalf("expected nil, got %+v", gerr)
		}

		if got == nil {
			t.Fatalf("expected channel, got nil")
		}
	})
}
