// Title：空白页面方便用户可以自己填词内容
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

type blankPanel struct {
	Content string `xml:"content"`
}

func (blankPanel blankPanel) generate(terminal, packageName string, employee Employee) []byte {

	var t *template.Template

	var buf bytes.Buffer

	t = template.New("blankpanel.html")

	t, err := t.ParseFiles("../lessgo/template/component/" + terminal + "/blankpanel.html")

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	data := make(map[string]interface{})

	data["BlankPanel"] = blankPanel
	data["Terminal"] = terminal
	data["Employee"] = employee

	err = t.Execute(&buf, data)

	if err != nil {
		Log.Error(err.Error())
		return []byte{}
	}

	return buf.Bytes()

}
