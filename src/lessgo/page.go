// Title：分页工具类
//
// Description:
//
// Author:black
//
// Createtime:2013-08-06 13:04
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-06 13:04 black 创建文档
package lessgo

import (
	"math"
)

//传统分页
type TraditionPage struct {
	CurrPageNo int      `json:currPageNo` //当前页
	PrePageNo  int      `json:prePageNo`  //上一页
	NextPageNo int      `json:nextPageNo` //下一页
	TotalPage  int      `json:totalPage`  //总共几页
	TotalNum   int      `json:totalNum`   //总共几条
	PageSize   int      `json:pageSize`   //每页几条
	HasNext    bool     `json:hasNext`    //是否有下一页
	HasPre     bool     `json:hasPre`     //是否有上一页
	Datas      []*Model `json:datas`      //数据集
}

//构建传统式Page对象
func BulidTraditionPage(currPageNo, pageSize, totalNum int, datas []*Model) (pageData *TraditionPage) {

	pageData = new(TraditionPage)

	pageData.CurrPageNo = currPageNo

	pageData.PageSize = pageSize

	pageData.TotalNum = totalNum

	pageData.TotalPage = int(math.Ceil(float64(totalNum) / float64(pageSize)))

	if currPageNo < pageData.TotalPage {
		pageData.HasNext = true
		pageData.NextPageNo = currPageNo + 1
	} else {
		pageData.HasNext = false
	}

	if currPageNo > 1 {
		pageData.HasPre = true
		pageData.PrePageNo = currPageNo - 1
	} else {
		pageData.HasPre = false
	}

	pageData.Datas = datas

	return pageData
}
