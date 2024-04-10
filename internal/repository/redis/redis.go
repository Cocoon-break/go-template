package redis

import (
	"context"
	"time"

	"go-template/pkg/zlog"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var _ Repo = (*repo)(nil)

type Repo interface {
	i()
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	TTL(key string) (time.Duration, error)
	Expire(key string, ttl time.Duration) bool
	ExpireAt(key string, ttl time.Time) bool
	Del(key string) bool
	Exists(keys ...string) bool
	Incr(key string) int64
	Close() error
}

type repo struct {
	prefix string
	client *redis.Client
}

func New(addr, pwd, prefix string, db int) (Repo, error) {
	client, err := connect(addr, pwd, prefix, db)
	if err != nil {
		return nil, err
	}
	return &repo{
		prefix: prefix,
		client: client,
	}, nil
}

// Set set some <key,value> into redis
func (c *repo) Set(key, value string, ttl time.Duration) error {
	ts := time.Now()
	defer logWithKey("set", key, ts)
	_, err := c.client.Set(context.Background(), key, value, ttl).Result()
	return err
}

// Get get some key from redis
func (c *repo) Get(key string) (string, error) {
	ts := time.Now()
	defer logWithKey("get", key, ts)
	return c.client.Get(context.Background(), key).Result()
}

// TTL get some key from redis
func (c *repo) TTL(key string) (time.Duration, error) {
	ttl, err := c.client.TTL(context.Background(), key).Result()
	if err != nil {
		return -1, errors.Wrapf(err, "redis get key: %s err", key)
	}
	return ttl, nil
}

func (c *repo) Del(key string) bool {
	ts := time.Now()
	defer logWithKey("del", key, ts)
	if key == "" {
		return true
	}
	value, _ := c.client.Del(context.Background(), key).Result()
	return value > 0
}

func (c *repo) Incr(key string) int64 {
	ts := time.Now()
	defer logWithKey("incr", key, ts)
	value, _ := c.client.Incr(context.Background(), key).Result()
	return value
}

// Expire expire some key
func (c *repo) Expire(key string, ttl time.Duration) bool {
	ok, _ := c.client.Expire(context.Background(), key, ttl).Result()
	return ok
}

// ExpireAt expire some key at some time
func (c *repo) ExpireAt(key string, ttl time.Time) bool {
	ok, _ := c.client.ExpireAt(context.Background(), key, ttl).Result()
	return ok
}

func (c *repo) Exists(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	value, _ := c.client.Exists(context.Background(), keys...).Result()
	return value > 0
}

// Close close redis client
func (c *repo) Close() error {
	return c.client.Close()
}

func (c *repo) i() {}

func connect(addr, pwd, prefix string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	_, err := client.Do(context.Background(), "ping", prefix).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func logWithKey(action, key string, start time.Time) {
	used := time.Since(start)
	fields := []zlog.Field{
		zlog.String("key", key),
		zlog.Duration("timeUsed", used),
	}
	if used > time.Second {
		zlog.Warn(action, fields...)
	} else {
		zlog.Debug(action, fields...)
	}
}
