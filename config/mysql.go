/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 12:21
 */

package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var MasterDB *sqlx.DB

func init() {
	var err error
	// 启动时就打开数据库连接
	if err = initEngine(); err != nil {
		panic(err)
	}

	// 测试数据库连接是否 OK
	if err = MasterDB.Ping(); err != nil {
		panic(err)
	}
}

func initEngine() error {
	var err error
	// 也可以使用MustConnect连接不成功就panic
	MasterDB, err = sqlx.Open("sqlite3", "file:./db/db.sqlite3?mode=rwc")
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return err
	}
	MasterDB.SetMaxOpenConns(20)
	MasterDB.SetMaxIdleConns(10)

	return nil
}
