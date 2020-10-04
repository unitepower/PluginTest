# Golang 插件化开发
Golang官方提供了`plugin`模块,该模块可以支持插件开发.

目前很多思路都是在开发过程中支持插件话,当主体程序写完后,不能够临时绑定插件.但是本文将带领你进行主体程序自动识别并加载、控制插件调用.

github地址

## 基本思路
插件化开发中,一定存在一个主体程序,对其他插件进行控制、处理、调度.

### 具有模拟业务的主体程序
我们首先开发一个简单的业务程序,进行两种输出.

1. 当时间秒数为`奇数`的时候,输出`hello`
2. 当时间秒数为`偶数`的时候,输出`world`
###### 主体代码
代码有一定的冗余,是为了模拟业务之间的调度
	
	主文件名:MainFile.go
```go
package main

import (
	"fmt"
	"time"
)

/*
@author: mxd
@create time: 2020/10/5
*/

// main 主体程序入口
func main() {

	// time.Now().Second 将会返回当前秒数
	nowSecond := time.Now().Second()
	doPrint(nowSecond)
	fmt.Println("Process Stop ========")
}

// 执行打印操作
func doPrint(nowSecond int) {
	if nowSecond%2 == 0 {
		printWorld() //偶数
	} else {
		printHello() //奇数
	}
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
	fmt.Println("Process On ==========")
}

```

