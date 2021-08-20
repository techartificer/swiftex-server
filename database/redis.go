package database

import (
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/techartificer/swiftex/config"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

var client *goredis.Client
var limmiterStore limiter.Store

func ConnectRedis() error {
	redisCfg := config.GetRedis()

	client = goredis.NewClient(&goredis.Options{
		Addr:     redisCfg.Host,
		Password: redisCfg.Password, // no password set
		DB:       0,                 // use default DB
	})
	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{Prefix: "rate_limit", CleanUpInterval: 1 * time.Minute})
	if err != nil {
		return err
	}
	limmiterStore = store
	return nil
}

func GetLimmiterStore() *limiter.Store {
	return &limmiterStore
}
