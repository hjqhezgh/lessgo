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
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hjqhezgh/commonlib"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//跳转至错误页面
func errMessage(w http.ResponseWriter, r *http.Request, errMsg string) {

	m := make(map[string]interface{})

	m["ErrMsg"] = errMsg
	m["SiteName"] = SiteName
	m["SiteIcon"] = SiteIcon

	commonlib.RenderTemplate(w, r, "err_message.html", m, nil, "../lessgo/template/err_message.html")
}

//注销
func loginOutAction(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, SESSION_USER)
	for v, _ := range session.Values {
		delete(session.Values, v)
	}
	session.Save(r, w)
	w.Write([]byte("success"))
}

//中心控制器
func homeAction(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]interface{})
	m["SiteName"] = SiteName
	m["SiteIcon"] = SiteIcon
	commonlib.RenderTemplate(w, r, "home.html", m, nil, "../lessgo/template/home.html")
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

		switch opera {
		case "home":
			employee := GetCurrentEmployee(r)

			if employee.UserId == "" {
				Log.Warn("用户未登陆")
				m["SiteName"] = SiteName
				m["SiteIcon"] = SiteIcon
				commonlib.RenderTemplate(w, r, "login.html", m, nil, "../lessgo/template/component/"+terminal+"/login.html")
				return
			}

			Log.Debug("路径：", r.URL.Path, "访问应用首页")

			content, err := ioutil.ReadFile("../etc/view/" + terminal + "/home.xml")

			if err != nil {
				Log.Error(err)
				errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
				return
			}

			packageName := terminal + "." + "home"

			w.Write(generate(content, terminal, packageName, r, employee))
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
			m["SiteName"] = SiteName
			m["SiteIcon"] = SiteIcon
			commonlib.RenderTemplate(w, r, "home.html", m, nil, "../lessgo/template/home.html")
		}
	}
}

//中心控制器
func independentAction(w http.ResponseWriter, r *http.Request) {

	strs := strings.Split(r.URL.Path, "/")

	terminal := strs[1]

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m := make(map[string]interface{})
		m["SiteName"] = SiteName
		m["SiteIcon"] = SiteIcon
		commonlib.RenderTemplate(w, r, "login.html", nil, nil, "../lessgo/template/component/"+terminal+"/login.html")
		return
	}

	Log.Debug("访问自定义路径：", r.URL.Path)

	view := ""

	for _, url := range urlList.Urls {

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

	packageName := terminal + "." + view

	w.Write(generate(content, terminal, packageName, r, employee))
}

//中心控制器
func loginAction(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	m["SiteName"] = SiteName
	m["SiteIcon"] = SiteIcon

	strs := strings.Split(r.URL.Path, "/")

	terminal := strs[1]

	Log.Debug("访问登陆页")

	commonlib.RenderTemplate(w, r, "login.html", m, nil, "../lessgo/template/component/"+terminal+"/login.html")
}

//分析URL得出当前url访问的实体模块，以及进行的操作，如果有错误，就去读取msg
func analyseUrl(url string) (_entity Entity, operation, terminal, msg string) {

	strs := strings.Split(url, "/")

	//首页的情况
	if len(strs) == 2 || strs[2] == "index.html" {
		return Entity{}, "home", strs[1], ""
	} else {

		_entity = getEntity(strs[2])

		if _entity.Id == "" {
			return Entity{}, "", "", "找不到该url下的相应实体"
		}

		if len(strs) == 3 {
			return _entity, "index", strs[1], ""
		}

		switch strs[3] {
		case "index.html":
			return _entity, "index", strs[1], ""
		case "add":
			return _entity, "add", strs[1], ""
		case "modify":
			return _entity, "modify", strs[1], ""
		case "save":
			return _entity, "save", strs[1], ""
		case "delete":
			return _entity, "delete", strs[1], ""
		case "page":
			return _entity, "page", strs[1], ""
		case "alldata":
			return _entity, "alldata", strs[1], ""
		case "load":
			return _entity, "load", strs[1], ""
		default:
			_, err := strconv.Atoi(strs[4])

			if err != nil {
				return Entity{}, "", "", "找不到该url下对应的操作"
			} else {
				return _entity, "detail", strs[1], ""
			}
		}
	}

	return
}

