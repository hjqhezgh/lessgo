// Title：统一的数据库服务
//
// Description:
//
// Author:black
//
// Createtime:2013-08-06 17:13
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-06 17:13 black 创建文档
package lessgo

import (
	"fmt"
	"math"
	"strings"
	"github.com/hjqhezgh/commonlib"
)

//构建查询的sql语句
func bulidSelectSql(entity Entity, colmuns []column) (string, string) {

	dataSql := "select " + entity.Id + "." + entity.Pk + ","

	if len(colmuns) > 0 {
		for _, column := range colmuns {
			if strings.Contains(column.Field, ".") {
				dataSql += column.Field + ","
			} else {
				dataSql += entity.Id + "." + column.Field + ","
			}
		}
	} else {
		for _, field := range entity.Fields {
			dataSql += entity.Id + "." + field.Name + ","
		}

		for _, ref := range entity.Refs {

			for _, field := range ref.Fields {
				dataSql += ref.Entity + "." + field.Name + ","
			}

		}
	}

	dataSql = commonlib.Substr(dataSql, 0, len(dataSql)-1)

	dataSql += " from " + entity.Id + " " + entity.Id

	for _, ref := range entity.Refs {
		dataSql += fmt.Sprintf(" left join %v %v on %v.%v=%v.%v", ref.Entity, ref.Entity, entity.Id, ref.Field, ref.Entity, ref.RefEntityField)
	}

	countSql := "select count(1) from " + entity.Id

	return countSql, dataSql
}

//构建查询
func bulidWhereSql(entity Entity, countSql, dataSql string, searchParam []search) (string, string, []interface{}) {

	params := []interface{}{}

	if len(searchParam) > 0 {
		countSql += " where 1=1 "
		dataSql += " where 1=1 "

		for _, search := range searchParam {
			if search.Value != "" {

				opera := ""

				switch search.SearchType {
				case "like":
					opera = "like"
				case "eq":
					opera = "="
				case "ge":
					opera = ">="
				case "le":
					opera = "<="
				case "gt":
					opera = ">"
				case "lt":
					opera = "<"
				}

				countSql += " and " + entity.Id + "." + search.Field + " " + opera + " ? "
				dataSql += " and " + entity.Id + "." + search.Field + " " + opera + " ? "

				if search.SearchType == "like" {
					params = append(params, "%"+search.Value+"%")
				} else {
					params = append(params, search.Value)
				}
			}
		}

		return countSql, dataSql, params
	} else {
		return countSql, dataSql, []interface{}{}
	}

	return countSql, dataSql, []interface{}{}
}

//查找分页数据
func findTraditionPage(entity Entity, currPageNo, pageSize int, searchParam []search, colmuns []column) (*commonlib.TraditionPage, error) {

	countSql, dataSql := bulidSelectSql(entity, colmuns)

	countSql, dataSql, params := bulidWhereSql(entity, countSql, dataSql, searchParam)

	Log.Debug(countSql)

	db := DBPool{}.GetMySQL()
	defer db.Close()

	rows, err := db.Query(countSql, params...)

	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	totalNum := 0

	if rows.Next() {
		err := rows.Scan(&totalNum)

		if err != nil {
			Log.Error(err.Error())
			return nil, err
		}
	}

	dataSql += " order by " + entity.Pk + " desc "
	dataSql += " limit ?,?"

	totalPage := int(math.Ceil(float64(totalNum) / float64(pageSize)))

	if currPageNo > totalPage {
		currPageNo = totalPage
	}

	params = append(params, (currPageNo-1)*pageSize, pageSize)
	Log.Debug(dataSql)

	rows, err = db.Query(dataSql, params...)

	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	objects := []interface {}{}

	for rows.Next() {

		model := new(Model)
		model.Entity = entity
		model.Id = 0
		model.Props = []*Prop{}

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for _, column := range colmuns {
			prop := new(Prop)
			prop.Name = column.Field
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			Log.Error(err.Error())
			return nil, err
		}

		objects = append(objects, model)
	}

	return commonlib.BulidTraditionPage(currPageNo, pageSize, totalNum, objects), nil
}

