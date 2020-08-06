package database

import (
	"github.com/aceld/zinx/zlog"
	"zinx-mj/database/dbface"
	"zinx-mj/database/mongodb"
)

/*
 * Descrp: 获取数据库接口
 * Create: zhangyi 2020-06-20 02:00:34
 */
var db dbface.IDataBase

func init() {
	db = mongodb.New()
}

func GetDB() dbface.IDataBase {
	return db
}

const (
	DBADDR = "localhost"
	DBPORT = 27017
	//DBPORT = 9191
)

func Connect() error {
	if err := db.Connect(DBADDR, DBPORT); err != nil {
		zlog.Error("Init, connect error")
		return err
	}
	if err := db.Ping(); err != nil {
		zlog.Error("Init, ping error")
		return err
	}
	zlog.Debug("connect mongo success")
	return nil
}
