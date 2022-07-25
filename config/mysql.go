/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 12:21
 */

package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MasterDB *sqlx.DB
var dns string
// mysql装在宿主机 通过ifconfig查看容器的ip
var MySQL = map[string]interface{}{
	"user":     "root",
	"password": "",
	"charset":  "utf8mb4",
	"host":     "172.17.0.1",
	"port":     "3306",
	"dbname":   "warlock",
}

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

func fillDns(mysqlConfig map[string]interface{}) {
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"])
}

func Init() error {
	var err error
	fillDns(MySQL)

	// 启动时就打开数据库连接
	if err = initEngine(); err != nil {
		fmt.Println("mysql is not open:", err)
		return err
	}
	return nil
}

func initEngine() error {
	var err error
	// 也可以使用MustConnect连接不成功就panic
	fillDns(MySQL)
	MasterDB, err = sqlx.Connect("mysql", dns)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return err
	}
	MasterDB.SetMaxOpenConns(20)
	MasterDB.SetMaxIdleConns(10)

	return nil
}
