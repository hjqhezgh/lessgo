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

type gridPanel struct {
	Entity       string       `xml:"entity,attr"`
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
}

//link目前可以支持，直接跳转，打开浏览器新窗口跳转，iframe弹窗，询问提示窗
//linkType=currentPage，
//以下为通用配置
//url 必填
//iconUrl 选填 如果有配置iconUrl，则会生成一个可点击的图标
//loadParamName 选填，不填就不带参数
//loadParamValue 如果loadParamName有值，则此配置必填，可取值为id 或者 this
type column struct {
	Field     string `xml:"field,attr"`
	Desc      string `xml:"desc,attr"`
	Hidden    string `xml:"hidden,attr"`
	LoadUrl   string `xml:"loadUrl,attr"`
	Formatter string `xml:"formatter"`
}

type action struct {
	Desc         string `xml:"desc,attr"`
	Url          string `xml:"url,attr"`
	ActionParams string `xml:"actionParams,attr"`
	LinkType     string `xml:"linkType,attr"`
}

type toolaction struct {
	Desc       string `xml:"desc,attr"`
	Url        string `xml:"url,attr"`
	LinkType   string `xml:"linkType,attr"`
	ColorClass string `xml:"colorClass,attr"`
	LoadUrl    string `xml:"loadUrl,attr"`

	//for mutiSelect
	ConfirmMsg string `xml:"confirmMsg,attr"`
	Params     string `xml:"params,attr"`
	Callback   string `xml:"callback"`

	//for addToCheckBox
	CheckboxDesc string `xml:"checkboxDesc,attr"`
}

type checkboxtool struct {
	Desc    string `xml:"desc,attr"`
	LoadUrl string `xml:"loadUrl,attr"`
	SaveUrl string `xml:"saveUrl,attr"`
}

type search struct {
	Field      string `xml:"field,attr"`
	SearchType string `xml:"searchType,attr"`
	InputType  string `xml:"inputType,attr"`
	LocalData  string `xml:"localData,attr"`
	Desc       string `xml:"desc,attr"`
	Url        string `xml:"url,attr"`
	ValueField string `xml:"valueField,attr"`
	DescField  string `xml:"descField,attr"`
	//存储实际的搜索值
	Value  string
	Char14 string `xml:"char14,attr"` //for 时间戳控件
	Char8  string `xml:"char8,attr"`  //for 时间日控件

	ParentSelect string `xml:"parentSelect,attr"` //for remoteSelect
}

func (gridpanel gridPanel) generate(entity Entity, terminal, packageName string) []byte {

	var t *template.Template

	var buf bytes.Buffer

	gridpanel.Id = packageName + "." + gridpanel.Id

	runtimeComponentContain[gridpanel.Id] = gridpanel

	t = template.New("gridpanel.html")

	t = t.Funcs(template.FuncMap{
		"getComponentId": getComponentId,
		"compareInt":     CompareInt,
		"compareString":  CompareString,
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
	data["SearchLength"] = len(gridpanel.Searchs)
	data["ActionLength"] = len(gridpanel.Actions)

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()

}