//处理实体的列表页请求
func dealEntityIndex(entity Entity, terminal string, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["SiteName"] = SiteName
		m["SiteIcon"] = SiteIcon
		commonlib.RenderTemplate(w, r, "login.html", m, nil, "../lessgo/template/component/"+terminal+"/login.html")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的列表页")

	content, err := ioutil.ReadFile("../etc/view/" + terminal + "/" + entity.Id + "/index.xml")

	if err != nil {
		Log.Error(err)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	packageName := terminal + "." + entity.Id + ".index"

	w.Write(generate(content, terminal, packageName, r, employee))
}

//处理实体的添加页请求
func dealEntityAdd(entity Entity, terminal string, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["SiteName"] = SiteName
		m["SiteIcon"] = SiteIcon
		commonlib.RenderTemplate(w, r, "login.html", m, nil, "../lessgo/template/component/"+terminal+"/login.html")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的添加页")

	content, err := ioutil.ReadFile("../etc/view/" + terminal + "/" + entity.Id + "/add.xml")

	if err != nil {
		Log.Error(err)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	packageName := terminal + "." + entity.Id + ".add"

	w.Write(generate(content, terminal, packageName, r, employee))
}

//处理实体的修改页请求
func dealEntityModify(entity Entity, terminal string, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["SiteName"] = SiteName
		m["SiteIcon"] = SiteIcon
		commonlib.RenderTemplate(w, r, "login.html", m, nil, "../lessgo/template/component/"+terminal+"/login.html")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的修改页")

	content, err := ioutil.ReadFile("../etc/view/" + terminal + "/" + entity.Id + "/modify.xml")

	if err != nil {
		Log.Error(err)
		errMessage(w, r, "出现错误，请联系IT部门，错误信息:"+err.Error())
		return
	}

	packageName := terminal + "." + entity.Id + ".modify"

	w.Write(generate(content, terminal, packageName, r, employee))
}

//处理实体的保存页请求
func dealEntitySave(_entity Entity, w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", _entity.Id, "的保存ajax请求")

	err := r.ParseForm()

	//异步请求绑定的组件Id
	componentId := r.FormValue("componentId")
	formpanel := runtimeComponentContain[componentId].(formPanel)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		Log.Warn(err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}

	idString := r.FormValue(_entity.Pk)

	_model := new(Model)
	_model.Props = []*Prop{}

	imgElements := []element{}

	for _, formElement := range formpanel.Elements {
		_prop := new(Prop)

		if formElement.Type == "image" { //图片类型需要做多表处理
			imgElements = append(imgElements, formElement)
		} else if formElement.Type == "currentTime" { //当前时间，一般用于createTime
			_prop.Name = formElement.Field
			if formElement.Char14 == "true" {
				_prop.Value = time.Now().Format("20060102150405")
			} else {
				_prop.Value = time.Now().Format("2006-01-02 15:04:05")
			}
			_model.Props = append(_model.Props, _prop)
		} else {
			_prop.Name = formElement.Field
			_prop.Value = r.FormValue(formElement.Field)
			_model.Props = append(_model.Props, _prop)
		}
	}

	if idString != "" { //修改的情况
		id, err := strconv.Atoi(idString)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			Log.Warn(err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		_model.Id = id

		err = modify(_entity, _model, formpanel.Elements)

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		m["success"] = true
		m["code"] = 200
		commonlib.OutputJson(w, m, " ")
		return
	} else {
		id, err := insert(_entity, _model, formpanel.Elements)

		for _, imgElment := range imgElements {

			filePath := r.FormValue(imgElment.Field)
			tmpFileName := ""

			if filePath != "" {
				imgEntity := getEntity(imgElment.ImageEntity)
				imageModel := new(Model)
				imageModel.Props = []*Prop{}

				fileNameProp := new(Prop)
				fileNameProp.Name = "filename"
				fileNameProp.Value = commonlib.SubstrByStEd(filePath, strings.LastIndex(filePath, "/")+1, strings.LastIndex(filePath, "."))
				tmpFileName = commonlib.Substr(filePath, strings.LastIndex(filePath, "/")+1, len(filePath))

				suffixProp := new(Prop)
				suffixProp.Name = "suffix"
				suffixProp.Value = commonlib.Substr(filePath, strings.LastIndex(filePath, ".")+1, len(filePath))

				refProp := new(Prop)
				refProp.Name = imgElment.Field
				refProp.Value = fmt.Sprint(id)

				imageModel.Props = append(imageModel.Props, fileNameProp)
				imageModel.Props = append(imageModel.Props, suffixProp)
				imageModel.Props = append(imageModel.Props, refProp)

				_, err := insert(imgEntity, imageModel, []element{
					element{
						Field: "filename",
					},
					element{
						Field: "suffix",
					},
					element{
						Field: imgElment.Field,
					},
				})

				if err != nil {
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					Log.Warn(err.Error())
					commonlib.OutputJson(w, m, " ")
					return
				}

				tmpFile, err := os.OpenFile("../tmp/"+tmpFileName, os.O_RDWR, 0777)

				if err != nil {
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					Log.Warn(err.Error())
					commonlib.OutputJson(w, m, " ")
					return
				}

				_, err = os.Stat(imgElment.ImagePath)

				if err != nil && os.IsNotExist(err) {
					Log.Info(imgElment.ImagePath, "文件夹不存在，创建")
					os.MkdirAll(imgElment.ImagePath, 0777)
				}

				disFile, err := os.Create(imgElment.ImagePath + "/" + tmpFileName)

				if err != nil {
					m["success"] = false
					m["code"] = 100
					m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
					Log.Warn(err.Error())
					commonlib.OutputJson(w, m, " ")
					return
				}

				io.Copy(disFile, tmpFile)

				os.Remove("../tmp/" + tmpFileName)
			}
		}

		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
			commonlib.OutputJson(w, m, " ")
			return
		}

		m["success"] = true
		m["code"] = 200
		commonlib.OutputJson(w, m, " ")
		return
	}

}

//处理实体的详细页请求
func dealEntityDetail(entity Entity, m map[string]interface{}, w http.ResponseWriter, r *http.Request) {

	//todo 目前还没实现此页面
	/*
		employee := GetCurrentEmployee(r)

		if employee.UserId == "" {
			Log.Warn("用户未登陆")
			commonlib.RenderTemplate(w, r, "login.html", m, nil, "../lessgo/template/login.html")
			return
		}

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

		commonlib.RenderTemplate(w, r, "entity_detail.html", m, template.FuncMap{"getPropValue": getPropValue}, "../lessgo/template/entity_detail.html")*/

}

//处理实体的删除页请求
func dealEntityDelete(entity Entity, w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的删除页")

	id := r.FormValue("id")

	err := deleteEntity(entity, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["success"] = true
	m["code"] = 200
	commonlib.OutputJson(w, m, " ")
	return

}

//处理实体的分页ajax请求
func dealEntityPage(entity Entity, w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的分页数据ajax请求")

	err := r.ParseForm()

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
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
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["PageData"] = pageData
	m["Entity"] = entity
	m["Gridpanel"] = gridpanel
	m["DataLength"] = len(pageData.Datas) - 1
	if len(pageData.Datas) > 0 {
		m["FieldLength"] = len(pageData.Datas[0].(*Model).Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_page.json", m, template.FuncMap{"getPropValue": GetPropValue, "compareInt": CompareInt, "compareString": CompareString, "dealJsonString": DealJsonString}, "../lessgo/template/entity_page.json")
}

//处理实体的所有数据ajax请求
func dealEntityAllData(entity Entity, w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的所有数据ajax请求")

	err := r.ParseForm()

	models, err := findAllData(entity)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["Models"] = models
	m["Entity"] = entity
	m["DataLength"] = len(models) - 1
	if len(models) > 0 {
		m["FieldLength"] = len(models[0].Props) - 1
	}

	commonlib.RenderTemplate(w, r, "entity_alldata.json", m, template.FuncMap{"getPropValue": GetPropValue, "compareInt": CompareInt}, "../lessgo/template/entity_alldata.json")
}

//处理实体的分页ajax请求
func dealEntityLoad(entity Entity, w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	employee := GetCurrentEmployee(r)

	if employee.UserId == "" {
		Log.Warn("用户未登陆")
		m["success"] = false
		m["code"] = 100
		m["msg"] = "用户未登陆"
		commonlib.OutputJson(w, m, " ")
		return
	}

	Log.Debug("路径：", r.URL.Path, "访问实体", entity.Id, "的load单实体ajax请求")

	vars := mux.Vars(r)
	id := vars["id"]

	model, err := findById(entity, id)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:" + err.Error()
		commonlib.OutputJson(w, m, " ")
		return
	}

	if model.Id == 0 {
		m["success"] = false
		m["code"] = 100
		m["msg"] = "出现错误，请联系IT部门，错误信息:找不到相应的实体"
		commonlib.OutputJson(w, m, " ")
		return
	}

	m["Entity"] = entity
	m["Model"] = model
	m["FieldLength"] = len(entity.Fields) - 1

	commonlib.RenderTemplate(w, r, "entity_load.json", m, template.FuncMap{"getPropValue": GetPropValue, "compareInt": CompareInt}, "../lessgo/template/entity_load.json")
}
