package groot

import (
	"fmt"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var client *goredislib.Client
var rs *redsync.Redsync

func InitRedis() error {
	if redisCfg.Addr == "" {
		fmt.Println("没有配置redis信息！")
		Debug("没有配置redis信息!")
		return nil
	}

	client = goredislib.NewClient(&goredislib.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.Db,       // use default DB
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
		if client == nil {
			InitRedis()
		}
	}
	if client == nil {
		panic("redis配置不正确为空！")
	}
	return client
}

func GetRedisSync() *redsync.Redsync {
	return rs
}
