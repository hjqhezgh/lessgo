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
type Model struct {
	Entity Entity
	Id     int
	Props  []*Prop
}

//通用属性
type Prop struct {
	Name  string
	Value string
}

//entity.xml
type entitys struct {
	XMLName xml.Name `xml:"entitys"`
	Entitys []Entity `xml:"entity"`
}

type Entity struct {
	Id     string  `xml:"id,attr"`
	Pk     string  `xml:"pk"`
	Fields []field `xml:"field"`
	Refs   []ref   `xml:"ref"`
}

//根据id查找出实体
func getEntity(id string) Entity {

	for _, entity := range entityList.Entitys {
		if entity.Id == id {
			return entity
		}
	}

	return Entity{}
}

type ref struct {
	Entity         string  `xml:"entity,attr"`
	Field          string  `xml:"field,attr"`
	RefEntityField string  `xml:"refEntityField,attr"`
	Fields         []field `xml:"field"`
	Alias	   string  `xml:"alias,attr"`
}

type field struct {
	Name string `xml:"name,attr"`
	Desc string `xml:"desc,attr"`
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
