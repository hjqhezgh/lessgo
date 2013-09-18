// Title：配置文件相关的模型
//
// Description:
//
// Author:black
//
// Createtime:2013-08-06 14:22
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-06 14:22 black 创建文档
package lessgo

import (
	"encoding/xml"
)

//通用模型，用于存储从数据库获取到的值
type model struct {
	Entity entity
	Id     int
	Props  []*prop
}

//通用属性
type prop struct {
	Name  string
	Value string
}

//entity.xml
type entitys struct {
	XMLName xml.Name `xml:"entitys"`
	Entitys []entity `xml:"entity"`
}

type entity struct {
	Id     string  `xml:"id,attr"`
	Pk     string  `xml:"pk"`
	Fields []field `xml:"field"`
	Refs   []ref   `xml:"ref"`
}

//根据id查找出实体
func getEntity(id string) entity {

	for _, entity := range entityList.Entitys {
		if entity.Id == id {
			return entity
		}
	}

	return entity{}
}

type ref struct {
	Entity         string  `xml:"entity,attr"`
	Field          string  `xml:"field,attr"`
	RefEntityField string  `xml:"refEntityField,attr"`
	Fields         []field `xml:"field"`
}

type field struct {
	Name string `xml:"name,attr"`
	Desc string `xml:"desc,attr"`
}

//nav.xml
type navs struct {
	XMLName xml.Name `xml:"navs"`
	Navs    []nav    `xml:"nav"`
}

type nav struct {
	Text  string `xml:"text,attr"`
	Items []item `xml:"item"`
}

type item struct {
	Path    string `xml:"path,attr"`
	Text    string `xml:"text,attr"`
	Outside bool   `xml:"outside,attr"`
}

//url.xml
type urls struct {
	XMLName   xml.Name `xml:"urls"`
	Urls      []url    `xml:"url"`
	Terminals []string `xml:"terminal"`
}

type url struct {
	Path string `xml:"path,attr"`
	View string `xml:"view,attr"`
}
