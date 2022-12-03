package groot

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

type OrmConfig struct {
	Driver          string `json:"driver"`
	Url             string `json:"url"`
	Debug           bool   `json:"debug"`
	MaxIdleConns    int    `json:"maxIdleConns"`    // SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	MaxOpenConns    int    `json:"maxOpenConns"`    // SetMaxOpenConns 设置打开数据库连接的最大数量。
	ConnMaxLifetime int    `json:"connMaxLifetime"` // SetConnMaxLifetime 设置了连接可复用的最大时间， 秒
}

const (
	DRIVER_MYSQL      = "mysql"
	DRIVER_POSTGRESQL = "pgsql"
	DRIVER_SQLITE     = "sqlite"
)

/**
mysql: "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
pg: "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
*/

var db *gorm.DB
var dbinstlock sync.Mutex

func InitDb() {
	if ormCfg.Url == "" || ormCfg.Driver == "" {
		return
	}
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks
	if ormCfg.Debug {
		logger.LogLevel = gormlogger.Info
	}
	if DRIVER_MYSQL == ormCfg.Driver {
		dbinst, err := gorm.Open(mysql.Open(ormCfg.Url), &gorm.Config{Logger: logger})
		if err != nil {
			panic(err)
		}
		db = dbinst
		initPool()
		return
	}
	if DRIVER_POSTGRESQL == ormCfg.Driver {
		dbinst, err := gorm.Open(postgres.New(postgres.Config{
			DSN: ormCfg.Url,
			// PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{Logger: logger})
		if err != nil {
			panic(err)
		}
		db = dbinst
		initPool()
		return
	}
	if DRIVER_SQLITE == ormCfg.Driver {
		dbinst, err := gorm.Open(sqlite.Open(ormCfg.Url), &gorm.Config{Logger: logger})
		if err != nil {
			panic(err)
		}
		db = dbinst
		initPool()
		return
	}
}

func initPool() {
	if db == nil {
		return
	}
	if ormCfg.MaxIdleConns <= 0 {
		ormCfg.MaxIdleConns = 10
	}
	if ormCfg.MaxOpenConns <= 0 {
		ormCfg.MaxOpenConns = 100
	}
	if ormCfg.ConnMaxLifetime <= 0 {
		ormCfg.ConnMaxLifetime = 3600
	}
	d, err := db.DB()
	if err != nil {
		panic(err)
	}
	d.SetMaxIdleConns(ormCfg.MaxIdleConns)
	d.SetMaxOpenConns(ormCfg.MaxOpenConns)
	d.SetConnMaxLifetime(time.Second * time.Duration(ormCfg.ConnMaxLifetime))
}

func GetDb() *gorm.DB {
	if db == nil {
		dbinstlock.Lock()
		defer dbinstlock.Unlock()
		if db == nil {
			InitDb()
		}
	}
	return db
}
