package activecheck

import (
	"fmt"
	"strings"
	"time"
)

// ActiveChecker trying to lock active mode on Start(), then Keepalive() lock while active
type ActiveChecker interface {
	// Start access to slice of Redis hosts, trying to set lock during provided period using Redlock algorithm
	Start(RedLocker, time.Duration) (bool, error)
	IsActive() bool
	SetActive()
	SetPassive()
	// Keepalive for active instance updates lock's expiration
	Keepalive() error

	GetValue() string
}

type AddressListFlags []string

func (alf *AddressListFlags) String() string {
	return fmt.Sprintf("%s", []string(*alf))
}

func (alf *AddressListFlags) Set(s string) error {
	const (
		connPrefix  = "tcp://"
		redisPrefix = "redis://"
	)
	s = strings.TrimPrefix(s, redisPrefix)
	if !strings.Contains(s, "://") {
		s = connPrefix + s
	}
	*alf = append(*alf, s)
	return nil
}

type activeCheck struct {
	lock         RedLocker
	active       bool
	lockDuration int64
	lockKey      string
	value        string
}

// NewActiveCheck instantiates ActiveChecker
func NewActiveCheck(lockKey string) ActiveChecker {
	return &activeCheck{
		lockKey: lockKey,
	}
}

func (ac *activeCheck) Start(redLock RedLocker, periodicalWaitTime time.Duration) (bool, error) {
	var err error

	if redLock == nil {
		return false, fmt.Errorf("nil Redlocker")
	}

	ac.lock = redLock

	ac.lockDuration = (periodicalWaitTime + 2*time.Second).Nanoseconds() / 1e6
	ac.value, err = ac.lock.Lock(ac.lockKey, ac.lockDuration)

	if err != nil {
		ac.SetPassive()
	} else {
		ac.SetActive()
	}

	return ac.active, nil
}

func (ac *activeCheck) IsActive() bool {
	return ac.active
}

func (ac *activeCheck) SetPassive() {
	ac.active = false
}

func (ac *activeCheck) SetActive() {
	ac.active = true
}

func (ac *activeCheck) Keepalive() error {
	_, err := ac.lock.ExpireLock(ac.lockKey, ac.value, ac.lockDuration)
	return err
}

func (ac *activeCheck) GetValue() string {
	return ac.value
}
