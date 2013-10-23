// Title：全局容器组件
//
// Description: 是一个view的根元素，是必需的，所有的元素都要建立在该容器之下
//
// Author:black
//
// Createtime:2013-08-08 09:27
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-08 09:27 black 创建文档
package lessgo

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"text/template"
)

//当viewport下面可以放置其他元素的时候，就扩展Viewport结构体
type viewport struct {
	XMLName          xml.Name          `xml:"viewport"`
	GridPanels       []gridPanel       `xml:"gridpanel"`
	FormPanels       []formPanel       `xml:"formpanel"`
	MutiFormPanels   []mutiFormPanel   `xml:"mutiformpanel"`
	CustomGridPanels []customGridPanel `xml:"customgridpanel"`
	CustomFormPanels []customFormPanel `xml:"customformpanel"`
	BlankPanels      []blankPanel      `xml:"blankpanel"`
	Crumbs           crumbs            `xml:"crumbs"`
}

//面包屑
type crumbs struct {
	Crumbs []crumb `xml:"crumb"`
}

type crumb struct {
	Url         string `xml:"url,attr"`
	Text        string `xml:"text,attr"`
	CurrentPage string `xml:"currentPage,attr"`
}

//扩展viewport的同时，记得同时扩展container
func (viewport viewport) generateViewport(terminal, packageName string, r *http.Request, employee Employee) []byte {

	content := ""

	for _, formpanel := range viewport.FormPanels {
		content += string(formpanel.generate(getEntity(formpanel.Entity), terminal, packageName, r))
	}

	for _, mutiformpanel := range viewport.MutiFormPanels {
		content += string(mutiformpanel.generate(terminal, packageName, r))
	}

	for _, gridpanel := range viewport.GridPanels {
		content += string(gridpanel.generate(getEntity(gridpanel.Entity), terminal, packageName))
	}

	for _, customgridpanel := range viewport.CustomGridPanels {
		content += string(customgridpanel.generate(terminal, packageName))
	}

	for _, customformpanel := range viewport.CustomFormPanels {
		content += string(customformpanel.generate(terminal, packageName))
	}

	for _, blankpanel := range viewport.BlankPanels {
		content += string(blankpanel.generate(terminal, packageName))
	}

	var t *template.Template

	var buf bytes.Buffer

	t = template.New("viewport.html")

	t = t.Funcs(template.FuncMap{
		"compareString": CompareString,
	})

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/viewport.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})
	data["Content"] = content
	data["Crumbs"] = viewport.Crumbs
	data["Employee"] = employee
	data["SiteName"] = SiteName
	data["SiteIcon"] = SiteIcon

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()
}
