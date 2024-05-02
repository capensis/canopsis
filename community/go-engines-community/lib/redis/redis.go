package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// Cache<Type> gives you constants to use for different caches.
const (
	CacheAlarm = iota
	CacheEntity
	CacheService

	LockStorage
	QueueStorage
	ApiCacheQueue
	AxePeriodicalLockStorage
	RuleTotalEntitiesStorage
	AlarmGroupStorage
	CorrelationLockStorage
	EngineRunInfo
	PBehaviorLockStorage
	ActionScenarioStorage
	EntityServiceStorage
	FIFOMessageStatisticsStorage
	// EngineLockStorage is used for all redis locks. It should be used by all engines.
	EngineLockStorage
)

// Env vars for redis session
const (
	EnvURL = "CPS_REDIS_URL"
)

// NewOptions handles redis.Options creation based on
// the surl, which must be on the following shape:
//
// redis://[nouser:password@]host:port/int
// int must be un number indicating the database
//
// If you have a password for the database, no user is required.
// But it is required to not leave the user empty to avoid url
// parsing error.
func NewOptions(surl string, db int, logger zerolog.Logger,
	reconnectCount int, minReconnectTimeout time.Duration) (*redis.Options, error) {
	redisURL, err := url.ParseRequestURI(surl)
	if err != nil {
		return nil, err
	}

	redisPassword := ""
	redisPasswordSet := false
	if redisURL.User != nil {
		redisPassword, redisPasswordSet = redisURL.User.Password()
		if !redisPasswordSet {
			redisPassword = ""
		}
	}

	redisDB := db
	if db < 0 {
		if len(redisURL.Path) < 2 {
			return nil, errors.New("no database specified in url")
		}
		rdb, err := strconv.ParseInt(redisURL.Path[1:], 10, 32)
		if err != nil {
			return nil, err
		}

		redisDB = int(rdb)
	}

	redisOptions := redis.Options{
		Addr:     redisURL.Host,
		Password: redisPassword,
		DB:       redisDB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Debug().Str("Addr", redisURL.Host).Int("DB", redisDB).Msg("New connection is established")
			return nil
		},
	}

	redisOptions.MaxRetries = reconnectCount
	redisOptions.MinRetryBackoff = minReconnectTimeout
	if redisOptions.MaxRetries > 0 && redisOptions.MinRetryBackoff > 0 {
		redisOptions.MaxRetryBackoff = redisOptions.MinRetryBackoff << redisOptions.MaxRetries
	}

	return &redisOptions, nil
}

// NewFailoverOptions handles redis.FailoverOptions creation based on the provided
// url, which must be on the following shape:
//
// redis-sentinel://[password@]host1[:port1][,host2[:port2]][,hostN[:portN]][/database][?
//
//	[timeout=timeout[d|h|m|s|ms|us|ns]][&sentinelMasterId=sentinelMasterId]]
//
// As well supported password parameter same as in NewOptions():
//
//	redis://[nouser:password@]host:port/int
//
// With this form "nouser" is ignored, and "password" extracted only.
func NewFailoverOptions(sURL string, db int, logger zerolog.Logger,
	reconnectCount int, minReconnectTimeout time.Duration) (*redis.FailoverOptions, error) {
	redisURL, err := url.ParseRequestURI(sURL)
	if err != nil {
		return nil, err
	}

	if redisURL.RawQuery == "" {
		return nil, errors.New("no master specified in the url")
	}

	failoverOptions := redis.FailoverOptions{}

	if redisURL.User != nil {
		if password, passwordFound := redisURL.User.Password(); passwordFound {
			if username := redisURL.User.Username(); username != "" {
				failoverOptions.SentinelUsername = username
				failoverOptions.SentinelPassword = password
				if redisURL.Query().Has("redisPassword") {
					failoverOptions.Password = redisURL.Query().Get("redisPassword")
				}
			} else {
				failoverOptions.Password = password
			}
		} else if password := redisURL.User.Username(); password != "" {
			failoverOptions.Password = password
		}
	} else if redisURL.Query().Has("redisPassword") {
		failoverOptions.Password = redisURL.Query().Get("redisPassword")
	}

	failoverOptions.SentinelAddrs = strings.Split(redisURL.Host, ",")

	failoverOptions.DB = db
	if db < 0 {
		if len(redisURL.Path) < 2 {
			return nil, errors.New("no database specified in url")
		}
		failoverOptions.DB, err = strconv.Atoi(redisURL.Path[1:])
		if err != nil {
			return nil, err
		}
	}

	failoverOptions.MasterName = redisURL.Query().Get("sentinelMasterId")

	if redisIdleTimeoutStr := redisURL.Query().Get("timeout"); redisIdleTimeoutStr != "" {
		if redisIdleTimeout, err := time.ParseDuration(redisIdleTimeoutStr); err == nil {
			failoverOptions.ConnMaxIdleTime = redisIdleTimeout
		} else {
			return nil, fmt.Errorf("redis-sentinel timeout parameter error %w", err)
		}
	}

	failoverOptions.OnConnect = func(ctx context.Context, cn *redis.Conn) error {
		redisURL.User = url.User("") // hide password
		logger.Debug().Str("Addr", redisURL.String()).Msg("New connection is established")
		return nil
	}

	failoverOptions.MaxRetries = reconnectCount
	failoverOptions.MinRetryBackoff = minReconnectTimeout
	if failoverOptions.MaxRetries > 0 && failoverOptions.MinRetryBackoff > 0 {
		failoverOptions.MaxRetryBackoff = failoverOptions.MinRetryBackoff << failoverOptions.MaxRetries
	}

	return &failoverOptions, nil
}

// NewSession creates a new connection to a Redis database.
// Configuration is base on EnvCpsRedisUrl.
func NewSession(ctx context.Context, db int, logger zerolog.Logger, reconnectCount int, minReconnectTimeout time.Duration) (*redis.Client, error) {
	connectUrl := os.Getenv(EnvURL)
	if connectUrl == "" {
		return nil, fmt.Errorf("environment variable %s empty", EnvURL)
	}

	var redisClient *redis.Client
	readTimeout := 3 * time.Second // redis.Options.ReadTimeout default value
	if strings.HasPrefix(connectUrl, "redis-sentinel://") {
		failoverOptions, err := NewFailoverOptions(connectUrl, db, logger, reconnectCount, minReconnectTimeout)
		if err != nil {
			return nil, err
		}
		redisClient = redis.NewFailoverClient(failoverOptions)
	} else {
		redisOptions, err := NewOptions(connectUrl, db, logger, reconnectCount, minReconnectTimeout)
		if err != nil {
			return nil, err
		}
		redisClient = redis.NewClient(redisOptions)
	}
	if minReconnectTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, readTimeout)
		defer cancel()
	}
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}

func IsConnectionError(err error) bool {
	netErr := &net.OpError{}
	if errors.As(err, &netErr) {
		return true
	}

	s := err.Error()
	if s == "ERR max number of clients reached" {
		return true
	}
	if strings.HasPrefix(s, "LOADING ") {
		return true
	}
	if strings.HasPrefix(s, "READONLY ") {
		return true
	}
	if strings.HasPrefix(s, "CLUSTERDOWN ") {
		return true
	}

	return false
}
