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
	"github.com/hjqhezgh/commonlib"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
)

func imageUpload(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]interface{})

	_, err := os.Stat("../tmp")

	if err != nil && os.IsNotExist(err) {
		Log.Info("tmp文件夹不存在，创建")
		os.Mkdir("../tmp", 0777)
	}

	fileInputName := r.FormValue("fileInputName")
	resolutionString := r.FormValue("resolution")
	maxWidthString := r.FormValue("maxWidth")
	maxHeightString := r.FormValue("maxHeight")
	minWidthString := r.FormValue("minWidth")
	minHeightString := r.FormValue("minHeight")
	//	maxSizeString := r.FormValue("maxSize")
	widths := r.FormValue("widths")

	fn, header, err := r.FormFile(fileInputName)
	fn1, header, err := r.FormFile(fileInputName)//读两次，一个用于分析，一个用于io，不然会出现错误

	if err != nil && os.IsNotExist(err) {
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}

	infos := strings.Split(header.Filename, ".")
	suffix := strings.ToLower(infos[len(infos)-1])

	newFileName := findRandomFileName(header.Filename)

	var i image.Image

	if suffix == "jpeg" || suffix == "jpg" {
		i, err = jpeg.Decode(fn)
	} else {
		i, _, err = image.Decode(fn)
	}

	if err != nil && os.IsNotExist(err) {
		m["success"] = false
		m["code"] = 100
		m["msg"] = err.Error()
		Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}

	b := i.Bounds()

	if maxWidthString != "" {
		maxWidth, _ := strconv.Atoi(maxWidthString)
		if b.Dx() > maxWidth {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "图片最大宽度不能超过" + maxWidthString
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if maxHeightString != "" {
		maxHeight, _ := strconv.Atoi(maxHeightString)
		if b.Dy() > maxHeight {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "图片最大高度不能超过" + maxHeightString
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if minWidthString != "" {
		minWidth, _ := strconv.Atoi(minWidthString)
		if b.Dx() < minWidth {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "图片最小宽度不能小于" + minWidthString
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if minHeightString != "" {
		minHeight, _ := strconv.Atoi(minHeightString)
		if b.Dy() < minHeight {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "图片最小高度不能小于" + minHeightString
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if resolutionString != "" {
		resolution, _ := strconv.ParseFloat(resolutionString, 64)
		if resolution != float64(b.Dx())/float64(b.Dy()) {
			m["success"] = false
			m["code"] = 100
			m["msg"] = "图片宽高比应为" + resolutionString
			commonlib.OutputJson(w, m, " ")
			return
		}
	}

	if widths != "" {
		widthsArray := strings.Split(widths, ",")

		tmpFileName := ""
		tmpFileNames := ""

		for index, widthString := range widthsArray {
			var i1 *image.RGBA

			width, _ := strconv.Atoi(widthString)
			height := (b.Dy() * width) / b.Dx()

			i1 = commonlib.Resample(i, b, width, height)
			i128 := commonlib.ResizeRGBA(i1, i1.Bounds(), width, height)

			var buf bytes.Buffer
			if err := png.Encode(&buf, i128); err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = err.Error()
				Log.Error("获取上传图片发生错误，信息如下：", err.Error())
				commonlib.OutputJson(w, m, " ")
				return
			}

			fo, err := os.Create(fmt.Sprint("../tmp/", newFileName, "_", width, ".", suffix))
			if err != nil {
				m["success"] = false
				m["code"] = 100
				m["msg"] = err.Error()
				Log.Error("获取上传图片发生错误，信息如下：", err.Error())
				commonlib.OutputJson(w, m, " ")
				return
			}
			defer fo.Close()
			writer := bufio.NewWriter(fo)
			buf.WriteTo(writer)

			tmpFileNames += fmt.Sprint("/tmp/", newFileName, "_", width, ".", suffix)

			if index < len(widthsArray)-1 {
				tmpFileNames += ","
			}

			if index == 0 {
				tmpFileName = fmt.Sprint(newFileName, "_", width, ".", suffix)
			}
		}

		m["success"] = true
		m["code"] = 200
		m["tmpfile"] = "/tmp/" + tmpFileName
		m["tmpfiles"] = tmpFileNames

		commonlib.OutputJson(w, m, " ")
		return

	} else {/*
		var i1 *image.RGBA

		i1 = commonlib.Resample(i, b, b.Dx(), b.Dy())
		i128 := commonlib.ResizeRGBA(i1, i1.Bounds(), b.Dx(), b.Dy())

		var buf bytes.Buffer
		if err := png.Encode(&buf, i128); err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = err.Error()
			Log.Error("获取上传图片发生错误，信息如下：", err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		fo, err := os.Create("../tmp/" + newFileName + "." + suffix)
		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = err.Error()
			Log.Error("获取上传图片发生错误，信息如下：", err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		defer fo.Close()
		writer := bufio.NewWriter(fo)
		buf.WriteTo(writer)*/

		fo, err := os.Create(fmt.Sprint("../tmp/", newFileName, ".", suffix))
		if err != nil {
			m["success"] = false
			m["code"] = 100
			m["msg"] = err.Error()
			Log.Error("获取上传图片发生错误，信息如下：", err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}
		defer fo.Close()

		_, err = io.Copy(fo, fn1)

		if err != nil && os.IsNotExist(err) {
			m["success"] = false
			m["code"] = 100
			m["msg"] = err.Error()
			Log.Error("获取上传图片发生错误，信息如下：", err.Error())
			commonlib.OutputJson(w, m, " ")
			return
		}

		m["success"] = true
		m["code"] = 200
		m["tmpfile"] = "/tmp/" + newFileName + "." + suffix

		commonlib.OutputJson(w, m, " ")
		return
	}
}

/*****
 * 获取上传图片的随机不重复文件名
 */
func findRandomFileName(sourceFileName string) string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	str := ""

	for i := 0; i < 4; i++ {
		str += fmt.Sprint(r.Intn(10))
	}

	return fmt.Sprint(time.Now().UnixNano(), str)
}

func kindeditorImageUpload(w http.ResponseWriter, r *http.Request) {

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
		commonlib.OutputJson(w, m, " ")
		return
	}

	newFileName := findRandomFileName(header.Filename)

	f, err := os.Create("../imageupload/" + newFileName)

	if err != nil {
		m["error"] = 1
		m["message"] = err.Error()
		Log.Error("获取上传图片发生错误，信息如下：", err.Error())
		commonlib.OutputJson(w, m, " ")
		return
	}

	defer f.Close()

	io.Copy(f, fn)

	m["error"] = 0
	m["url"] = "/imageupload/" + newFileName

	commonlib.OutputJson(w, m, " ")
}
