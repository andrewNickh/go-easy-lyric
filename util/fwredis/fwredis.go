package fwredis

import (
	"easy-lyric/config"
	"easy-lyric/util/log"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"sync"
)

type FwRedis struct {
	client redis.UniversalClient
}

var redisCache *FwRedis
var onceRedis sync.Once

func GetRedisCache() redis.UniversalClient {
	return redisCache.client
}

func InitRedis() {
	onceRedis.Do(func() {
		redisCache = new(FwRedis)
		redisCache.connectDB(config.Instance.RedisCache)
	})
}

func (r *FwRedis) connectDB(conf config.RedisConfig) {
	host := conf.Host

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Host,
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("redis connect ping failed, err:", err)
	} else {
		log.Debug("redis connect ping response:", pong)
		log.Info("redis connection established:", host)
		r.client = client
	}
}
