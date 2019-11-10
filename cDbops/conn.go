package cDbops

import (
	"AleCode/clog"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	dbConn *sql.DB
	err    error
)

func Init() {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, addrs, database)
	dbConn, err = sql.Open("mysql", dsn)
	if err != nil {
		clog.Error("get mysql conn failed,err:%v", err)
		_ = dbConn.Close()
		return
	}
	//设置连接的最大连接周期，超时自动关闭
	dbConn.SetConnMaxLifetime(connMaxLifetime * time.Second)
	//设置最大连接数
	dbConn.SetMaxOpenConns(maxOpenConns)
	//设置闲置连接数
	dbConn.SetMaxIdleConns(maxIdleConns)
}

func GetDb() *sql.DB {
	return dbConn
}
