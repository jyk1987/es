package data

import (
	"errors"
	"gitee.com/jyk1987/es/log"
	"gitee.com/jyk1987/es/tool"
	"reflect"
)

// Request 调用服务的时候发出去的数据
type Request struct {
	NodeName   string   // 节点名称
	Path       string   //服务包路径
	Method     string   //服务名
	Parameters [][]byte //调用参数
}

func (r *Request) SetParameters(parameter ...interface{}) error {
	count := len(parameter)
	for i := 0; i < count; i++ {
		e := r.AddParameter(parameter[i])
		if e != nil {
			return e
		}
	}
	return nil
}
func (r *Request) AddParameter(parameter interface{}) error {
	if r.Parameters == nil {
		r.Parameters = make([][]byte, 0)
	}
	b, e := tool.EncodeData(parameter)
	if e != nil {
		log.Log.Error("参数转换失败！", parameter)
		return e
	}
	r.Parameters = append(r.Parameters, b)
	return nil
}

// Result 服务执行结果
type Result struct {
	Binary [][]byte // 方法返回的数据
}

func NewResult(vs []reflect.Value) (*Result, error) {
	r := new(Result)
	if vs == nil || len(vs) == 0 {
		return r, nil
	}
	count := len(vs)
	r.Binary = make([][]byte, count)
	for i := 0; i < count; i++ {
		b, e := tool.EncodeDataByType(vs[i])
		if e != nil {
			log.Log.Error("增加返回数据出错：", e)
			return r, e
		}
		r.Binary[i] = b
	}
	return r, nil
}

func (r *Result) GetData() [][]byte {
	return r.Binary
}
func (r *Result) SetData(data [][]byte) {
	r.Binary = data
}

// AddData 增加数据
func (r *Result) AddData(typeValue reflect.Value) error {
	if r.Binary == nil {
		r.Binary = make([][]byte, 0)
	}
	b, e := tool.EncodeDataByType(typeValue)
	if e != nil {
		log.Log.Error("增加返回数据出错：", e)
		return e
	}
	r.Binary = append(r.Binary, b)
	return nil
}

// GetResult 利用回调方法获取数据
func (r *Result) GetResult(function interface{}) error {
	if r == nil {
		return errors.New("返回数据为nil无法使用")
	}
	ft := reflect.TypeOf(function)
	if ft.Kind() != reflect.Func {
		return errors.New("必须传入一个方法")
	}
	dataLen := len(r.Binary)
	funcInLen := ft.NumIn()
	if dataLen != funcInLen {
		e := errors.New("方法参数个数不同")
		log.Log.Error(e.Error())
		return e
	}
	fo := reflect.ValueOf(function)
	in := make([]reflect.Value, funcInLen)
	for i := 0; i < ft.NumIn(); i++ {
		inType := ft.In(i)
		obj, e := tool.DecodeDataByType(r.Binary[i], inType)
		if e != nil {
			return e
		}
		if obj == nil {
			in[i] = reflect.New(inType).Elem()
		} else {
			in[i] = reflect.ValueOf(obj)
		}
	}
	fo.Call(in)
	return nil
}
