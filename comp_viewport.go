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

//当viewport下面可以放置其他元素的时候，就扩展Viewport结构体，同时记得扩展容器结构体
type Viewport struct {
	XMLName        xml.Name        `xml:"viewport"`
	GridPanels     []GridPanel     `xml:"gridpanel"`
	FormPanels     []FormPanel     `xml:"formpanel"`
	MutiFormPanels []MutiFormPanel `xml:"mutiformpanel"`
}

//扩展viewport的同时，记得同时扩展container
func (viewport Viewport) GenerateViewport(terminal, packageName string, r *http.Request) []byte {

	content := ""

	for _, gridpanel := range viewport.GridPanels {
		content += string(gridpanel.GenerateGridPanel(GetEntity(gridpanel.Entity), terminal, packageName))
	}

	for _, formpanel := range viewport.FormPanels {
		content += string(formpanel.GenerateFormPanel(GetEntity(formpanel.Entity), terminal, packageName, r))
	}

	for _, mutiformpanel := range viewport.MutiFormPanels {
		content += string(mutiformpanel.GenerateMutiFormPanel(terminal, packageName, r))
	}

	var t *template.Template

	var buf bytes.Buffer

	t = template.New("viewport.html")

	t, err := t.ParseFiles("../lessgo/template/component/"+terminal+"/viewport.html", "../lessgo/template/component/"+terminal+"/nav.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})
	data["Content"] = content
	data["Nav"] = navs

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()
}
