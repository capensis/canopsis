package redis_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRedisOptions(t *testing.T) {
	redisOptions, err := redis.NewOptions("redis://user:password@host:7777", 0, log.NewTestLogger(), 0, 0)
	if err != nil {
		t.Fatalf("redis options error: %v", err)
	}

	if redisOptions.Password != "password" {
		t.Fatalf("redis bad password: %s", redisOptions.Password)
	}

	if redisOptions.Addr != "host:7777" {
		t.Fatalf("redis bad addr: %s", redisOptions.Addr)
	}

	if redisOptions.DB != 0 {
		t.Fatalf("redis bad database: %d", redisOptions.DB)
	}
}

func TestBadRedisOptions(t *testing.T) {
	ctx := context.Background()

	Convey("Testing bad redis urls", t, func() {
		Convey("Bad url", func() {
			_, err := redis.NewOptions("bla://nrausitenrste,anursiet", -1, log.NewTestLogger(), 0, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Bad url - judgement day", func() {
			_, err := redis.NewOptions("", -1, log.NewTestLogger(), 0, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Bad url - access denied", func() {
			_, err := redis.NewOptions("redis://user@localhost/0", -1, log.NewTestLogger(), 0, 0)
			So(err, ShouldBeNil)
		})

		Convey("Bad db", func() {
			_, err := redis.NewOptions("redis://localhost/bleurk", -1, log.NewTestLogger(), 0, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Bad url - i'm calling you", func() {
			oldredisurl := os.Getenv(redis.EnvURL)
			os.Setenv(redis.EnvURL, "redis://anrsuitenrstau,nrasutie;;;;...")
			_, err := redis.NewSession(ctx, -1, log.NewTestLogger(), 0, 0)
			os.Setenv(redis.EnvURL, oldredisurl)
			So(err, ShouldNotBeNil)
		})

		Convey("Good url, wrong database - Bagdad cafe", func() {
			oldredisurl := os.Getenv(redis.EnvURL)
			os.Setenv(redis.EnvURL, "redis://Canyouexplainwhatyouwannado/0")
			_, err := redis.NewSession(ctx, -1, log.NewTestLogger(), 0, 0)
			os.Setenv(redis.EnvURL, oldredisurl)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestNewRedisSession(t *testing.T) {
	_, err := redis.NewSession(context.Background(), redis.CacheAlarm, log.NewTestLogger(), 0, 0)
	if err != nil {
		t.Fatalf("redis client error: %v", err)
	}
}

func BenchmarkRedisSpeed(b *testing.B) {
	ctx := context.Background()
	client, _ := redis.NewSession(ctx, redis.CacheAlarm, log.NewTestLogger(), 0, 0)

	for i := 1; i < b.N; i++ {
		rid := "titi_" + strconv.Itoa(i)
		if i%4 == 0 {
			err := client.Set(ctx, rid, "toto", time.Hour*4)
			if err.Err() != nil {
				fmt.Printf("error setting key %s: %v", rid, err)
			}
		}
		_ = client.Exists(ctx, rid)
	}
}

func TestNewFailoverOptions(t *testing.T) {
	redisOptions, err := redis.NewFailoverOptions("redis-sentinel://user:password@host:7777?sentinelMasterId=prime", 0, log.NewTestLogger(), 0, 0)
	if err != nil {
		t.Fatalf("redis options error: %v", err)
	}

	if redisOptions.Password != "" {
		t.Fatalf("redis bad password: %s", redisOptions.Password)
	}

	if redisOptions.SentinelUsername != "user" {
		t.Fatalf("redis bad sentinel username: %s", redisOptions.SentinelUsername)
	}

	if redisOptions.SentinelPassword != "password" {
		t.Fatalf("redis bad sentinel password: %s", redisOptions.SentinelPassword)
	}

	if len(redisOptions.SentinelAddrs) != 1 {
		t.Fatalf("redis bad SentinelAddrs: %s", redisOptions.SentinelAddrs)
	}

	if redisOptions.SentinelAddrs[0] != "host:7777" {
		t.Fatalf("redis bad addr: %s", redisOptions.SentinelAddrs)
	}

	if redisOptions.DB != 0 {
		t.Fatalf("redis bad database: %d", redisOptions.DB)
	}

	redisOptions, err = redis.NewFailoverOptions("redis-sentinel://password@host1:7777,host2:7778/3?timeout=1s&sentinelMasterId=prime", -1, log.NewTestLogger(), 0, 0)
	if err != nil {
		t.Fatalf("redis options error: %v", err)
	}

	if redisOptions.Password != "password" {
		t.Fatalf("redis bad password: %s", redisOptions.Password)
	}

	if len(redisOptions.SentinelAddrs) != 2 {
		t.Fatalf("redis bad SentinelAddrs: %s", redisOptions.SentinelAddrs)
	}

	if redisOptions.SentinelAddrs[0] != "host1:7777" && redisOptions.SentinelAddrs[1] != "host2:7777" {
		t.Fatalf("redis bad SentinelAddrs: %s", redisOptions.SentinelAddrs)
	}

	if redisOptions.DB != 3 {
		t.Fatalf("redis bad database: %d", redisOptions.DB)
	}

	if redisOptions.MasterName != "prime" {
		t.Fatalf("redis bad master: %s", redisOptions.MasterName)
	}

	redisOptions, err = redis.NewFailoverOptions("redis-sentinel://password@host1:7777,host2:7778/6?timeout=1s&sentinelMasterId=prime", 3, log.NewTestLogger(), 0, 0)
	if err != nil {
		t.Fatalf("redis options error: %v", err)
	}
	if redisOptions.Password != "password" {
		t.Fatalf("redis bad password: %s", redisOptions.Password)
	}
	if redisOptions.DB != 3 {
		t.Fatalf("redis bad database: %d", redisOptions.DB)
	}
	redisOptions, err = redis.NewFailoverOptions("redis-sentinel://:password@host1:7777,host2:7778/?timeout=1s&sentinelMasterId=prime", 3, log.NewTestLogger(), 0, 0)
	if err != nil {
		t.Fatalf("redis options error: %v", err)
	}
	if redisOptions.Password != "password" {
		t.Fatalf("redis bad password: %s", redisOptions.Password)
	}
}
