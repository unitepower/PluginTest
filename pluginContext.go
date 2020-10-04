package PluginTest

/*
@author: mxd
@create time: 2020/10/5
*/
// 定义主体程序流程

const (
	GetTimeActive   = "get_time_active"   //获取时间的流程
	DoPrintActive   = "do_print_active"   //执行打印的流程
	PrintItemActive = "print_item_active" //执行分别打印
)

// 存储插件信息
type PluginBaseInfo struct {
	Name        string // 插件名称
	ActiveFlag  string // 插件执行的位置
	ActivePoint bool   // 插件的执行点
	Functions   string // 插件可用函数
	// 对于可用函数，可以写为数组，来暴露更多的方法，使用一些信息标注调用时间
	// 使用更细的粒度控制插件暴露的API
}