输出如下:![主体代码运行截图](https://img-blog.csdnimg.cn/20201005042724188.png#pic_center)

### 简单的插件
然后我们编写一个插件代码
插件代码的入口`package`也要为`main`,但是可以不包含`main`方法

设定插件逻辑为当当前秒数为奇数的时候,同时输出当前时间(与hello的判定不是一个时间)

	插件文件名:HelloPlugin.go
###### 插件代码
```go
package main

import (
	"fmt"
	"time"
)

/*
@author: mxd
@create time: 2020/10/5
*/

// 打印当前时间
func PrintNowTime(){
	fmt.Println(time.Now().Second())
}
```
在当前目录下,执行插件生成指令:
```shell
$ go build --buildmode=plugin -o HelloPlugin.so HelloPlugin.go
```
当前目录下就会多出来一个文件`HelloPlugin.so`
然后,我们让主程序加载该插件
###### 修改的主体代码
```go
package main

import (
	"fmt"
	"plugin"
	"time"
)

/*
@author: mxd
@create time: 2020/10/5
*/

// main 主体程序入口
func main() {

	// time.Now().Second 将会返回当前秒数
	nowSecond := time.Now().Second()
	doPrint(nowSecond)
	fmt.Println("Process Stop ========")
}

// 执行打印操作
func doPrint(nowSecond int) {
	if nowSecond%2 == 0 {
		printWorld() //偶数
	} else {
		printHello() //奇数
	}
}

// 执行打印hello
func printHello() {
	// 执行插件调用
	if pluginFunc != nil{
		//将存储的信息转换为函数
		if targetFunc, ok := pluginFunc.(func()); ok {
			targetFunc()
		}
	}
	fmt.Println("hello")
}

// 执行打印world
func printWorld() {
	fmt.Println("world")
}

// 定义插件信息
const pluginFile = "HelloPlugin.so"

// 存储插件中将要被调用的方法或变量
var pluginFunc plugin.Symbol

// init 函数将于 main 函数之前运行
func init() {

	// 查找插件文件
	pluginFile, err := plugin.Open(pluginFile)

	if err != nil {
		fmt.Println("An error occurred while opening the plug-in")
	} else{
		// 查找目标函数
		targetFunc, err := pluginFile.Lookup("PrintNowTime")
		if err != nil {
			fmt.Println("An error occurred while search target func")
		}

		pluginFunc = targetFunc
	}

	fmt.Println("Process On ==========")
}

```

运行效果如下
![插件化代码](https://img-blog.csdnimg.cn/2020100504380360.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQyMDM4NDA3,size_16,color_FFFFFF,t_70#pic_center)
如上,我们的主体文件已经写好,我们不需要再修改生成后的可执行文件,如果需要扩展代码,仅需要修改插件代码,然后生成`so`文件替换即可.

## 插件进阶-批量化
### 批量化
我们需要考虑到一个问题,如果我们要支持很多的插件,一个一个写的化,很容易导致我们的主体文件膨胀,因为我们将插件文件写死,无法完成自动识别,因此,我们要为主体文件提供自动识别的功能,自动加载插件

#### 自动读取文件夹下的插件
我们可以单独设置一个名为`plugins`的文件夹来保存所有插件
首先我么在项目根目录创建一个文件夹`plugins`
我们将刚刚写好的插件代码移动到 `plugins`文件夹下,同时为了符合**golang标准布局**,我们将主文件移动到`cmd`文件夹下. 此时项目目录如下:
![项目目录结构](https://img-blog.csdnimg.cn/20201005045004292.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQyMDM4NDA3,size_16,color_FFFFFF,t_70#pic_center)
然后,在项目跟(与cmd、plugins同级)目录下新建一个文件,用来处理与业务无关的`util.go`代码.
```go
package PluginTest

import (
	"fmt"
	"io/ioutil"
	"path"
)

/*
@author: mxd
@create time: 2020/10/5
*/

//FindFile 将会打开指定目录，并返回该目录下的所有文件
func FindFile(directoryPath string) []string {
	// 尝试打开文件夹
	baseFile, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("An error occurred while open file :[" + directoryPath + "] .")
		fmt.Println(err)
		return nil
	}

	// 定义返回数据
	var res []string

	for _, fileItem := range baseFile {

		// 文件夹类型继续递归查找
		if fileItem.IsDir() {
			// 加上前缀路径，合成正确的相对或绝对路径
			innerFiles := FindFile(path.Join(directoryPath, fileItem.Name()))
			// 合并结果集
			res = append(res, innerFiles...)
		} else {
			// 这里可以添加过滤，但是会提高方法的复杂度
			/*
				if path.Ext(fileItem.Name()) == ".so"{
					...
				}
			*/
			res = append(res, path.Join(directoryPath, fileItem.Name()))
		}
	}

	return res
}
```
然后修改MainFile文件,让主文件读取插件文件夹
###### MainFile更新代码如下
```go
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

// main 主体程序入口
func main() {

	// time.Now().Second 将会返回当前秒数
	nowSecond := time.Now().Second()
	doPrint(nowSecond)
	fmt.Println("Process Stop ========")
}

// 执行打印操作
func doPrint(nowSecond int) {
	if nowSecond%2 == 0 {
		printWorld() //偶数
	} else {
		printHello() //奇数
	}
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

	for _, pluginItem := range(pluginsFiles){
		fmt.Println(pluginItem)
	}

	fmt.Println("Process On ==========")
}

```
运行结果如下
![在这里插入图片描述](https://img-blog.csdnimg.cn/20201005052846829.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQyMDM4NDA3,size_16,color_FFFFFF,t_70#pic_center)

#### 插件装载
插件装载很简单,但是让插件运行需要们指定一个函数,所有插件都要必须实现该方法,但是如果批量后,我们无法确定插件的运行时机,因此我们会在装载后,直接运行插件,测试我们的批量装载是可行的.
首先我们需要创建一个单独处理插件的文件`pluginSupport.go`
###### 插件支持代码如下
```go
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
	Name       string
	TargetFunc plugin.Symbol
}

// 所有插件必须实现该方法
const TargetFuncName = "TargetFunc"

// LoadAllPlugin 将会过滤一次传入的targetFile,同时将so后缀的文件装载，并返回一个插件信息集合
func LoadAllPlugin(targetFile []string) []PluginItem {
	var res []PluginItem

	for _, fileItem := range targetFile {
		// 过滤插件文件
		if path.Ext(fileItem) == "so" {
			pluginFile, err := plugin.Open(fileItem)
			if err != nil {
				fmt.Println("An error occurred while load plugin : [" + fileItem + "]")
				fmt.Println(err)
			}

			//查找指定函数或符号
			targetFunc, err := pluginFile.Lookup(TargetFuncName)
			if err != nil {
				fmt.Println("An error occurred while search target func : [" + fileItem + "]")
				fmt.Println(err)
			}

			//采集插件信息
			pluginInfo := PluginItem{
				Name:       fileItem,
				TargetFunc: targetFunc,
			}

			// 进行调用
			if f, ok := targetFunc.(func()); ok {
				f()
			}

			res = append(res, pluginInfo)
		}
	}
	return res
}

```

修改`main`函数,使主函数支持该逻辑调用
```go
// init 函数将于 main 函数之前运行
func init() {

	// 读取plugin文件夹
	pluginsFiles := PluginTest.FindFile("plugins")

	// 装载插件
	pluginItems := PluginTest.LoadAllPlugin(pluginsFiles)
	fmt.Println(pluginItems)

	fmt.Println("Process On ==========")
}
```

修改插件代码,使其具有`TargetFunc`方法
```go
package main

import (
	"fmt"
	"time"
)

/*
@author: mxd
@create time: 2020/10/5
*/

func TargetFunc(){
	PrintNowTime()
}

// 打印当前时间
func PrintNowTime(){
	fmt.Println(time.Now().Second())
}
```

生成so文件
```shell
$ cd plugins
$ go build --buildmode=plugin -o HelloPlugin.so HelloPlugin.go 
```
然后运行主文件
![批量插件化开发](https://img-blog.csdnimg.cn/20201005053250731.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQyMDM4NDA3,size_16,color_FFFFFF,t_70#pic_center)
#### 增加插件进行测试
增加插件,但是不修改主代码逻辑
该插件实现逻辑如下:
1. 当前时间秒数\< 30 : 打印`你`
2. 当前时间秒数\>=30: 打印`好`

创建文件:`TestPlugin.go`

由于两个代码中都含有同样的方法`TargetFunc`,编辑器会报错,所以将`HelloPlugin.go`文件中的相关代码注释掉即可(不执行go build命令)

> 实际开发过程中,插件和主体程序是不会混在一起的,但是这里考虑方便才写到一起的

代码如下:
```go
package main

import (
	"fmt"
	"time"
)

/*
@author: mxd
@create time: 2020/10/5
*/

func TargetFunc() {
	nowSecond := time.Now().Second()
	if nowSecond % 2 == 0{
		fmt.Println("好")
	} else {
		fmt.Println("你")
	}
}
```
然后生成so文件:
```shell
$ cd plugins
$ go build --buildmode=plugin -o TestPlugin.so TestPlugin.go
```
然后,不用修改任何代码,直接运行主文件
![在这里插入图片描述](https://img-blog.csdnimg.cn/20201005054132606.png#pic_center)
可以看到,插件已经加载成功,并被执行

致此,我们已经能够自动加载、执行插件了.

## 插件进阶-流程控制、原程序的方法调用
### 流程控制
上一进阶最后面,我们发现了一个问题,我们只能调用一个方法,而且无法控制插件的调用时机,那么我们在插件中,写入一些信息,让主体程序识别,然后在合适的时候进行调用.
首先,我们声明一个插件信息结构体,所有插件填写正确的插件信息,才能被调用

我们的主体应用流程如下:
1. 获取当前时间
2. 调用打印函数进行打印
3. 分别打印

所以我们定义如下插件共享信息
```go
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
```
我们首先修改第一个插件信息(HelloPlugin.go)
代码如下:
```go
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
	Name:        "hello plugin",
	ActiveFlag:  PluginTest.PrintItemActive,
	ActivePoint: true,           //目标任务执行前运行
	Functions:   "PrintNowTime", //可用函数名
}

// 打印当前时间
func PrintNowTime() {
	fmt.Println(time.Now().Second())
}
```
修改主体程序,使其支持插件运行控制

首先修改PluginSupport.go,使其能够获取插件更多的信息,同时添加一个方法,控制调用流程
```go
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
```

修改主文件,添加流程定义,当然,你可以利用上下文让流程定更优雅些
```go
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
```

编译HelloPlugin.go插件

运行主文件
![在这里插入图片描述](https://img-blog.csdnimg.cn/20201005062248596.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQyMDM4NDA3,size_16,color_FFFFFF,t_70#pic_center)
### 调用原有程序的API
直接import原代码,然后调用即可
我们尝试利用FindFile函数,输出当前目录下的所有文件
TestPlugin.go代码:
```go
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

```
编译、运行主文件
![在这里插入图片描述](https://img-blog.csdnimg.cn/2020100506300982.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQyMDM4NDA3,size_16,color_FFFFFF,t_70#pic_center)
我们在此过程中,并未修改主要逻辑代码,却对行为进行了修改.

## 插件进阶-传参、更加优雅的调用
这里已经偏离的插件的控制了,已经涉及到模块的控制了.

对于不通的项目,已经无法依靠小幅度修改进行适配了,因此这里仅提供几种思路,不提供具体的逻辑实现.
### 1. 上下文
依靠程序上下文,我们可以做很多事情
我们让所有插件都接受上下文作为参数,而上下文对于插件和主体程序是共享的,因此可以依靠上下文传递变量、或者更多的信息.

日志体系完全可以依靠这种方式植入,同时,插件能够控制更多的行为和数据.

我们还可以依靠上下文控制插件能够调用的方法.
我们在每一次调用方法的时候,使用包装器或者其他手段让上下文自动更新,而上下文更新的同时去调用插件,这样,我们就和插件降耦了,而且,本身上下文也可以作为参数,提供给程序主体进行调用控制,所以我们是和上下文耦合的.

### 2. 参数写死
这样做的好处是,快速开发,如果我们按照方法1的方式进行开发,整个应用会变得特别臃肿:上下文、插件、流程、静态变量等众多模块将会被引入.
但是缺点也显而易见,不论是主体应用还是插件本身的维护成本很高.

