package main

import (
	"fmt"
	"time"
)

/*
@author: mxd
@create time: 2020/10/5
*/
//
//var PluginBaseInfo = PluginTest.PluginBaseInfo{
//	Name:        "hello plugin",
//	ActiveFlag:  PluginTest.PrintItemActive,
//	ActivePoint: true,           //目标任务执行前运行
//	Functions:   "PrintNowTime", //可用函数名
//}

// 打印当前时间
func PrintNowTime() {
	fmt.Println(time.Now().Second())
}
