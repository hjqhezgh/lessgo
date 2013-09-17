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

type MutiFormPanel struct {
	Id            string        `xml:"id,attr"`
	PublicElement PublicElement `xml:"publicelement"`
	FormTabs      []FormTab     `xml:"formtab"`
}

//公共输入域容器
type PublicElement struct {
	Elements []Element `xml:"element"`
}

type FormTab struct {
	Elements []Element `xml:"element"`
	Entity   string    `xml:"entity,attr"`
	Desc     string    `xml:"desc,attr"`
}

func (mutiFormPanel MutiFormPanel) GenerateMutiFormPanel(terminal, packageName string, r *http.Request) []byte {

	var t *template.Template

	var buf bytes.Buffer

	mutiFormPanel.Id = packageName + "." + mutiFormPanel.Id

	runtimeComponentContain[mutiFormPanel.Id] = mutiFormPanel

	t = template.New("mutiformpanel.html")

	t = t.Funcs(template.FuncMap{
		"getComponentId":  getComponentId,
		"compareInt":      compareInt,
		"compareString":   compareString,
		"getPropValue":    getPropValue,
		"dealHTMLEscaper": dealHTMLEscaper,
	})

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/mutiformpanel.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})

	data["MutiFormPanel"] = mutiFormPanel
	data["Terminal"] = terminal

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()
}
