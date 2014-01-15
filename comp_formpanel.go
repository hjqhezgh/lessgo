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
	"net/http"
	"text/template"
)

type formPanel struct {
	Entity          string       `xml:"entity,attr"`
	Load            string       `xml:"load,attr"`
	Id              string       `xml:"id,attr"`
	PageId          string       `xml:"pageId,attr"`
	Title           string       `xml:"title,attr"`
	Elements        []element    `xml:"element"`
	FormButtons     []formButton `xml:"formButton"`
	BeforeSave      string       `xml:"beforeSave"`
	AfterRender     string       `xml:"afterRender"`
	AfterSave       string       `xml:"afterSave"`
	Inwindow        string       `xml:"inwindow,attr"`
	HideSaveButton  string       `xml:"hideSaveButton,attr"`
	HideResetButton string       `xml:"hideResetButton,attr"`
}

type formButton struct {
	Desc        string `xml:"desc,attr"`
	ButtonClass string `xml:"buttonClass,attr"`
	Handler     string `xml:"handler"`
}

type element struct {
	Field    string `xml:"field,attr"`
	Desc     string `xml:"desc,attr"`
	Type     string `xml:"type,attr"`
	Readonly string `xml:"readonly,attr"`
	Validate string `xml:"validate,attr"`
	Tip      string `xml:"tip,attr"`

	LocalData string `xml:"localData,attr"` //for 本地下拉框

	Url        string `xml:"url,attr"`        //for 远程下拉框
	ValueField string `xml:"valueField,attr"` //for 远程下拉框，多选框控件
	DescField  string `xml:"descField,attr"`  //for 远程下拉框，多选框控件

	DefaultValue string `xml:"defaultValue,attr"` //for 隐藏域、本地下拉框

	RefTable string `xml:"refTable,attr"` //for 多选框控件
	SelfId   string `xml:"selfId,attr"`   //for 多选框控件
	RefId    string `xml:"refId,attr"`    //for 多选框控件

	ImageEntity string `xml:"imageEntity,attr"` //for Image控件
	ImagePath   string `xml:"imagePath,attr"`   //for Image控件
	Widths      string `xml:"widths,attr"`      //for Image控件
	MaxWidth    string `xml:"maxWidth,attr"`    //for Image控件
	MaxHeight   string `xml:"maxHeight,attr"`   //for Image控件
	MinWidth    string `xml:"minWidth,attr"`    //for Image控件
	MinHeight   string `xml:"minHeight,attr"`   //for Image控件
	MaxSize     string `xml:"maxSize,attr"`     //for Image控件
	Resolution  string `xml:"resolution,attr"`  //for Image控件
	ImageType   string `xml:"imageType,attr"`   //for Image控件

	UploadUrl string `xml:"uploadUrl,attr"` //for HTML编辑器控件

	Char14 string `xml:"char14,attr"` //for 时间戳控件
	Char8  string `xml:"char8,attr"`  //for 时间日控件

	ParentSelect string `xml:"parentSelect,attr"` //for remoteSelect
	Params       string `xml:"params,attr"`       //for remoteSelect

	FileUploadUrl string `xml:"fileUploadUrl,attr"` //for fileupload
}

func (formpanel formPanel) generate(entity Entity, terminal, packageName string, r *http.Request, employee Employee) []byte {

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
	data["Employee"] = employee

	if formpanel.Load == "true" {
		id := r.FormValue("id")
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
