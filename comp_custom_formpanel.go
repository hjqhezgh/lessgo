// Title：自定义表单
//
// Description:
//
// Author:black
//
// Createtime:2013-08-09 16:48
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-09 16:48 black 创建文档
package lessgo

import (
	"bytes"
	"text/template"
)

type customFormPanel struct {
	Load            string       `xml:"load,attr"`
	LoadUrl         string       `xml:"loadUrl,attr"`
	SaveUrl         string       `xml:"saveUrl,attr"`
	Id              string       `xml:"id,attr"`
	PageId          string       `xml:"pageId,attr"`
	Title           string       `xml:"title,attr"`
	Elements        []element    `xml:"element"`
	FormButtons     []formButton `xml:"formButton"`
	BeforeSave      string       `xml:"beforeSave"`
	AfterRender     string       `xml:"afterRender"`
	AfterSave       string       `xml:"afterSave"`
	AfterLoad       string       `xml:"afterLoad"`
	FailCallback    string       `xml:"failCallback"`
	Inwindow        string       `xml:"inwindow,attr"`
	HideSaveButton  string       `xml:"hideSaveButton,attr"`
	HideResetButton string       `xml:"hideResetButton,attr"`
}

func (formpanel customFormPanel) generate(terminal, packageName string, employee Employee) []byte {

	var t *template.Template

	var buf bytes.Buffer

	formpanel.Id = packageName + "." + formpanel.Id

	runtimeComponentContain[formpanel.Id] = formpanel

	t = template.New("customformpanel.html")

	t = t.Funcs(template.FuncMap{
		"getComponentId":  getComponentId,
		"compareInt":      CompareInt,
		"compareString":   CompareString,
		"getPropValue":    GetPropValue,
		"dealHTMLEscaper": DealHTMLEscaper,
	})

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/customformpanel.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})

	data["CustomFormPanel"] = formpanel
	data["Terminal"] = terminal
	data["Employee"] = employee

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()

}
