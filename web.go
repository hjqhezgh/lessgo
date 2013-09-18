// Title：web服务相关
//
// Description:
//
// Author:black
//
// Createtime:2013-08-06 15:43
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-06 15:43 black 创建文档
package lessgo

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"github.com/hjqhezgh/commonlib"
)

//跳转至错误页面
func errMessage(w http.ResponseWriter, r *http.Request, errMsg string) {

	m := make(map[string]interface{})

	m["ErrMsg"] = errMsg

	m["Nav"] = navs

	commonlib.RenderTemplate(w, r, "err_message.html", m, nil, "../template/err_message.html", "../template/nav.html")
}

//中心控制器
func homeAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	commonlib.RenderTemplate(w, r, "home.html", m, nil, "../template/home.html")
}

//中心控制器
func commonAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	entity, opera, terminal, msg := analyseUrl(r.URL.Path)

	if msg != "" {
		Log.Warn(msg, r.URL.Path)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+msg)
		return
	} else {

		err := analyNav(terminal)

		if err != nil {
			errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
			return
		}

		//将导航数据放入页面上
		m["Nav"] = navs

		switch opera {
		case "home":
			Log.Debug("路径：", r.URL.Path, "访问应用首页")

			content, err := ioutil.ReadFile("../etc/view/" + terminal + "/home.xml")

			if err != nil {
				Log.Error(err)
				errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
				return
			}

			packageName := terminal + "." + "home"

			w.Write(generate(content, terminal, packageName, r))
		case "index":
			dealEntityIndex(entity, terminal, m, w, r)
		case "add":
			dealEntityAdd(entity, terminal, m, w, r)
		case "modify":
			dealEntityModify(entity, terminal, m, w, r)
		case "save":
			dealEntitySave(entity, w, r)
		case "delete":
			dealEntityDelete(entity, w, r)
		case "detail":
			dealEntityDetail(entity, m, w, r)
		case "page":
			dealEntityPage(entity, w, r)
		case "alldata":
			dealEntityAllData(entity, w, r)
		case "load":
			dealEntityLoad(entity, w, r)
		default:
			Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的未知页")
			commonlib.RenderTemplate(w, r, "home.html", m, nil, "../template/home.html", "../template/nav.html")
		}
	}
}

