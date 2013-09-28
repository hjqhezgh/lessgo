// Title：表格及其下属控件
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
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

type formPanel struct {
	Entity   string    `xml:"entity,attr"`
	Load     string    `xml:"load,attr"`
	Id       string    `xml:"id,attr"`
	Elements []element `xml:"element"`
}

type element struct {
	Field        string `xml:"field,attr"`
	Desc         string `xml:"desc,attr"`
	Type         string `xml:"type,attr"`
	LocalData    string `xml:"localData,attr"`
	Url          string `xml:"url,attr"`
	ValueField   string `xml:"valueField,attr"`
	DescField    string `xml:"descField,attr"`
	Readonly     string `xml:"readonly,attr"`
	DefaultValue string `xml:"defaultValue,attr"`
	Validate     string `xml:"validate,attr"`
	ImageEntity  string `xml:"imageEntity,attr"`
	ImagePath    string `xml:"imagePath,attr"`
	Resolutions  string `xml:"resolutions,attr"`
	RefTable     string `xml:"refTable,attr"`
	SelfId       string `xml:"selfId,attr"`
	RefId        string `xml:"refId,attr"`
}

func (formpanel formPanel) generate(entity Entity, terminal, packageName string, r *http.Request) []byte {

	var t *template.Template

	var buf bytes.Buffer

	formpanel.Id = packageName + "." + formpanel.Id

	runtimeComponentContain[formpanel.Id] = formpanel

	t = template.New("formpanel.html")

	t = t.Funcs(template.FuncMap{
		"getComponentId":  getComponentId,
		"compareInt":      CompareInt,
		"compareString":   CompareString,
		"getPropValue":    GetPropValue,
		"dealHTMLEscaper": DealHTMLEscaper,
	})

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/formpanel.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})

	data["Formpanel"] = formpanel
	data["Entity"] = entity
	data["Terminal"] = terminal

	if formpanel.Load == "true" {
		vars := mux.Vars(r)
		id := vars["id"]
		model, err := findById(entity, id)

		if err != nil {
			Log.Error(err.Error())
			return []byte{}
		}

		data["Model"] = model
	}

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()

}
