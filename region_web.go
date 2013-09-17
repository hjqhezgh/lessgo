// Title：地区相关的web服务
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
	"net/http"
	"github.com/hjqhezgh/commonlib"
)

func Regions(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil {
		m["success"] = false
		m["reason"] = "请求解析异常"
		commonlib.OutputJson(w, m ," ")
		return
	}

	code := r.FormValue("code")

	regions, err := FindRegionByParentCode(code)

	if err != nil {
		m["success"] = false
		m["reason"] = "服务器异常"
		commonlib.OutputJson(w, m," ")
		return
	}

	m["success"] = true
	m["regions"] = regions

	commonlib.OutputJson(w, m," ")

	return
}
