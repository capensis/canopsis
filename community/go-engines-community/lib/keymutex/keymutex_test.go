package keymutex_test

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/keymutex"
	"sync"
	"testing"
	"time"
)

const waitTimeout = time.Millisecond * 10

func TestKeyMutex_Lock_GivenKey_ShouldLockKey(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	key := "test"
	mx := keymutex.New()
	go func() {
		mx.Lock(key)
		done <- true
	}()

	waitDone(t, done)
}

func TestKeyMutex_Lock_GiveMultipleKeys_ShouldLockKeys(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	key1 := "test1"
	key2 := "test2"
	mx := keymutex.New()
	go func() {
		mx.Lock(key1)
		mx.Lock(key2)
		done <- true
	}()

	waitDone(t, done)
}

func TestKeyMutex_Lock_GivenMultipleLocks_ShouldWaitUnlockBeforeNextLock(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	key := "test"
	var unlockTime, lockTime time.Time
	var err error
	mx := keymutex.New()
	go func() {
		mx.Lock(key)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			mx.Lock(key)
			lockTime = time.Now()
		}()

		time.Sleep(waitTimeout / 2)
		err = mx.Unlock(key)
		unlockTime = time.Now()

		wg.Wait()
		done <- true
	}()

	waitDone(t, done)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if unlockTime.After(lockTime) {
		t.Errorf("expected lock after unlock")
	}
}

func TestKeyMutex_LockMultiple_GivenMultipleLocks_ShouldWaitUnlockBeforeNextLock(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	keys := []string{"test1", "test2"}
	mx := keymutex.New()
	var err error
	var unlockTime, lockTime1, lockTime2 time.Time

	go func() {
		mx.LockMultiple(keys...)
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			mx.Lock(keys[0])
			lockTime1 = time.Now()
		}()

		go func() {
			defer wg.Done()
			mx.Lock(keys[1])
			lockTime2 = time.Now()
		}()

		time.Sleep(waitTimeout / 2)
		err = mx.UnlockMultiple(keys...)
		unlockTime = time.Now()

		wg.Wait()

		done <- true
	}()

	waitDone(t, done)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if unlockTime.After(lockTime1) {
		t.Errorf("expected lock after unlock")
	}

	if unlockTime.After(lockTime2) {
		t.Errorf("expected lock after unlock")
	}
}

func TestKeyMutex_Unlock_GivenLock_ShouldUnlock(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	key := "test"
	mx := keymutex.New()
	var err error
	go func() {
		mx.Lock(key)
		err = mx.Unlock(key)
		done <- true
	}()

	waitDone(t, done)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestKeyMutex_Unlock_GivenMultipleUnlock_ShouldReturnError(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	key := "test"
	mx := keymutex.New()
	var err error
	go func() {
		mx.Lock(key)
		firstUnlockErr := mx.Unlock(key)
		if firstUnlockErr == nil {
			err = mx.Unlock(key)
		} else {
			t.Errorf("expected no error but got %v", firstUnlockErr)
		}

		done <- true
	}()

	waitDone(t, done)

	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

func TestKeyMutex_Unlock_GivenMoreLocksThenUnlocks_ShouldReturnError(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	key := "test"
	mx := keymutex.New()
	var err error
	go func() {
		mx.Lock(key)

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			mx.Lock(key)
			firstUnlockErr := mx.Unlock(key)
			if firstUnlockErr == nil {
				err = mx.Unlock(key)
			} else {
				t.Errorf("expected no error but got %v", firstUnlockErr)
			}
		}()

		time.Sleep(waitTimeout / 2)
		unlockErr := mx.Unlock(key)
		if unlockErr != nil {
			t.Errorf("expected no error but got %v", unlockErr)
		}

		wg.Wait()
		done <- true
	}()

	waitDone(t, done)

	if err == nil {
		t.Errorf("expected error but got nothing")
	}
}

func TestKeyMutex_Unlock_GivenNoLock_ShouldReturnError(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	key := "test"
	var err error
	mx := keymutex.New()

	go func() {
		err = mx.Unlock(key)
		done <- true
	}()

	waitDone(t, done)

	if err == nil {
		t.Error("expected error but got nothing")
	}
}

func TestKeyMutex_UnlockMultiple_GivenNoLock_ShouldReturnError(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	keys := []string{"test1", "test2", "test3", "test4"}
	var err error
	mx := keymutex.New()

	go func() {
		mx.LockMultiple(keys[0], keys[2])
		err = mx.UnlockMultiple(keys...)
		done <- true
	}()

	waitDone(t, done)

	if err == nil {
		t.Error("expected error but got nothing")
	}
}

func TestKeyMutex_UnlockMultiple_GivenMultipleUnlock_ShouldReturnError(t *testing.T) {
	done := make(chan bool)
	defer close(done)

	keys := []string{"test1", "test2"}
	var err error
	mx := keymutex.New()

	go func() {
		mx.LockMultiple(keys...)
		firstUnlockErr := mx.UnlockMultiple(keys...)
		if firstUnlockErr == nil {
			err = mx.UnlockMultiple(keys...)
		} else {
			t.Errorf("expected no error but got %v", firstUnlockErr)
		}
		done <- true
	}()

	waitDone(t, done)

	if err == nil {
		t.Error("expected error but got nothing")
	}
}

func BenchmarkKeyMutex_Lock(b *testing.B) {
	for name, keys := range genKeys(10, 10000, 10) {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				mx := keymutex.New()

				for _, key := range keys {
					mx.Lock(key)
				}
			}
		})
	}
}

func BenchmarkKeyMutex_Lock_Unlock(b *testing.B) {
	for name, keys := range genKeys(10, 10000, 10) {
		b.Run(name, func(b *testing.B) {
			mx := keymutex.New()

			for i := 0; i < b.N; i++ {
				for _, key := range keys {
					mx.Lock(key)
				}
				for _, key := range keys {
					_ = mx.Unlock(key)
				}
			}
		})
	}
}

func BenchmarkKeyMutex_LockMultiple(b *testing.B) {
	for name, keys := range genKeys(10, 10000, 10) {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				mx := keymutex.New()
				mx.LockMultiple(keys...)
			}
		})
	}
}

func BenchmarkKeyMutex_LockMultiple_UnlockMultiple(b *testing.B) {
	for name, keys := range genKeys(10, 10000, 10) {
		b.Run(name, func(b *testing.B) {
			mx := keymutex.New()

			for i := 0; i < b.N; i++ {
				mx.LockMultiple(keys...)
				_ = mx.UnlockMultiple(keys...)
			}
		})
	}
}

func waitDone(t *testing.T, done <-chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Error("timeout expired")
	case _, ok := <-done:
		if !ok {
			t.Error("channel closed")
		}
	}
}

func genKeys(minLen, maxLen, step int) map[string][]string {
	keys := make([]string, maxLen)
	for i := 0; i < maxLen; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}

	keysByLen := make(map[string][]string)
	for i := minLen; i <= maxLen; i *= step {
		keysByLen[fmt.Sprintf("len %d", i)] = keys[:i]
	}

	return keysByLen
}
