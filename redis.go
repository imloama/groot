package groot

import (
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var client *goredislib.Client
var rs *redsync.Redsync

func InitRedis() error {
	err := LoadConfig()
	if err != nil {
		return err
	}
	if serverCfg.Redis == nil {
		Debug("没有配置redis信息！")
		return nil
	}

	client = goredislib.NewClient(&goredislib.Options{
		Addr:     serverCfg.Redis.Addr,
		Password: serverCfg.Redis.Password, // no password set
		DB:       serverCfg.Redis.Db,       // use default DB
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
		InitRedis()
	}
	return client
}

func GetRedisSync() *redsync.Redsync {
	return rs
}
