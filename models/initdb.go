package models

import (
	"database/sql"
	"sync"
	"time"
	"walksmile/conf"
)

var db *sql.DB
var onceDB sync.Once

func InitDB() {
	db, _ = sql.Open("mysql", conf.Conf.Mysql.DataSource)
	db.SetConnMaxLifetime(time.Duration(conf.Conf.Mysql.ConnMaxLifetime))
	db.SetMaxIdleConns(conf.Conf.Mysql.MaxIdleConns)
	db.SetMaxOpenConns(conf.Conf.Mysql.MaxOpenConns)

}

func GetDB() *sql.DB {
	if db == nil {
		onceDB.Do(InitDB)
	}
	return db

}
