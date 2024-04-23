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

var (
	dbSlave    *gorm.DB
	sqlDBSlave *sql.DB
)

func DBSlave() *gorm.DB {
	return dbSlave
}

func initSlave() {
	gormConf := &gorm.Config{}
	gormConf.Logger = logger.New(milog.New(log.GetLogger().GetWriter(), "\r\n[db]", log.LstdFlags),
		logger.Config{
			SlowThreshold:             3 * time.Second,
			LogLevel:                  If(config.Instance.ShowSql, logger.Info, logger.Warn).(logger.LogLevel),
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		})
	err := openSlaveDB(config.Instance.SlaveMySqlUrl, gormConf,
		config.Instance.SlaveMySqlMaxIdle, config.Instance.SlaveMySqlMaxOpen)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	log.Info("Slave MySQL connection established")
}

func openSlaveDB(dsn string, config *gorm.Config, maxIdleConns, maxOpenConns int) (err error) {
	if config == nil {
		config = &gorm.Config{}
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		}
	}

	if dbSlave, err = gorm.Open(mysql.Open(dsn), config); err != nil {
		log.Errorf("opens slave database failed: %v", err.Error())
		return
	}

	if sqlDBSlave, err = db.DB(); err == nil {
		sqlDBSlave.SetMaxIdleConns(maxIdleConns)
		sqlDBSlave.SetMaxOpenConns(maxOpenConns)
	} else {
		log.Error(err)
	}

	return
}
