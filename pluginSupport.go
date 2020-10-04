package PluginTest

import (
	"fmt"
	"path"
	"plugin"
)

/*
@author: mxd
@create time: 2020/10/5
*/

// PluginItem 存储着插件的信息
type PluginItem struct {
	Name           string
	PluginBaseInfo PluginBaseInfo
	PluginItem     *plugin.Plugin
}

// 所有插件必须实现该方法
const BaseInfo = "PluginBaseInfo"

// LoadAllPlugin 将会过滤一次传入的targetFile,同时将so后缀的文件装载，并返回一个插件信息集合
func LoadAllPlugin(targetFile []string) []PluginItem {
	var res []PluginItem

	for _, fileItem := range targetFile {
		// 过滤插件文件
		if path.Ext(fileItem) == ".so" {
			pluginFile, err := plugin.Open(fileItem)
			if err != nil {
				fmt.Println("An error occurred while load plugin : [" + fileItem + "]")
				fmt.Println(err)
			}

			//查找指定函数或符号
			targetFunc, err := pluginFile.Lookup(BaseInfo)
			if err != nil {
				fmt.Println("An error occurred while search target info : [" + fileItem + "]")
				fmt.Println(err)
			}

			baseInfo, ok := targetFunc.(*PluginBaseInfo)
			if !ok {
				fmt.Println("Can find base info.")
			}

			//采集插件信息
			pluginInfo := PluginItem{
				Name:           fileItem,
				PluginBaseInfo: *baseInfo,
				PluginItem:     pluginFile,
			}

			res = append(res, pluginInfo)
		}
	}
	return res
}

// DoInvokePlugin 会根据当前状态执行插件调用
func DoInvokePlugin(pluginsItems [] PluginItem, nowActive string, nowPoint bool){
	for _, pluginItem := range pluginsItems{
		// 判断流程
		if pluginItem.PluginBaseInfo.ActiveFlag == nowActive{
			// 判断执行点
			if nowPoint == pluginItem.PluginBaseInfo.ActivePoint{
				funcName := pluginItem.PluginBaseInfo.Functions
				funcItem, err := pluginItem.PluginItem.Lookup(funcName)

				if err != nil{
					fmt.Println("Can't find target func in [" + pluginItem.Name +"].")
					continue
				}

				if f, ok := funcItem.(func()); ok{
					f()
				}
			}
		}
	}
}