//中心控制器
func independentAction(w http.ResponseWriter, r *http.Request) {

	Log.Debug("访问自定义路径：", r.URL.Path)

	strs := strings.Split(r.URL.Path, "/")

	terminal := strs[1]

	err := analyNav(terminal)

	if err != nil {
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	view := ""

	for _, url := range urls.Urls {

		if url.Path == r.URL.Path {
			view = url.View
		}

	}

	if view == "" {
		Log.Warn("路径：", r.URL.Path, "找不到页面")
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+"找不到页面")
		return
	}

	content, err := ioutil.ReadFile("../etc/view/" + view)

	if err != nil {
		Log.Error(err)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	packageName := terminal + "." + r.URL.Path

	w.Write(generate(content, terminal, packageName, r))
}

//解析导航
func analyNav(terminal string) error {

	navs = Navs{}

	content, err := ioutil.ReadFile("../etc/view/" + terminal + "/nav.xml")

	if err != nil {
		Log.Error(err)
		return err
	}

	err = xml.Unmarshal(content, &navs)

	if err != nil {
		Log.Error(err)
		return err
	}

	return nil
}

//分析URL得出当前url访问的实体模块，以及进行的操作，如果有错误，就去读取msg
func analyseUrl(url string) (entity Entity, operation, terminal, msg string) {

	strs := strings.Split(url, "/")

	//首页的情况
	if len(strs) == 2 || strs[2] == "index.html" {
		return Entity{}, "home", strs[1], ""
	} else {

		entity = GetEntity(strs[2])

		if entity.Id == "" {
			return Entity{}, "", "", "找不到该url下的相应实体"
		}

		if len(strs) == 3 {
			return entity, "index", strs[1], ""
		}

		switch strs[3] {
		case "index.html":
			return entity, "index", strs[1], ""
		case "add":
			return entity, "add", strs[1], ""
		case "modify":
			return entity, "modify", strs[1], ""
		case "save":
			return entity, "save", strs[1], ""
		case "delete":
			return entity, "delete", strs[1], ""
		case "page":
			return entity, "page", strs[1], ""
		case "alldata":
			return entity, "alldata", strs[1], ""
		case "load":
			return entity, "load", strs[1], ""
		default:
			_, err := strconv.Atoi(strs[4])

			if err != nil {
				return Entity{}, "", "", "找不到该url下对应的操作"
			} else {
				return entity, "detail", strs[1], ""
			}
		}
	}

	return
}

//处理实体的列表页请求
func dealEntityIndex(entity Entity, terminal string, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的列表页")

	content, err := ioutil.ReadFile("../etc/view/" + terminal + "/" + entity.Id + "/index.xml")

	if err != nil {
		Log.Error(err)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	packageName := terminal + "." + entity.Id + ".index"

	w.Write(generate(content, terminal, packageName, r))
}

//处理实体的添加页请求
func dealEntityAdd(entity Entity, terminal string, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的添加页")

	content, err := ioutil.ReadFile("../etc/view/" + terminal + "/" + entity.Id + "/add.xml")

	if err != nil {
		Log.Error(err)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	packageName := terminal + "." + entity.Id + ".add"

	w.Write(generate(content, terminal, packageName, r))
}

//处理实体的修改页请求
func dealEntityModify(entity Entity, terminal string, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的修改页")

	content, err := ioutil.ReadFile("../etc/view/" + terminal + "/" + entity.Id + "/modify.xml")

	if err != nil {
		Log.Error(err)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	packageName := terminal + "." + entity.Id + ".modify"

	w.Write(generate(content, terminal, packageName, r))
}

//处理实体的保存页请求
func dealEntitySave(entity Entity, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的保存ajax请求")

	err := r.ParseForm()

	//异步请求绑定的组件Id
	componentId := r.FormValue("componentId")
	formpanel := runtimeComponentContain[componentId].(formPanel)

	m := make(map[string]interface{})

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		Log.Warn(err.Error())
		commonlib.OutputJson(w, m," ")
		return
	}

	idString := r.FormValue(entity.Pk)

	model := new(Model)
	model.Props = []*Prop{}

	for _, formElement := range formpanel.Elements {
		prop := new(Prop)

		if formElement.Type == "image" { //图片类型需要做多表处理

			filePath := r.FormValue(formElement.Field)
			tmpFileName := ""

			if filePath != "" {
				imgEntity := GetEntity(formElement.ImageEntity)
				imageModel := new(Model)
				imageModel.Props = []*Prop{}

				fileNameProp := new(Prop)
				fileNameProp.Name = "filename"
				fileNameProp.Value = commonlib.SubstrByStEd(filePath, strings.LastIndex(filePath, "/")+1, strings.LastIndex(filePath, "."))
				tmpFileName = commonlib.Substr(filePath, strings.LastIndex(filePath, "/")+1, len(filePath))

				suffixProp := new(Prop)
				suffixProp.Name = "suffix"
				suffixProp.Value = commonlib.Substr(filePath, strings.LastIndex(filePath, ".")+1, len(filePath))

				imageModel.Props = append(imageModel.Props, fileNameProp)
				imageModel.Props = append(imageModel.Props, suffixProp)

				imgId, err := insert(imgEntity, imageModel, []element{
					element{
						Field: "filename",
					},
					element{
						Field: "suffix",
					},
				})

				if err != nil {
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					Log.Warn(err.Error())
					commonlib.OutputJson(w, m," ")
					return
				}

				prop.Name = formElement.Field
				prop.Value = fmt.Sprint(imgId)
				model.Props = append(model.Props, prop)

				tmpFile, err := os.OpenFile("../tmp/"+tmpFileName, os.O_RDWR, 0777)

				if err != nil {
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					Log.Warn(err.Error())
					commonlib.OutputJson(w, m," ")
					return
				}

				_, err = os.Stat(formElement.ImagePath)

				if err != nil && os.IsNotExist(err) {
					Log.Info(formElement.ImagePath, "文件夹不存在，创建")
					os.Mkdir(formElement.ImagePath, 0777)
				}

				disFile, err := os.Create(formElement.ImagePath + "/" + tmpFileName)

				if err != nil {
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					Log.Warn(err.Error())
					commonlib.OutputJson(w, m," ")
					return
				}

				io.Copy(disFile, tmpFile)

				os.Remove("../tmp/" + tmpFileName)
			} else {
				prop.Name = formElement.Field
				prop.Value = fmt.Sprint(0)
				model.Props = append(model.Props, prop)
			}

		} else {
			prop.Name = formElement.Field
			prop.Value = r.FormValue(formElement.Field)
			model.Props = append(model.Props, prop)
		}
	}

	if idString != "" { //修改的情况
		id, err := strconv.Atoi(idString)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			Log.Warn(err.Error())
			commonlib.OutputJson(w, m," ")
			return
		}

		model.Id = id

		err = modify(entity, model, formpanel.Elements)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m," ")
			return
		}

		m["success"] = true
		m["code"] = 200
		commonlib.OutputJson(w, m," ")
		return
	} else {
		_, err = insert(entity, model, formpanel.Elements)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m," ")
			return
		}

		m["success"] = true
		m["code"] = 200
		commonlib.OutputJson(w, m," ")
		return
	}

}

