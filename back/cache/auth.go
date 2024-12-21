package cache

import (
	"app/model"
	"time"

	"github.com/patrickmn/go-cache"
)

type AuthCache struct {
	cache *cache.Cache
}

// NewAuthCache создает новый экземпляр AuthCache с инициализированным кэшем
func NewAuthCache(defaultExpiration, cleanupInterval time.Duration) *AuthCache {
	return &AuthCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Set добавляет значение в кэш с указанным ключом, используя время жизни по умолчанию
func (a *AuthCache) Set(key string, value model.CachedUser) {
	a.cache.SetDefault(key, value)
}

// Get извлекает значение из кэша по ключу
func (a *AuthCache) Get(key string) (*model.CachedUser, bool) {
	value, found := a.cache.Get(key)
	if found {
		user, ok := value.(model.CachedUser)
		if ok {
			return &user, true
		}
		return nil, false
	}
	return nil, false
}

// Delete удаляет значение из кэша по ключу
func (a *AuthCache) Delete(key string) {
	a.cache.Delete(key)
}
