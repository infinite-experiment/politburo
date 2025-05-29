package common

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheService struct {
	cache *cache.Cache
}

func NewCacheService(defaultExpirationSeconds, cleanUpIntervalSeconds int) *CacheService {

	defaultExpiration := time.Duration(defaultExpirationSeconds) * time.Second
	cleanUpInterval := time.Duration(cleanUpIntervalSeconds) * time.Second
	c := cache.New(defaultExpiration, cleanUpInterval)
	return &CacheService{cache: c}
}

func (cs *CacheService) Set(key string, value interface{}, duration time.Duration) {
	cs.cache.Set(key, value, duration*time.Minute)
}

func (cs *CacheService) Get(key string) (interface{}, bool) {
	return cs.cache.Get(key)
}

func (cs *CacheService) Delete(key string) {
	cs.cache.Delete(key)
}
