// Title：
//
// Description:
//
// Author:black
//
// Createtime:2013-08-28 09:55
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-08-28 09:55 black 创建文档
package lessgo

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/hjqhezgh/commonlib"
)

func ImageUpload(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	_, err := os.Stat("../tmp")

	if err != nil && os.IsNotExist(err) {
		Log.Info("tmp文件夹不存在，创建")
		os.Mkdir("../tmp", 0777)
	}

	fn, header, err := r.FormFile("pid")

	if err != nil && os.IsNotExist(err) {
		m["success"] = false
		m["code"] = 100
		Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m," ")
		return
	}

	newFileName := findRandomFileName(header.Filename)

	f, err := os.Create("../tmp/" + newFileName)

	if err != nil {
		m["success"] = false
		m["code"] = 100
		Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m," ")
		return
	}

	defer f.Close()

	io.Copy(f, fn)

	m["success"] = true
	m["code"] = 200
	m["tmpfile"] = "/tmp/" + newFileName

	commonlib.OutputJson(w, m," ")
}

/*****
 * 获取上传图片的随机不重复文件名
 */
func findRandomFileName(sourceFileName string) string {

	suffix := commonlib.Substr(sourceFileName, strings.LastIndex(sourceFileName, "."), len(sourceFileName))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	str := ""

	for i := 0; i < 4; i++ {
		str += fmt.Sprint(r.Intn(10))
	}

	return fmt.Sprint(time.Now().UnixNano(), str, suffix)
}

func KindeditorImageUpload(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	_, err := os.Stat("../imageupload")

	if err != nil && os.IsNotExist(err) {
		Log.Info("imageupload，创建")
		os.Mkdir("../imageupload", 0777)
	}

	fn, header, err := r.FormFile("imgFile")

	if err != nil && os.IsNotExist(err) {
		m["error"] = 1
		m["message"] = err.Error()
		Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m," ")
		return
	}

	newFileName := findRandomFileName(header.Filename)

	f, err := os.Create("../imageupload/" + newFileName)

	if err != nil {
		m["error"] = 1
		m["message"] = err.Error()
		Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m," ")
		return
	}

	defer f.Close()

	io.Copy(f, fn)

	m["error"] = 0
	m["url"] = "/imageupload/" + newFileName

	commonlib.OutputJson(w, m," ")
}
