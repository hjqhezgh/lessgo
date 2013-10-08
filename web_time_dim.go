// Title：时间维度web服务
//
// Description:时间维度web服务
//
// Author:black
//
// Createtime:2013-08-19 17:48
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-19 17:48 black 创建文档
package lessgo

import (
	"github.com/hjqhezgh/commonlib"
	"net/http"
)

func years(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	years, err := FindYear()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	m["datas"] = years

	commonlib.OutputJson(w, m, " ")

	return
}

func months(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	year := r.FormValue("year")

	months, err := FindMonth(year)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	m["datas"] = months

	commonlib.OutputJson(w, m, " ")

	return
}

func weeks(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	year := r.FormValue("year")
	month := r.FormValue("month")

	weeks, err := FindWeek(year, month)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	m["datas"] = weeks

	commonlib.OutputJson(w, m, " ")

	return
}
