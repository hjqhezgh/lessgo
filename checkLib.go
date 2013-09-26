// Title：
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
	"fmt"
	"os"
	"os/exec"
)

//检查文件是否存在
func IsExists(filePath string) bool{
	fi, _ := os.Stat(filePath)
	return fi != nil
}

//下载
func DownLoad(url, libName string) bool{
	cmd := exec.Command("/usr/local/bin/wget", "-c", url, "-O", libName)
	fmt.Println(cmd)
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println(buf)
		fmt.Println(err)
		return false
	}
	return true
}

//解压
func DecompressionZip(libName string) bool{

	cmd := exec.Command("unzip", libName)
	fmt.Println(cmd)
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println(buf)
		fmt.Println(err)
		return false
	}
	return true
}

func CheckLib() {
	//下载到本地后的包名
	libName := "liuli.zip"
	//下载url
	url := "https://github.com/lauly/Demo/blob/master/liuli.zip?raw=true"
	//不存在
	if !IsExists("liuli.zip") {
		fmt.Println("文件不存在!")
		//下载
		if DownLoad(url, libName) {
			//解压
			fmt.Println("下载文件!")
			if DecompressionZip(libName) {
				fmt.Println("解压成功!")
			} else {
				fmt.Println("解压失败!")
			}
		}
	}
}
func main() {
	CheckLib()
}

