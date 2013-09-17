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
	_ "github.com/go-sql-driver/mysql"
)

type DBPool struct {
}

var MySQLPool chan *sql.DB

func (dbPool DBPool) getMySQL() *sql.DB {
	if MySQLPool == nil {
		MySQLPool = make(chan *sql.DB, config.MaxPoolSize)
	}
	if len(MySQLPool) == 0 {
		go func() {
			for i := 0; i < config.MaxPoolSize/2; i++ {
				db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", config.DbUserName, config.DbPwd, config.DbUrl, config.DbName))
				if err != nil {
					Log.Warn(err)
					continue
				}
				dbPool.putMySQL(db)
			}
		}()
	}
	return <-MySQLPool
}

func (dbPool DBPool) putMySQL(conn *sql.DB) {
	if MySQLPool == nil {
		MySQLPool = make(chan *sql.DB, config.MaxPoolSize)
	}
	if len(MySQLPool) == config.MaxPoolSize {
		conn.Close()
		return
	}
	MySQLPool <- conn
}
