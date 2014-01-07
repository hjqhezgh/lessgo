// Title：
//
// Description:
//
// Author:black
//
// Createtime:2014-01-06 18:01
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2014-01-06 18:01 black 创建文档
package lessgo

import (
	"github.com/hjqhezgh/commonlib"
	"database/sql"
	"math"
)

//获取数据库总页数
func GetTotalPage(pageSize int,db *sql.DB,sql string ,params []interface {}) (totalPage,totalNum int,err error){

	rows, err := db.Query(sql, params...)

	if err != nil {
		Log.Error(err.Error())
		return 0,0,err
	}

	if rows.Next() {
		err := rows.Scan(&totalNum)

		if err != nil {
			Log.Error(err.Error())
			return 0,0,err
		}
	}

	totalPage = int(math.Ceil(float64(totalNum) / float64(pageSize)))

	return totalPage,totalNum,nil
}

func GetFillObjectPage(db *sql.DB,sql string,currPageNo,pageSize,totalNum int ,params []interface {}) (*commonlib.TraditionPage,error){

	rows, err := db.Query(sql, params...)

	if err != nil {
		Log.Error(err.Error())
		return nil,err
	}

	objects := []interface{}{}

	columns,err := rows.Columns()

	if err != nil {
		Log.Error(err.Error())
		return nil,err
	}

	for rows.Next() {

		model := new(Model)

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for index,column := range columns {
			if index > 0 {//第一个列必须是id
				prop := new(Prop)
				prop.Name = column
				prop.Value = ""
				fillObjects = append(fillObjects, &prop.Value)
				model.Props = append(model.Props, prop)
			}
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			Log.Error(err.Error())
			return nil,err
		}

		objects = append(objects, model)
	}


	return commonlib.BulidTraditionPage(currPageNo, pageSize, totalNum, objects),nil
}