//处理实体的详细页请求
func dealEntityDetail(entity Entity, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的详细信息页")

	vars := mux.Vars(r)
	id := vars["id"] //先假设这个是活动的ID

	model, err := findById(entity, id)

	if err != nil {
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	m["Entity"] = entity
	m["Model"] = model

	commonlib.RenderTemplate(w, r, "entity_detail.html", m, template.FuncMap{"getPropValue": getPropValue}, "../template/entity_detail.html", "../template/nav.html")

}

//处理实体的删除页请求
func dealEntityDelete(entity Entity, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的删除页")

	m := make(map[string]interface{})

	vars := mux.Vars(r)
	id := vars["id"]

	err := delete(entity, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m," ")
		return
	}

	m["success"] = true
	m["code"] = 200
	commonlib.OutputJson(w, m," ")
	return

}

//处理实体的分页ajax请求
func dealEntityPage(entity Entity, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的分页数据ajax请求")

	err := r.ParseForm()

	m := make(map[string]interface{})

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m," ")
		return
	}

	pageNoString := r.FormValue("page")
	pageNo := 1
	if pageNoString != "" {
		pageNo, err = strconv.Atoi(pageNoString)
		if err != nil {
			pageNo = 1
			Log.Warn("错误的pageNo:", pageNo)
		}
	}

	pageSizeString := r.FormValue("rows")
	pageSize := 10
	if pageSizeString != "" {
		pageSize, err = strconv.Atoi(pageSizeString)
		if err != nil {
			Log.Warn("错误的pageSize:", pageSize)
		}
	}

	//异步请求绑定的组件Id
	componentId := r.FormValue("componentId")
	gridpanel := runtimeComponentContain[componentId].(gridPanel)

	searchParam := []search{}

	for key, value := range r.Form {
		//滤除分页组件自带的参数，其他参数都认定为搜索参数， Fixme 后期改进这块机制
		if key != "page" && key != "rows" && key != "componentId" && key != "_search" && key != "nd" && key != "sidx" && key != "sord" && key != "filters" {
			strs := strings.Split(key, "-")
			search := search{
				Field:      strs[0],
				SearchType: strs[1],
			}

			if len(value) > 0 {
				search.Value = value[0]
			} else {
				search.Value = ""
			}

			searchParam = append(searchParam, search)
		}
	}

	pageData, err := findTraditionPage(entity, pageNo, pageSize, searchParam, gridpanel.Columns)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m," ")
		return
	}

	m["PageData"] = pageData
	m["Entity"] = entity
	m["Gridpanel"] = gridpanel
	m["DataLength"] = len(pageData.Datas) - 1
	if len(pageData.Datas) > 0 {
		m["FieldLength"] = len(pageData.Datas[0].(*Model).Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_page.json", m, template.FuncMap{"getPropValue": getPropValue, "compareInt": compareInt, "dealJsonString": dealJsonString}, "../template/entity_page.json")
}

//处理实体的所有数据ajax请求
func dealEntityAllData(entity Entity, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的所有数据ajax请求")

	err := r.ParseForm()

	m := make(map[string]interface{})

	models, err := findAllData(entity)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m," ")
		return
	}

	m["Models"] = models
	m["Entity"] = entity
	m["DataLength"] = len(models) - 1
	if len(models) > 0 {
		m["FieldLength"] = len(models[0].Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_alldata.json", m, template.FuncMap{"getPropValue": getPropValue, "compareInt": compareInt}, "../template/entity_alldata.json")
}

//处理实体的分页ajax请求
func dealEntityLoad(entity Entity, w http.ResponseWriter, r *http.Request) {

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的load单实体ajax请求")

	m := make(map[string]interface{})

	vars := mux.Vars(r)
	id := vars["id"]

	model, err := findById(entity, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m," ")
		return
	}

	if model.Id == 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:找不到相应的实体"
		commonlib.OutputJson(w, m," ")
		return
	}

	m["Entity"] = entity
	m["Model"] = model
	m["FieldLength"] = len(entity.Fields) - 1

	commonlib.RenderTemplate(w, r, "entity_load.json", m, template.FuncMap{"getPropValue": getPropValue, "compareInt": compareInt}, "../template/entity_load.json")
}

