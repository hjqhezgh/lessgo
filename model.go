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

//config.xml
type Config struct {
	XMLName     xml.Name `xml:"config"`
	Port        int      `xml:"port"`
	DbUrl       string   `xml:"dbUrl"`
	DbName      string   `xml:"dbName"`
	DbUserName  string   `xml:"dbUserName"`
	DbPwd       string   `xml:"dbPwd"`
	MaxPoolSize int      `xml:"maxPoolSize"`
}

//entity.xml
type Entitys struct {
	XMLName xml.Name `xml:"entitys"`
	Entitys []Entity `xml:"entity"`
}

type Entity struct {
	Id     string  `xml:"id,attr"`
	Pk     string  `xml:"pk"`
	Fields []Field `xml:"field"`
	Refs   []Ref   `xml:"ref"`
}

//根据id查找出实体
func GetEntity(id string) Entity {

	for _, entity := range entitys.Entitys {
		if entity.Id == id {
			return entity
		}
	}

	return Entity{}
}

type Ref struct {
	Entity         string  `xml:"entity,attr"`
	Field          string  `xml:"field,attr"`
	RefEntityField string  `xml:"refEntityField,attr"`
	Fields         []Field `xml:"field"`
}

type Field struct {
	Name string `xml:"name,attr"`
	Desc string `xml:"desc,attr"`
}

//nav.xml
type Navs struct {
	XMLName xml.Name `xml:"navs"`
	Navs    []Nav    `xml:"nav"`
}

type Nav struct {
	Text  string `xml:"text,attr"`
	Items []Item `xml:"item"`
}

type Item struct {
	Path    string `xml:"path,attr"`
	Text    string `xml:"text,attr"`
	Outside bool   `xml:"outside,attr"`
}

//url.xml
type Urls struct {
	XMLName   xml.Name `xml:"urls"`
	Urls      []Url    `xml:"url"`
	Terminals []string `xml:"terminal"`
}

type Url struct {
	Path string `xml:"path,attr"`
	View string `xml:"view,attr"`
}
