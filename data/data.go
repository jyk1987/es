package data

import (
	"errors"
	"reflect"
)

// Request 调用服务的时候发出去的数据
type Request struct {
	NodeName string        // 节点名称
	Path     string        //服务包路径
	Method   string        //服务名
	Args     []interface{} //调用参数
}

// Result 服务执行结果
type Result struct {
	//Error error //执行出错内容，此错误不是远程方法返回的错误，而是服务调用过程出错，或者远程方法执行报错（非正常执行错误）
	Returns []interface{} // 方法返回的数据
}

func (r *Result) GetResult(funcInstance interface{}) error {
	ft := reflect.TypeOf(funcInstance)
	if ft.Kind() != reflect.Func {
		return errors.New("必须传入一个方法")
	}

	returnLen := len(r.Returns)
	funcInLen := ft.NumIn()
	if returnLen != funcInLen {
		return errors.New("方法参数个数不同")
	}
	fo := reflect.ValueOf(funcInstance)
	in := make([]reflect.Value, returnLen)
	for i := 0; i < ft.NumIn(); i++ {
		value := r.Returns[i]
		typeValue := reflect.ValueOf(value)
		inType := ft.In(i)
		// valueType := value.Type()
		// if inType.Kind() != valueType.Kind() {
		// 	return errors.New(fmt.Sprintf(
		// 		"第%v个参数类型不相符，返回数据类型%v，方法参数类型%v",
		// 		i, valueType.Kind(), inType.Kind(),
		// 	))
		// }
		switch inType.Kind() {
		// TODO:不同类型的转换
		case reflect.Struct:
			in[i] = typeValue.Convert(inType)
		default:
			in[i] = typeValue.Convert(inType)
		}

	}
	fo.Call(in)
	return nil
}
