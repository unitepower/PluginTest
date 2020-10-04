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
