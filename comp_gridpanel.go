// Title：表格，及其下属控件
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

type GridPanel struct {
	Entity   string   `xml:"entity,attr"`
	PageSize int      `xml:"pageSize,attr"`
	Id       string   `xml:"id,attr"`
	Title    string   `xml:"title,attr"`
	Width    string   `xml:"width,attr"`
	Height   string   `xml:"height,attr"`
	Columns  []Column `xml:"column"`
	Actions  []Action `xml:"action"`
	Searchs  []Search `xml:"search"`
}

type Column struct {
	Field string `xml:"field,attr"`
	Desc  string `xml:"desc,attr"`
}

type Action struct {
	Desc   string `xml:"desc,attr"`
	Action string `xml:"action,attr"`
	Url    string `xml:"url,attr"`
}

type Search struct {
	Field      string `xml:"field,attr"`
	SearchType string `xml:"searchType,attr"`
	InputType  string `xml:"inputType,attr"`
	LocalData  string `xml:"localData,attr"`
	Desc       string `xml:"desc,attr"`
	Url        string `xml:"url,attr"`
	ValueField string `xml:"valueField,attr"`
	DescField  string `xml:"descField,attr"`
	//存储实际的搜索值
	Value string
}

func (gridpanel GridPanel) GenerateGridPanel(entity Entity, terminal, packageName string) []byte {

	var t *template.Template

	var buf bytes.Buffer

	gridpanel.Id = packageName + "." + gridpanel.Id

	runtimeComponentContain[gridpanel.Id] = gridpanel

	t = template.New("gridpanel.html")

	t = t.Funcs(template.FuncMap{
		"getComponentId": getComponentId,
		"compareInt":     compareInt,
	})

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/gridpanel.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})

	data["Gridpanel"] = gridpanel
	data["Entity"] = entity
	data["Terminal"] = terminal

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()

}
