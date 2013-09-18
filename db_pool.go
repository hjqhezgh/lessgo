// Title：数据库连接池
//
// Description:
//
// Author:Bill Cai
//
// Createtime:2013-08-06 14:15
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-07-05 09:53 Bill 创建文档
package lessgo

import (
	"database/sql"
	"fmt"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

var mySQLPool chan *sql.DB

func GetMySQL() *sql.DB {

	maxPoolSizeString,_ := Config.GetValue("lessgo","maxPoolSize")
	maxPoolSize,_ := strconv.Atoi(maxPoolSizeString)

	if mySQLPool == nil {
		mySQLPool = make(chan *sql.DB, maxPoolSize)
	}

	dbUrl,_ := Config.GetValue("lessgo","dbUrl")
	dbName,_ := Config.GetValue("lessgo","dbName")
	dbUserName,_ := Config.GetValue("lessgo","dbUserName")
	dbPwd,_ := Config.GetValue("lessgo","dbPwd")

	if len(mySQLPool) == 0 {
		go func() {
			for i := 0; i < maxPoolSize/2; i++ {
				db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", dbUserName, dbPwd, dbUrl, dbName))
				if err != nil {
					Log.Warn(err)
					continue
				}
				putMySQL(db)
			}
		}()
	}
	return <-mySQLPool
}

func putMySQL(conn *sql.DB) {

	maxPoolSizeString,_ := Config.GetValue("lessgo","maxPoolSize")
	maxPoolSize,_ := strconv.Atoi(maxPoolSizeString)

	if mySQLPool == nil {
		mySQLPool = make(chan *sql.DB, maxPoolSize)
	}

	if len(mySQLPool) == maxPoolSize {
		conn.Close()
		return
	}

	mySQLPool <- conn
}
