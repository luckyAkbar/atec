package common

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sweet-go/stdlib/helper"
)

const lockKeyPrefix = "lock"

const (
	defaultLockExpiry = 10 * time.Second
	defaultLockTries  = 1
)

// DistributedLockerIface is the interface for distributed locker
type DistributedLockerIface interface {
	GetLock(key string) (*RedsyncMutexWrapper, error)
	IsLocked(ctx context.Context, key string) (bool, error)
}

// DistributedLocker is a distributed locker
type DistributedLocker struct {
	lockConnPool redsyncredis.Pool

	defaultLockExpiry time.Duration
	defaultLockTries  int
}

// NewDistributedLocker create a new distributed locker
func NewDistributedLocker(redisClient *redis.Client) *DistributedLocker {
	pool := goredis.NewPool(redisClient)

	return &DistributedLocker{
		lockConnPool: pool,

		// all these value is hardcoded.
		defaultLockExpiry: defaultLockExpiry,
		defaultLockTries:  defaultLockTries,
	}
}

// GetLock will get the lock for the given key
func (dl *DistributedLocker) GetLock(key string) (*RedsyncMutexWrapper, error) {
	rs := redsync.New(dl.lockConnPool)
	mutex := rs.NewMutex(lockKeyPrefix+key,
		redsync.WithExpiry(dl.defaultLockExpiry),
		redsync.WithTries(dl.defaultLockTries),
	)

	if err := mutex.Lock(); err != nil {
		return nil, err
	}

	return NewRedsyncMutexWrapper(mutex), nil
}

// IsLocked will check if the key is locked or not
func (dl *DistributedLocker) IsLocked(ctx context.Context, key string) (bool, error) {
	client, err := dl.lockConnPool.Get(ctx)
	if err != nil {
		return false, err
	}

	defer helper.WrapCloser(client.Close)

	_, err = client.Get(lockKeyPrefix + key)
	switch err {
	default:
		return false, err
	case redis.Nil:
		return false, nil
	case nil:
		return true, nil
	}
}

// RedsyncMutex is the interface for wrapping redsync mutex
// used functionality. Thus, if you need to use other redsync
// functionality, you can add it here.
type RedsyncMutex interface {
	Lock() error
	Unlock() (bool, error)
}

// RedsyncMutexWrapper is the wrapper for redsync mutex
type RedsyncMutexWrapper struct {
	mutex RedsyncMutex
}

// NewRedsyncMutexWrapper creates a new RedsyncMutexWrapper
func NewRedsyncMutexWrapper(mutex RedsyncMutex) *RedsyncMutexWrapper {
	return &RedsyncMutexWrapper{
		mutex: mutex,
	}
}

// Lock call the Lock method of the mutex
func (r *RedsyncMutexWrapper) Lock() error {
	return r.mutex.Lock()
}

// Unlock call the Unlock method of the mutex
func (r *RedsyncMutexWrapper) Unlock() (bool, error) {
	return r.mutex.Unlock()
}
