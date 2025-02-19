package repository

// CacheKey is the standart cache key type
// and should be used to represent a cache key
// from a resource that is unique or have special use
type CacheKey string

// all the cache key defined here will represent a resource that is unique, or have special use thus the cache key is constant.
// for example, the AllActivePackageCacheKey is a cache key that is used to store all active packages in the cache is unique across the system.
// if updates happen to the related resource of cache value, you must obtain the lock to the resource.
const (
	AllActivePackageCacheKey CacheKey = "github.com/luckyAkbar/atec:cache-key:const:all_active_packages"
)
