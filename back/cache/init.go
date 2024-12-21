package cache

import (
	"app/config"
	"time"
)

var Auth *AuthCache

func Init() error {
	conf := config.File.AuthConfig

	TTL := time.Duration(conf.TimeToLive) * time.Minute
	CI := time.Duration(conf.CleanupInterval) * time.Minute

	Auth = NewAuthCache(TTL, CI)
	return nil
}
