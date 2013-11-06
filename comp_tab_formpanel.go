// Title：多表编辑表单
//
// Description: 可以同时录入多表数据的表单控件，目前只支持添加
//
// Author:black
//
// Createtime:2013-08-19 11:55
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-19 11:55 black 创建文档
package lessgo

import (
	"bytes"
	"net/http"
	"text/template"
)

type tabFormPanel struct {
	Id              string       `xml:"id,attr"`
	FormTabs        []formTab    `xml:"formtab"`
	Load            string       `xml:"load,attr"`
	LoadUrl         string       `xml:"loadUrl,attr"`
	SaveUrl         string       `xml:"saveUrl,attr"`
	Title           string       `xml:"title,attr"`
	FormButtons     []formButton `xml:"formButton"`
	BeforeSave      string       `xml:"beforeSave"`
	AfterRender     string       `xml:"afterRender"`
	AfterSave       string       `xml:"afterSave"`
	Inwindow        string       `xml:"inwindow,attr"`
	HideSaveButton  string       `xml:"hideSaveButton,attr"`
	HideResetButton string       `xml:"hideResetButton,attr"`
}

type formTab struct {
	Elements []element `xml:"element"`
	Desc     string    `xml:"desc,attr"`
}

func (tabFormPanel tabFormPanel) generate(terminal, packageName string, r *http.Request) []byte {

	var t *template.Template

	var buf bytes.Buffer

	tabFormPanel.Id = packageName + "." + tabFormPanel.Id

	runtimeComponentContain[tabFormPanel.Id] = tabFormPanel

	t = template.New("tabformpanel.html")

	t = t.Funcs(template.FuncMap{
		"getComponentId":  getComponentId,
		"compareInt":      CompareInt,
		"compareString":   CompareString,
		"getPropValue":    GetPropValue,
		"dealHTMLEscaper": DealHTMLEscaper,
	})

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/tabformpanel.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})

	data["TabFormPanel"] = tabFormPanel
	data["Terminal"] = terminal

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()
}
