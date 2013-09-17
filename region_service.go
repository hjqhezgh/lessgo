// Title：地区相关的服务
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
	"database/sql"
)

// 根据父节点地区码获取该节点下的地区列表
func FindRegionByParentCode(parentCode string) (regions []*Region, err error) {
	db := DBPool{}.getMySQL()
	defer db.Close()

	regions = []*Region{}

	var rows *sql.Rows

	if parentCode == "" { //根节点的情况
		rows, err = db.Query("SELECT code,name FROM region WHERE parent_code is null")
	} else {
		rows, err = db.Query("SELECT code,name FROM region WHERE parent_code = ?", parentCode)
	}

	if err != nil {
		Log.Error(err.Error())
		return regions, err
	}

	for rows.Next() {
		region := new(Region)
		err := rows.Scan(&region.Code, &region.Name)
		if err != nil {
			Log.Warn(err.Error())
			return regions, err
		}
		regions = append(regions, region)
	}

	return regions, nil
}

//根据地区号，获取地区全称
func GetAreaNameByCode(area string) (areaName string) {
	if area == "" || len(area) != 6 {
		return ""
	}

	db := DBPool{}.getMySQL()
	defer db.Close()

	province := area[0:2] + "0000"
	city := area[0:4] + "00"

	var tmp string

	//查询省份名字
	rows, err := db.Query("SELECT name FROM region WHERE code = ?", province)
	if err != nil {
		return ""
	}

	if rows.Next() {
		rows.Scan(&tmp)
		areaName = tmp
	}

	if province == city {
		return areaName
	}

	//查询城市名字
	rows, err = db.Query("SELECT name FROM region WHERE code = ?", city)
	if err != nil {
		return ""
	}

	if rows.Next() {
		rows.Scan(&tmp)
		areaName = areaName + tmp
	}

	if area == city {
		return areaName
	}

	//查询区的名字
	rows, err = db.Query("SELECT name FROM region WHERE code = ?", area)
	if err != nil {
		return ""
	}

	if rows.Next() {
		rows.Scan(&tmp)
		areaName = areaName + tmp
	}

	return areaName
}

//根据地区号，获取地区名称
func GetOneAreaNameByCode(area string) (areaName string) {
	if area == "" || len(area) != 6 {
		return ""
	}

	db := DBPool{}.getMySQL()
	defer db.Close()

	//查询区的名字
	rows, err := db.Query("SELECT name FROM region WHERE code = ?", area)
	if err != nil {
		return ""
	}

	tmp := ""

	if rows.Next() {
		rows.Scan(&tmp)
		areaName = areaName + tmp
	}

	return areaName
}
