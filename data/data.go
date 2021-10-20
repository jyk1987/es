package data

import (
	"errors"
	"gitee.com/jyk1987/es/log"
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
	b, e := EncodeData(parameter)
	if e != nil {
		log.Log.Error("参数转换失败！", parameter)
		return e
	}
	r.Parameters = append(r.Parameters, b)
	return nil
}

// Result 服务执行结果
type Result struct {
	data [][]byte // 方法返回的数据
}

func NewResult(vs []reflect.Value) (*Result, error) {
	r := new(Result)
	if vs == nil || len(vs) == 0 {
		return r, nil
	}
	count := len(vs)
	r.data = make([][]byte, count)
	for i := 0; i < count; i++ {
		b, e := EncodeDataByType(vs[i])
		if e != nil {
			log.Log.Error("增加返回数据出错：", e)
			return r, e
		}
		r.data[i] = b
	}
	return r, nil
}

func (r *Result) GetData() [][]byte {
	return r.data
}
func (r *Result) SetData(data [][]byte) {
	r.data = data
}

// AddData 增加数据
func (r *Result) AddData(typeValue reflect.Value) error {
	if r.data == nil {
		r.data = make([][]byte, 0)
	}
	b, e := EncodeDataByType(typeValue)
	if e != nil {
		log.Log.Error("增加返回数据出错：", e)
		return e
	}
	r.data = append(r.data, b)
	return nil
}

// GetResult 利用回调方法获取数据
func (r *Result) GetResult(function interface{}) error {
	ft := reflect.TypeOf(function)
	if ft.Kind() != reflect.Func {
		return errors.New("必须传入一个方法")
	}
	dataLen := len(r.data)
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
		obj, e := DecodeDataByType(r.data[i], inType)
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
