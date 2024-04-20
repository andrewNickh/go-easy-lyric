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
	"os"
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
	newLogger := logger.New(
		milog.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	gormConf.Logger = newLogger
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