//多实体保存ajax请求处理器
func mutiSavaAction(w http.ResponseWriter, r *http.Request) {

	Log.Debug("访问多表保存ajax路径：", r.URL.Path)

	err := r.ParseForm()

	m := make(map[string]interface{})
	modelMap := make(map[string]*Model)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m," ")
		return
	}

	//异步请求绑定的组件Id
	componentId := r.FormValue("componentId")
	mutiFormPanel := runtimeComponentContain[componentId].(mutiFormPanel)

	//公有属性赋值
	for _, tab := range mutiFormPanel.FormTabs {
		model := new(Model)
		modelMap[tab.Entity] = model
		model.Props = []*Prop{}

		for _, element := range mutiFormPanel.PublicElement.Elements {
			prop := new(Prop)
			prop.Name = element.Field
			prop.Value = r.FormValue(element.Field)
			model.Props = append(model.Props, prop)
		}

	}

	//对各实体自己的属性赋值
	for key, _ := range r.Form {
		if strings.Contains(key, ".") {
			strs := strings.Split(key, ".")

			model := modelMap[strs[0]]
			prop := new(Prop)
			prop.Name = strs[1]
			prop.Value = r.FormValue(key)
			model.Props = append(model.Props, prop)
		}
	}

	//数据插入
	//对各实体自己的属性赋值
	for _, tab := range mutiFormPanel.FormTabs {

		elements := tab.Elements

		for _, element := range mutiFormPanel.PublicElement.Elements {
			elements = append(elements, element)
		}

		_, err = insert(GetEntity(tab.Entity), modelMap[tab.Entity], elements)
		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m," ")
			return
		}
	}

	m["success"] = true
	m["code"] = 200
	commonlib.OutputJson(w, m," ")
	return
}
