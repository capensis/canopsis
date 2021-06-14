package activecheck

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
)

const (
	// DefaultRetryCount is the max retry times for lock acquire
	DefaultRetryCount = 10

	// DefaultRetryDelay is upper wait time in millisecond for lock acquire retry
	DefaultRetryDelay = 200

	// ClockDriftFactor is clock drift factor, more information refers to doc
	ClockDriftFactor = 0.01

	// UnlockScript is redis lua script to release a lock
	UnlockScript = `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
        `
)

// RedLocker interface to implement Reddlock
type RedLocker interface {
	// Lock acquires lock for specified resource and period in ms, returns assigned value
	Lock(resource string, period int64) (string, error)
	// ExpireLock changes expiration for resource with value to new period in ms, returns remained period
	ExpireLock(resource, value string, period int64) (int64, error)
	// SetRetryCount updates acquire lock attempts count
	SetRetryCount(int)
	// SetRetryDelay updates retry delays in milliseconds between acquire lock attempts
	SetRetryDelay(int)
}

// RedLock holds the redis lock
type RedLock struct {
	retryCount  int
	retryDelay  int
	driftFactor float64

	clients []*RedClient
	quorum  int
}

// RedClient holds client to redis
type RedClient struct {
	addr string
	cli  *redis.Client
}

func parseConnString(addr string) (*redis.Options, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	opts := &redis.Options{
		Network: u.Scheme,
		Addr:    u.Host,
	}

	dbStr := strings.Trim(u.Path, "/")
	if dbStr == "" {
		dbStr = "0"
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, err
	}
	opts.DB = db

	password, ok := u.User.Password()
	if ok {
		opts.Password = password
	}

	for k, v := range u.Query() {
		timeout := 0
		switch k {
		case "DialTimeout", "ReadTimeout", "WriteTimeout":
			timeout, err = strconv.Atoi(v[0])
			if err != nil {
				return nil, fmt.Errorf("%s %v conversion error %s", k, v[0], err)
			}
		}
		switch k {
		case "DialTimeout":
			opts.DialTimeout = time.Duration(timeout)
		case "ReadTimeout":
			opts.ReadTimeout = time.Duration(timeout)
		case "WriteTimeout":
			opts.WriteTimeout = time.Duration(timeout)
		}
	}

	return opts, nil
}

// NewRedLock creates a RedLock
func NewRedLock(addrs []string) (RedLocker, error) {
	if (len(addrs) == 1 && addrs[0] == "") || len(addrs)%2 == 0 {
		return nil, fmt.Errorf("error redis server list length: %d", len(addrs))
	}

	clients := []*RedClient{}
	for _, addr := range addrs {
		opts, err := parseConnString(addr)
		if err != nil {
			return nil, err
		}
		cli := redis.NewClient(opts)
		clients = append(clients, &RedClient{addr, cli})
	}

	return &RedLock{
		retryCount:  DefaultRetryCount,
		retryDelay:  DefaultRetryDelay,
		driftFactor: ClockDriftFactor,
		quorum:      len(addrs)/2 + 1,
		clients:     clients,
	}, nil
}

// SetRetryCount sets acquire lock retry count
func (r *RedLock) SetRetryCount(count int) {
	if count <= 0 {
		return
	}
	r.retryCount = count
}

// SetRetryDelay sets acquire lock retry max internal in millisecond
func (r *RedLock) SetRetryDelay(delay int) {
	if delay <= 0 {
		return
	}
	r.retryDelay = delay
}

func getRandStr() string {
	b := make([]byte, 16)
	crand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func lockInstance(client *RedClient, resource string, val string, ttl int, c chan bool) {
	if client.cli == nil {
		c <- false
		return
	}
	reply := client.cli.SetNX(resource, val, time.Duration(ttl)*time.Millisecond)
	if reply.Err() != nil || !reply.Val() {
		c <- false
		return
	}
	c <- true
}

func setExpireInstance(client *RedClient, resource string, val string, ttl int, c chan bool) {
	if client.cli == nil {
		c <- false
		return
	}
	if resourceVal := client.cli.Get(resource); resourceVal.Err() != nil || resourceVal.Val() != val {
		c <- false
		return
	}
	reply := client.cli.Expire(resource, time.Duration(ttl)*time.Millisecond)
	if reply.Err() != nil || !reply.Val() {
		c <- false
		return
	}
	c <- true
}

func unlockInstance(client *RedClient, resource, val string, wg *sync.WaitGroup) {
	defer wg.Done()

	if client.cli != nil {
		client.cli.Eval(UnlockScript, []string{resource}, val)
	}
}

// Lock acquires a distribute lock
func (r *RedLock) Lock(resource string, ttl int64) (string, error) {
	val := getRandStr()
	for i := 0; i < r.retryCount; i++ {
		c := make(chan bool, len(r.clients))
		success := 0
		start := time.Now()

		for _, cli := range r.clients {
			go lockInstance(cli, resource, val, int(ttl), c)
		}
		for j := 0; j < len(r.clients); j++ {
			if <-c {
				success++
			}
		}
		close(c)

		drift := int64(float64(ttl)*r.driftFactor) + 2
		costTime := time.Since(start).Nanoseconds() / 1e6
		validityTime := ttl - costTime - drift

		if success >= r.quorum && validityTime > 0 {
			return val, nil
		}

		r.UnLock(resource, val)
		// Wait a random delay before to retry
		time.Sleep(time.Duration(rand.Intn(r.retryDelay)) * time.Millisecond)
	}

	return val, fmt.Errorf("failed to require lock %s", resource)
}

// UnLock releases an acquired lock
func (r *RedLock) UnLock(resource, val string) {
	var wg sync.WaitGroup
	wg.Add(len(r.clients))
	for _, cli := range r.clients {
		go unlockInstance(cli, resource, val, &wg)
	}
	wg.Wait()
}

// ExpireLock changes distributed lock expiration
func (r *RedLock) ExpireLock(resource, val string, ttl int64) (int64, error) {
	c := make(chan bool, len(r.clients))
	success := 0
	start := time.Now()

	for _, cli := range r.clients {
		go setExpireInstance(cli, resource, val, int(ttl), c)
	}
	for j := 0; j < len(r.clients); j++ {
		if <-c {
			success++
		}
	}
	close(c)

	drift := int64(float64(ttl)*r.driftFactor) + 2
	costTime := time.Since(start).Nanoseconds() / 1e6
	validityTime := ttl - costTime - drift
	if success >= r.quorum && validityTime > 0 {
		return validityTime, nil
	}

	r.UnLock(resource, val)

	return 0, fmt.Errorf("failed to require lock %s", resource)
}
