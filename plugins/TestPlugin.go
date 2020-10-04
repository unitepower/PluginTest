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


var PluginBaseInfo = PluginTest.PluginBaseInfo{
	Name:        "test plugin",
	ActiveFlag:  PluginTest.DoPrintActive,
	ActivePoint: false,           //目标任务执行前运行
	Functions:   "TargetFunc", //可用函数名
}

func TargetFunc() {
	nowSecond := time.Now().Second()
	if nowSecond%2 == 0 {
		fmt.Println("好")
	} else {
		fmt.Println("你")
	}
	files := PluginTest.FindFile("./")
	for _, fileItem := range files {
		fmt.Println(fileItem)
	}
}
