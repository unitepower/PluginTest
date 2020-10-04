package main

import (
	"PluginTest"
	"fmt"
	"time"
)

/*
@author: mxd
@create time: 2020/10/5
*/

// 存储所有插件信息
var PluginItems []PluginTest.PluginItem

// main 主体程序入口
func main() {

	// time.Now().Second 将会返回当前秒数

	PluginTest.DoInvokePlugin(PluginItems, PluginTest.GetTimeActive, true)
	nowSecond := time.Now().Second()
	PluginTest.DoInvokePlugin(PluginItems, PluginTest.GetTimeActive, false)

	PluginTest.DoInvokePlugin(PluginItems, PluginTest.DoPrintActive, true)
	doPrint(nowSecond)
	PluginTest.DoInvokePlugin(PluginItems, PluginTest.DoPrintActive, false)
	fmt.Println("Process Stop ========")
}

// 执行打印操作
func doPrint(nowSecond int) {
	PluginTest.DoInvokePlugin(PluginItems, PluginTest.PrintItemActive, true)
	if nowSecond%2 == 0 {
		printWorld() //偶数
	} else {
		printHello() //奇数
	}
	PluginTest.DoInvokePlugin(PluginItems, PluginTest.PrintItemActive, false)
}

// 执行打印hello
func printHello() {
	fmt.Println("hello")
}

// 执行打印world
func printWorld() {
	fmt.Println("world")
}

// init 函数将于 main 函数之前运行
func init() {

	// 读取plugin文件夹
	pluginsFiles := PluginTest.FindFile("plugins")

	// 装载插件
	PluginItems = PluginTest.LoadAllPlugin(pluginsFiles)

	fmt.Println("Process On ==========")
}
