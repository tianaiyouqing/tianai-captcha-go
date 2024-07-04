package application

import (
	"github.com/patrickmn/go-cache"
	"github.com/tianaiyouqing/tianai-captcha-go/common/model"
	"time"
)

type CacheStore interface {
	getCache(key string) (value map[string]any, ok bool)
	getAndRemoveCache(key string) (value map[string]any, ok bool)
	SetCache(key string, data map[string]any, captchaInfo *model.ImageCaptchaInfo) error
}

// MemoryCacheStore 内存本地实现
type MemoryCacheStore struct {
	cache *cache.Cache
}

func NewMemoryCacheStore(expired time.Duration, cleanUpInterval time.Duration) *MemoryCacheStore {
	return &MemoryCacheStore{
		cache: cache.New(expired, cleanUpInterval),
	}
}

func (self *MemoryCacheStore) getCache(key string) (value map[string]any, ok bool) {
	data, ok := self.cache.Get(key)
	if ok {
		return data.(map[string]any), ok
	}
	return nil, ok
}

func (self *MemoryCacheStore) getAndRemoveCache(key string) (value map[string]any, ok bool) {
	data, ok := self.cache.Get(key)
	if ok {
		self.cache.Delete(key)
		return data.(map[string]any), ok
	}
	return nil, ok
}

func (self *MemoryCacheStore) SetCache(key string, data map[string]any, _ *model.ImageCaptchaInfo) error {
	self.cache.Set(key, data, cache.DefaultExpiration)
	return nil
}
