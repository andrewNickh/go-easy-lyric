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
	newLogger := logger.New(
		milog.New(log.GetLogger().GetWriter(), "\r\n[db]", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	gormConf.Logger = newLogger
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

	log.Info("连接mysql服务(从库)成功")
	return
}
