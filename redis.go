package groot

import (
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
)

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

var client *goredislib.Client
var rs *redsync.Redsync

func Init(keys ...string) error {
	addr := "127.0.0.1"
	password := ""
	db := 0
	if len(keys) == 0 {
		var cfg RedisConfig
		err := UnmarshalConfigByKey("redis", &cfg)
		if err != nil {
			Error("redis init error", zap.String("err", err.Error()))
			return err
		}
		addr = cfg.Addr
		password = cfg.Password
		db = cfg.Db
	} else if len(keys) >= 3 {
		addr = keys[0]
		password = keys[1]
		db = 0
	} else if len(keys) == 2 {

	} else {

	}
	client = goredislib.NewClient(&goredislib.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs = redsync.New(pool)
	return nil
}

func GetRedisClient() *goredislib.Client {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()
		Init()
	}
	return client
}

func GetRedisSync() *redsync.Redsync {
	return rs
}
