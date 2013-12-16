// Title：自定义表格
//
// Description:
//
// Author:black
//
// Createtime:2013-08-08 09:29
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-08 09:29 black 创建文档
package lessgo

import (
	"bytes"
	"text/template"
)

type customGridPanel struct {
	Url          string       `xml:"url,attr"`
	PageSize     int          `xml:"pageSize,attr"`
	LoadUrl      string       `xml:"loadUrl,attr"`
	Id           string       `xml:"id,attr"`
	PageId       string       `xml:"pageId,attr"`
	Title        string       `xml:"title,attr"`
	Width        string       `xml:"width,attr"`
	Height       string       `xml:"height,attr"`
	MutiSelect   string       `xml:"mutiSelect,attr"`
	Columns      []column     `xml:"column"`
	Actions      []action     `xml:"action"`
	Searchs      []search     `xml:"search"`
	Checkboxtool checkboxtool `xml:"checkboxtool"`
	ToolActions  []toolaction `xml:"toolaction"`
	BeforeRender string       `xml:"beforeRender"`
	AfterRender  string       `xml:"afterRender"`
	CustomSearch string       `xml:"customSearch"`
	ActionWidth  string       `xml:"actionWidth,attr"`
}

func (gridpanel customGridPanel) generate(terminal, packageName string) []byte {

	var t *template.Template

	var buf bytes.Buffer

	gridpanel.Id = packageName + "." + gridpanel.Id

	runtimeComponentContain[gridpanel.Id] = gridpanel

	t = template.New("customgridpanel.html")

	t = t.Funcs(template.FuncMap{
		"getComponentId": getComponentId,
		"compareInt":     CompareInt,
		"compareString":  CompareString,
	})

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/customgridpanel.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})

	data["CustomGridPanel"] = gridpanel
	data["Terminal"] = terminal
	data["SearchLength"] = len(gridpanel.Searchs)
	data["ActionLength"] = len(gridpanel.Actions)

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()

}
