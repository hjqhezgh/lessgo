// Title：扩展下log4go，用来支持文件名和行号
//
// Description:
//
// Author:Black
//
// Createtime:2013-07-22 16:17
//
// Version:1.0
//
// 修改历史:版本号 修改日期 修改人 修改说明
//1.0 2013-05-23 13:32 Black 创建
package lessgo

import (
	"fmt"
	"runtime"
	"strings"
	"github.com/hjqhezgh/commonlib"
)

type MyLogger struct {
}

func (log *MyLogger) Error(arg0 ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = commonlib.Substr(file, strings.LastIndex(file, "/")+1, 100)
	tmplog.Error(fmt.Sprint("[文件：", file, "，行：", line, "] ", fmt.Sprint(arg0...)))
}

func (log *MyLogger) Debug(arg0 ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = commonlib.Substr(file, strings.LastIndex(file, "/")+1, 100)
	tmplog.Debug(fmt.Sprint("[文件：", file, "，行：", line, "] ", fmt.Sprint(arg0...)))
}

func (log *MyLogger) Info(arg0 ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = commonlib.Substr(file, strings.LastIndex(file, "/")+1, 100)
	tmplog.Info(fmt.Sprint("[文件：", file, "，行：", line, "] ", fmt.Sprint(arg0...)))
}

func (log *MyLogger) Warn(arg0 ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = commonlib.Substr(file, strings.LastIndex(file, "/")+1, 100)
	tmplog.Warn(fmt.Sprint("[文件：", file, "，行：", line, "] ", fmt.Sprint(arg0...)))
}
