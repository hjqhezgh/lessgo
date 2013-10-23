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
	"github.com/Unknwon/goconfig"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/moovweb/log4go"
	"io/ioutil"
	"net/http"
)

var (
	tmplog     log4go.Logger
	Log        *MyLogger //提供公用的日志方式
	Config     *goconfig.ConfigFile
	entityList entitys
	urlList    urls
	Store      *sessions.CookieStore
	SiteName   string
	SiteIcon   string
)

func init() {
	Config, _ = goconfig.LoadConfigFile("../etc/config.ini")

	logFilePath, _ := Config.GetValue("lessgo", "logFilePath")

	tmplog = make(log4go.Logger)
	tmplog.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

	//单位是字节
	fw := log4go.NewFileLogWriter(logFilePath, false).SetRotateSize(10 * 1024 * 1024).SetRotate(true)
	tmplog.AddFilter("log", log4go.INFO, fw)
	Log = new(MyLogger)

	cookieSecret, _ := Config.GetValue("lessgo", "cookieSecret")
	Store = sessions.NewCookieStore([]byte(cookieSecret))

	SiteName,_ = Config.GetValue("lessgo", "siteName")
	SiteIcon,_ = Config.GetValue("lessgo", "siteIcon")
}

//解析配置文件内容至内存中
func analyse() error {

	content, err := ioutil.ReadFile("../etc/entity.xml")
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
func ConfigLessgo() *mux.Router {

	err := analyse()

	if err != nil {
		panic(err)
	}

	err = checkLib()
	if err != nil {
		panic(err)
	}

	http.Handle("/lessgo/", http.FileServer(http.Dir("../")))
	http.Handle("/tmp/", http.FileServer(http.Dir("../")))
	http.Handle("/imageupload/", http.FileServer(http.Dir("../")))

	r := mux.NewRouter()

	r.HandleFunc("/", homeAction)

	r.HandleFunc("/loginOut", loginOutAction)

	//这里的把每个实体的url规约好，暂时不去改变，将来再考虑配置 FIXME

	for _, terminal := range urlList.Terminals {

		r.HandleFunc("/"+terminal, commonAction)
		r.HandleFunc("/"+terminal+"/index.html", commonAction)
		r.HandleFunc("/"+terminal+"/login.html", loginAction)

		for _, entity := range entityList.Entitys {
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/index.html", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id, commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/{id:[0-9]+}", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/add", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/modify", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/save", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/page", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/alldata", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/load"+"/{id:[0-9]+}", commonAction)
			r.HandleFunc("/"+terminal+"/"+entity.Id+"/delete", commonAction)
		}

	}

	for _, url := range urlList.Urls {
		r.HandleFunc(url.Path, independentAction)
	}

	r.HandleFunc("/queryMenus", QueryMenusAction)

	r.HandleFunc("/region/regions", regions)

	r.HandleFunc("/timedim/years", years)
	r.HandleFunc("/timedim/months", months)
	r.HandleFunc("/timedim/weeks", weeks)

	r.HandleFunc("/mutisave", mutiSavaAction)

	r.HandleFunc("/imgageuplaod", imageUpload)

	r.HandleFunc("/kindeditorImageUpload", kindeditorImageUpload)

	//	http.Handle("/", r)

	fmt.Println("lessgo配置完成")

	return r
}
