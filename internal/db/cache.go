package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jpillora/backoff"
	"github.com/luckyAkbar/atec/internal/common"
	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// NilValue is used to store and indicate a key's value is nil, thus can be
// handled in a special way by the caller / receiver
const NilValue = "NIL"

// all known errors that might be returned by this module.
// you should handle these errors in the caller / receiver
var (
	ErrCacheKeyNotFound   = errors.New("key not found")
	ErrCacheNil           = errors.New("nil value stored on cache value")
	ErrLockWaitingTooLong = errors.New("wait too long for acquiring lock or read the locked value")
)

const (
	defaultWaitTime        = 15 * time.Second
	defaultNilCacheTimeout = 1 * time.Minute
	defaultMinBackoff      = 20 * time.Millisecond
	defaultMaxBackoff      = 200 * time.Millisecond
)

// RedisConnOpts redis connection options
type RedisConnOpts struct {
	Addr               string
	Password           string
	DB                 int
	MinIdleConns       int
	ConnMaxLifetimeSec int
}

// NewRedisClient create a new redis client
func NewRedisClient(opts RedisConnOpts) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:            opts.Addr,
		DB:              opts.DB,
		Password:        opts.Password,
		MinIdleConns:    opts.MinIdleConns,
		ConnMaxLifetime: time.Second * time.Duration(opts.ConnMaxLifetimeSec),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	logrus.Info("connected to redis")

	return client
}

// CacheKeeperIface cache keeper interface
type CacheKeeperIface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	SetNil(ctx context.Context, key string, expiration ...time.Duration) error
	SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error
	GetOrLock(ctx context.Context, key string) (string, *common.RedsyncMutexWrapper, error)
	AcquireLock(key string) (*common.RedsyncMutexWrapper, error)
	Del(ctx context.Context, key string) error
}

// CacheKeeper cache keeper
type CacheKeeper struct {
	redisClient            *redis.Client
	distributedLock        common.DistributedLockerIface
	defaultNilCacheTimeout time.Duration
	defaultMinBackoff      time.Duration
	defaultMaxBackoff      time.Duration
}

// NewCacheKeeper create a new cache keeper
func NewCacheKeeper(redisClient *redis.Client, distributedLock *common.DistributedLocker) *CacheKeeper {
	return &CacheKeeper{
		redisClient: redisClient,

		// if see the need to make this configurable, please use Setters
		defaultNilCacheTimeout: defaultNilCacheTimeout,
		distributedLock:        distributedLock,
		defaultMinBackoff:      defaultMinBackoff,
		defaultMaxBackoff:      defaultMaxBackoff,
	}
}

// Get get a key from cache
func (ck *CacheKeeper) Get(ctx context.Context, key string) (string, error) {
	val, err := ck.redisClient.Get(ctx, key).Result()
	switch err {
	default:
		return "", err
	case redis.Nil:
		return "", ErrCacheKeyNotFound
	case nil:
		break
	}

	if val == NilValue {
		return "", ErrCacheNil
	}

	return val, nil
}

// Set set a key with value
func (ck *CacheKeeper) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return ck.redisClient.Set(ctx, key, value, expiration).Err()
}

// SetJSON set a key with value in JSON format
func (ck *CacheKeeper) SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ck.Set(ctx, key, string(val), expiration)
}

// SetNil set a key with nil value defined as NilValue = "NIL" constant
func (ck *CacheKeeper) SetNil(ctx context.Context, key string, expiration ...time.Duration) error {
	exp := ck.defaultNilCacheTimeout
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	return ck.redisClient.Set(ctx, key, NilValue, exp).Err()
}

// GetOrLock get a key from cache, if not found, try to lock it and return the mutex
func (ck *CacheKeeper) GetOrLock(ctx context.Context, key string) (string, *common.RedsyncMutexWrapper, error) {
	cacheValue, err := ck.Get(ctx, key)
	switch err {
	default:
		return "", nil, err
	case nil:
		return cacheValue, nil, nil
	case ErrCacheKeyNotFound:
		break
	}

	// if mutex successfully obtained, try to lock it and return the mutex
	mutex, err := ck.AcquireLock(key)
	if err == nil {
		return "", mutex, nil
	}

	// backoff mechanism start here
	backoffStartTime := time.Now()
	backoffRule := backoff.Backoff{
		Min:    ck.defaultMinBackoff,
		Max:    ck.defaultMaxBackoff,
		Jitter: true,
	}

	for {
		if time.Since(backoffStartTime) > defaultWaitTime {
			return "", nil, ErrLockWaitingTooLong
		}

		time.Sleep(backoffRule.Duration())

		isLocked, err := ck.distributedLock.IsLocked(ctx, key)
		if err != nil {
			logrus.WithError(err).Error("failed to check if key is locked. decide to retry")

			continue
		}

		if isLocked {
			continue
		}

		cacheValue, err = ck.Get(ctx, key)
		switch err {
		default:
			return "", nil, err
		case nil:
			return cacheValue, nil, nil
		case ErrCacheKeyNotFound:
		}

		mutex, err := ck.distributedLock.GetLock(key)
		if err == nil {
			return "", mutex, mutex.Lock()
		}

		logrus.WithError(err).Debug("still failed to acquire lock, decide to continue retry")
	}
}

// AcquireLock acquire a lock for a given key
func (ck *CacheKeeper) AcquireLock(key string) (*common.RedsyncMutexWrapper, error) {
	return ck.distributedLock.GetLock(key)
}

// Del delete a key from cache
func (ck *CacheKeeper) Del(ctx context.Context, key string) error {
	return ck.redisClient.Del(ctx, key).Err()
}