//查找分页数据
func findAllData(entity Entity) ([]*Model, error) {

	_, dataSql := bulidSelectSql(entity, []column{})

	Log.Debug(dataSql)

	db := DBPool{}.GetMySQL()
	defer db.Close()

	dataSql += " order by " + entity.Pk + " desc "

	Log.Debug(dataSql)

	objects := []*Model{}

	rows, err := db.Query(dataSql)

	if err != nil {
		Log.Error(err.Error())
		return objects, err
	}

	for rows.Next() {

		model := new(Model)
		model.Entity = entity
		model.Id = 0
		model.Props = []*Prop{}

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for _, field := range entity.Fields {
			prop := new(Prop)
			prop.Name = field.Name
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		for _, ref := range entity.Refs {
			for _, field := range ref.Fields {
				prop := new(Prop)
				prop.Name = ref.Entity + "." + field.Name
				prop.Value = ""
				fillObjects = append(fillObjects, &prop.Value)
				model.Props = append(model.Props, prop)
			}
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			Log.Error(err.Error())
			return []*Model{}, err
		}

		objects = append(objects, model)
	}

	return objects, nil
}

//添加数据
func insert(entity Entity, model *Model, elements []element) (id int, err error) {

	sql := "insert into " + entity.Id + "("
	valueSql := " values ("

	params := []interface{}{}

	for index, element := range elements {
		sql += element.Field

		valueSql += "?"

		if index < len(elements)-1 {
			sql += ","
			valueSql += ","
		}

		params = append(params, getPropValue(model, element.Field))
	}
	fmt.Println(params)
	sql += ")"
	valueSql += ")"
	sql += valueSql

	Log.Debug(sql)

	db := DBPool{}.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sql)

	if err != nil {
		Log.Warn(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(params...)

	if err != nil {
		Log.Warn(err.Error())
		return 0, err
	}

	id64, err := res.LastInsertId()

	if err != nil {
		Log.Warn(err.Error())

		return 0, err
	}

	return int(id64), nil
}

//根据id查找对象
func findById(entity Entity, id string) (*Model, error) {

	dataSql := "select " + entity.Pk + ","

	for index, field := range entity.Fields {
		dataSql += field.Name
		if index < len(entity.Fields)-1 {
			dataSql += ","
		}
	}

	dataSql += " from " + entity.Id + " where " + entity.Pk + "=?"

	Log.Debug(dataSql)

	db := DBPool{}.GetMySQL()
	defer db.Close()

	rows, err := db.Query(dataSql, id)

	if err != nil {
		Log.Warn(err.Error())
		return nil, err
	}

	model := new(Model)
	model.Entity = entity
	model.Id = 0
	model.Props = []*Prop{}

	if rows.Next() {

		fillObjects := []interface{}{}

		fillObjects = append(fillObjects, &model.Id)

		for _, field := range entity.Fields {
			prop := new(Prop)
			prop.Name = field.Name
			prop.Value = ""
			fillObjects = append(fillObjects, &prop.Value)
			model.Props = append(model.Props, prop)
		}

		err = commonlib.PutRecord(rows, fillObjects...)

		if err != nil {
			return nil, err
		}
	}

	return model, nil
}

//修改
func modify(entity Entity, model *Model, elements []element) error {
	sql := "update " + entity.Id + " set "

	params := []interface{}{}

	for index, element := range elements {
		sql += element.Field + "=?"

		if index < len(elements)-1 {
			sql += ","
		}

		params = append(params, getPropValue(model, element.Field))
	}

	sql += " where " + entity.Pk + "=?"
	params = append(params, model.Id)

	Log.Debug(sql)

	db := DBPool{}.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sql)
	if err != nil {
		Log.Warn(err.Error())
		return err
	}

	_, err = stmt.Exec(params...)

	if err != nil {
		Log.Warn(err.Error())
		return err
	}

	return nil
}

//删除
func delete(entity Entity, id string) error {
	sql := "delete from " + entity.Id + " where " + entity.Pk + "=?"

	Log.Debug(sql)

	db := DBPool{}.GetMySQL()
	defer db.Close()

	stmt, err := db.Prepare(sql)
	if err != nil {
		Log.Warn(err.Error())
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		Log.Warn(err.Error())
		return err
	}

	return nil
}
