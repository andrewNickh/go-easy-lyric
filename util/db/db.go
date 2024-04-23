package db

import (
	"database/sql"
	"easy-lyric/config"
	"easy-lyric/util/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	milog "log"
	"time"
)

func InitDB() {
	initMaster()
	initSlave()
}

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

func Master() *gorm.DB {
	return db
}

func initMaster() {
	gormConf := &gorm.Config{}
	gormConf.Logger = logger.New(milog.New(log.GetLogger().GetWriter(), "\r\n[db]", log.LstdFlags),
		logger.Config{
			SlowThreshold:             3 * time.Second,
			LogLevel:                  If(config.Instance.ShowSql, logger.Info, logger.Warn).(logger.LogLevel),
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		})
	err := openDB(config.Instance.MySqlUrl, gormConf,
		config.Instance.MySqlMaxIdle, config.Instance.MySqlMaxOpen)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Info("MySQL connection established")
}

func openDB(dsn string, config *gorm.Config, maxIdleConns, maxOpenConns int) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		}
	}

	if db, err = gorm.Open(mysql.Open(dsn), config); err != nil {
		log.Infof("opens database failed: %v", err.Error())
		return
	}

	if sqlDB, err = db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	} else {
		log.Error(err)
	}
	return
}

func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}
