// Title：逻辑流程
//
// Description:
//
// Author:black
//
// Createtime:2013-08-06 14:15
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-06 14:15 black 创建文档
package lessgo

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/moovweb/log4go"
	"io/ioutil"
	"net/http"
)

var (
	tmplog  log4go.Logger
	Log     *MyLogger
	config  Config
	entityList entitys
	navList    navs
	urlList    urls
)

func init() {
	tmplog = make(log4go.Logger)
	tmplog.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())
	//单位是字节
	fw := log4go.NewFileLogWriter("code.log", false).SetRotateSize(10 * 1024 * 1024).SetRotate(true)
	tmplog.AddFilter("log", log4go.INFO, fw)
	Log = new(MyLogger)
}

//解析配置文件内容至内存中
func Analyse() error {

	content, err := ioutil.ReadFile("../etc/config.xml")

	if err != nil {
		Log.Error(err)
		return err
	}

	err = xml.Unmarshal(content, &config)
	if err != nil {
		Log.Error(err)
		return err
	}

	content, err = ioutil.ReadFile("../etc/entity.xml")
	err = xml.Unmarshal(content, &entityList)

	if err != nil {
		Log.Error(err)
		return err
	}

	content, err = ioutil.ReadFile("../etc/url.xml")
	err = xml.Unmarshal(content, &urlList)

	if err != nil {
		Log.Error(err)
		return err
	}

	return nil
}

//启动应用
func ConfigServer() {

	http.Handle("/js/", http.FileServer(http.Dir("../static")))
	http.Handle("/css/", http.FileServer(http.Dir("../static")))
	http.Handle("/img/", http.FileServer(http.Dir("../static")))
	http.Handle("/tmp/", http.FileServer(http.Dir("../")))
	http.Handle("/imageupload/", http.FileServer(http.Dir("../")))

	r := mux.NewRouter()

	r.HandleFunc("/", homeAction)

	//这里的把每个实体的url规约好，暂时不去改变，将来再考虑配置 FIXME

	for _, terminal := range urlList.Terminals {

		r.HandleFunc("/"+terminal, commonAction)
		r.HandleFunc("/"+terminal+"/index.html", commonAction)

		for _, entity := range entityList.Entitys {
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/index.html", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id, commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/{id:[0-9]+}", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/add", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/modify"+"/{id:[0-9]+}", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/save", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/page", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/alldata", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/load"+"/{id:[0-9]+}", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/delete"+"/{id:[0-9]+}", commonAction)
		}

	}

	for _, url := range urlList.Urls {
		r.HandleFunc(url.Path, independentAction)
	}

	r.HandleFunc("/region/regions", regions)

	r.HandleFunc("/timedim/years", years)
	r.HandleFunc("/timedim/months", months)
	r.HandleFunc("/timedim/weeks", weeks)

	r.HandleFunc("/mutisave", mutiSavaAction)

	r.HandleFunc("/imgageuplaod", imageUpload)

	r.HandleFunc("/kindeditorImageUpload", kindeditorImageUpload)

	http.Handle("/", r)

	fmt.Println("Server start on:", config.Port)

	//Log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
