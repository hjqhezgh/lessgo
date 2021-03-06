// Title：时间维度服务
//
// Description:
//
// Author:black
//
// Createtime:2013-08-19 17:29
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-19 17:29 black 创建文档
package lessgo

import (
	"github.com/hjqhezgh/commonlib"
)

type Week struct {
	WeekKey     string `json:"weekKey"`
	CurrentWeek string `json:"currentWeek"`
}

//返回系统支持的年份
func FindYear() ([]string, error) {

	db := GetMySQL()

	defer db.Close()

	years := []string{}

	rows, err := db.Query("select distinct(current_year) year from time_dim order by year")

	if err != nil {
		Log.Error(err.Error())
		return years, err
	}

	for rows.Next() {
		year := ""
		err := rows.Scan(&year)
		if err != nil {
			Log.Warn(err.Error())
			return []string{}, err
		}
		years = append(years, year)
	}

	return years, nil
}

//返回指定年份下的月份
func FindMonth(year string) ([]string, error) {

	db := GetMySQL()

	defer db.Close()

	months := []string{}

	rows, err := db.Query("select distinct(current_month) month from time_dim where current_year=? order by month", year)

	if err != nil {
		Log.Error(err.Error())
		return months, err
	}

	for rows.Next() {
		month := ""
		err := rows.Scan(&month)
		if err != nil {
			Log.Warn(err.Error())
			return []string{}, err
		}
		months = append(months, month)
	}

	return months, nil
}

//返回指定年份、月份下的周信息
func FindWeek(year, month string) ([]*Week, error) {
	db := GetMySQL()

	defer db.Close()

	weeks := []*Week{}

	rows, err := db.Query("select distinct(week_key) week_key,week_of_month from time_dim where current_year=? and current_month=? order by week_key", year, month)

	if err != nil {
		Log.Error(err.Error())
		return weeks, err
	}

	for rows.Next() {
		week := new(Week)
		err := rows.Scan(&week.WeekKey, &week.CurrentWeek)
		if err != nil {
			Log.Warn(err.Error())
			return []*Week{}, err
		}
		weeks = append(weeks, week)
	}

	return weeks, nil
}

//根据时间key值获取时间区间，false表示没有找到相应的时间区间
func FindRangeTimeDim(yearKey, monthKey, weekKey string) (startTime, endTime string, flag bool) {

	db := GetMySQL()

	defer db.Close()

	sql := ""

	param := ""

	if weekKey != "" {
		sql = "select min(day_date),max(day_date) from time_dim where week_key=?"
		param = weekKey
	} else if monthKey != "" {
		sql = "select min(day_date),max(day_date) from time_dim where month_key=?"
		param = monthKey
	} else {
		sql = "select min(day_date),max(day_date) from time_dim where current_year=?"
		param = yearKey
	}

	Log.Debug(sql)

	rows, err := db.Query(sql, param)

	if err != nil {
		Log.Error(err.Error())
		return "", "", false
	}

	if rows.Next() {
		err := commonlib.PutRecord(rows, &startTime, &endTime)
		if err != nil {
			Log.Warn(err.Error())
			return "", "", false
		}

		if startTime != "" && endTime != "" {
			return startTime, endTime, true
		}

	}

	return "", "", false
}
