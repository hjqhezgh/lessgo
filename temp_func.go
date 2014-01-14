// Title：模板中用的函数集
//
// Description:模板中用的函数集
//
// Author:black
//
// Createtime:2013-08-07 00:47
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-07 00:47 black 创建文档
package lessgo

import (
	"fmt"
	"html/template"
	"math/rand"
	"strings"
	"time"
)

//获取通用model的指定属性的值
func GetPropValue(model *Model, propName string) string {

	if model != nil {
		for _, prop := range model.Props {
			if prop.Name == propName {
				return prop.Value
			}
		}
		Log.Debug("找不到实体", model.Entity.Id, "的属性", propName, "对应的值")
	}

	return ""
}

func SetPropValue(model *Model, propName,newValue string) {

	if model != nil {
		for _, prop := range model.Props {
			if prop.Name == propName {
				prop.Value = newValue
				return
			}
		}
		Log.Debug("找不到实体", model.Entity.Id, "的属性")
	}
}

//获取随机的组件id
func getComponentId(componentType string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	str := ""

	for i := 0; i < 4; i++ {
		str += fmt.Sprint(r.Intn(10))
	}

	return componentType + str
}

//模板中的int类型比较比较
func CompareInt(a, b int, compareType string) (flag bool) {

	switch compareType {

	case "eq":
		if a == b {
			flag = true
		} else {
			flag = false
		}
	case "gt":
		if a > b {
			flag = true
		} else {
			flag = false
		}
	case "ge":
		if a >= b {
			flag = true
		} else {
			flag = false
		}
	case "lt":
		if a < b {
			flag = true
		} else {
			flag = false
		}
	case "le":
		if a <= b {
			flag = true
		} else {
			flag = false
		}
	default:
		if a == b {
			flag = true
		} else {
			flag = false
		}
	}

	return flag
}

//模板中字符串比较
func CompareString(a, b string) bool {
	if a == b {
		return true
	} else {
		return false
	}
}

//替换json字符中的换行等特殊符号
func DealJsonString(str string) string {
	str = strings.Replace(str, "\n", " ", -1)
	str = strings.Replace(str, "\n\r", " ", -1)
	str = strings.Replace(str, "\r\n", " ", -1)
	str = strings.Replace(str, "\r", " ", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	//to fixed
	str = strings.Replace(str, "'", "\\\"", -1)

	return str
}

func DealHTMLEscaper(str string) string {
	return template.HTMLEscaper(str)
}
