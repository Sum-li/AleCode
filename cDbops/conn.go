package cDbops

import (
	"AleCode/clog"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var dbConn *sql.DB

func newConn(username,password,server,database string) (conn *sql.DB,err error) {
	var dsn string
	dsn = username+":"+password+"@tcp("+server+")/"+database+"?charset=utf8&parseTime=True&loc=Local"
	conn,err = sql.Open("mysql",dsn)
	if err != nil {
		clog.Error("get mysql conn failed,err:%v",err)
		_ = conn.Close()
		return
	}
	return
}

func init() {
	var err error
	dbConn,err =  newConn("user","pswd","127.0.0.1:3306","test")
	if err != nil {
		clog.Error("get mysql conn failed,err:%v",err)
		panic(err.Error())
	}
	//设置连接的最大连接周期，超时自动关闭
	dbConn.SetConnMaxLifetime(100 * time.Second)
	//设置最大连接数
	dbConn.SetMaxOpenConns(100)
	//设置闲置连接数
	dbConn.SetMaxIdleConns(15)
}
