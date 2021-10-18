package data

import (
	"errors"
	"reflect"
)

// Request 调用服务的时候发出去的数据
type Request struct {
	Node   string        // 节点名称
	Path   string        //服务包路径
	Method string        //服务名
	Args   []interface{} //调用参数
}

// Result 服务执行结果
type Result struct {
	//Error error //执行出错内容，此错误不是远程方法返回的错误，而是服务调用过程出错，或者远程方法执行报错（非正常执行错误）
	Returns []reflect.Value // 方法返回的数据
}

func (r *Result) GetResult(funcInstance interface{}) error {

	fo := reflect.ValueOf(funcInstance)
	if fo.Kind() != reflect.Func {
		return errors.New("参数必须为方法！")
	}
	fo.Call(r.Returns)
	return nil
}
