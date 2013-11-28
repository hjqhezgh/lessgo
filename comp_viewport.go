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
	Window           string            `xml:"window,attr"`
	GridPanels       []gridPanel       `xml:"gridpanel"`
	FormPanels       []formPanel       `xml:"formpanel"`
	TabFormPanels    []tabFormPanel    `xml:"tabformpanel"`
	CustomGridPanels []customGridPanel `xml:"customgridpanel"`
	CustomFormPanels []customFormPanel `xml:"customformpanel"`
	BlankPanels      []blankPanel      `xml:"blankpanel"`
	Crumbs           crumbs            `xml:"crumbs"`
	CustomJs         string            `xml:"customJs"`
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

	for _, tabFormPanel := range viewport.TabFormPanels {
		content += string(tabFormPanel.generate(terminal, packageName, r))
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

	data := make(map[string]interface{})
	data["Content"] = content
	data["Crumbs"] = viewport.Crumbs
	data["Crumbs"] = viewport.Crumbs
	data["Employee"] = employee
	data["SiteName"] = SiteName
	data["CustomJs"] = viewport.CustomJs

	var t *template.Template

	var buf bytes.Buffer

	tempName := ""

	if viewport.Window == "true" {
		t = template.New("window.html")
		data["ParentComponentId"] = r.FormValue("parentComponentId")
		data["ParentWindowName"] = r.FormValue("parentWindowName")
		tempName = "../lessgo/template/component/" + terminal + "/window.html"
	} else {
		t = template.New("viewport.html")
		tempName = "../lessgo/template/component/" + terminal + "/viewport.html"
	}

	t = t.Funcs(template.FuncMap{
		"compareString": CompareString,
	})

	t, err := t.ParseFiles(tempName)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()
}
