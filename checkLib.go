// Title：静态资源更新
//
// Description:
//
// Author:samurai
//
// Createtime:2013-09-26 11:24
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//
// 1.0 2013-09-26 11:24 samurai 创建文档package checkLib
package lessgo

import (
	"os"
	"os/exec"
)

//检查文件是否存在
func isExists(filePath string) bool{
	fi, _ := os.Stat(filePath)
	return fi != nil
}

//下载
func downLoad(url, libName string) error{
	cmd := exec.Command("/usr/local/bin/wget", "-c", url, "-O", libName)
	Log.Debug(cmd)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

//解压
func decompressionZip(libName string) error{
	cmd := exec.Command("unzip", libName)
	Log.Debug(cmd)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func checkLib() error{
	//下载到本地后的包名
	libName := "lessgo.zip"
	//下载url
	url,_ := Config.GetValue("lessgo", "staticZipUrl")

	//不存在
	if isExists(libName) {
		Log.Debug("资源包文件不存在，开始下载")
		//下载
		err := downLoad(url, libName)

		if err!= nil{
			return err
		}else{
			err = decompressionZip(libName)

			if err!= nil{
				return err
			}
		}

		return nil
	}

	return nil
}
