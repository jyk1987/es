package data

import (
	"gitee.com/jyk1987/es/log"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

type ESDataType uint

const (
	ESDataNil    ESDataType = 0
	ESDataString ESDataType = 1
	ESDataNumber ESDataType = 2
	ESDataBool   ESDataType = 3
	ESDataJson   ESDataType = 4
	ESDataBinary ESDataType = 5
	ESDataError  ESDataType = 100
)

type ESData struct {
	Type   ESDataType // 数据类型
	Number float64    //存储数值数据，所有的熟知了行不管整数还是浮点都转成float64存储
	Bool   bool       //存储bool数据
	Binary []byte     //其他数据都存储二级制流中
}

func NewESData(value reflect.Value) *ESData {
	d := new(ESData)
	d.Type = ESDataNil
	data := value.Interface()
	switch value.Kind() {
	case reflect.Bool:
		//fmt.Println("转换bool")
		d.Type = ESDataBool
		d.Bool = data.(bool)
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
		//fmt.Println("转换number")
		d.Type = ESDataNumber
		d.Number = data.(float64)
	case reflect.String:
		d.Type = ESDataString
		//fmt.Println("转换string")
		d.Binary = []byte(data.(string))
	case reflect.Interface:
		//fmt.Println("转换interface")
		if !value.IsNil() {
			if value.Type().Name() == "error" {
				d.Type = ESDataError
				d.Binary = []byte(value.MethodByName("Error").Call([]reflect.Value{})[0].String())
			} else {
				d.Type = ESDataJson
				var json = jsoniter.ConfigCompatibleWithStandardLibrary
				d.Binary, _ = json.Marshal(&data)
			}
		}
	case reflect.Struct:
		//fmt.Println("转换struct")
		d.Type = ESDataJson
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		d.Binary, _ = json.Marshal(&data)
	case reflect.Ptr:
		//fmt.Println("转换ptr")
		if !value.IsNil() {
			d.Type = ESDataJson
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			d.Binary, _ = json.Marshal(&data)
		}
	case reflect.Slice:
		if !value.IsNil() {
			if value.Type().String() == "[]uint8" {
				//fmt.Println("转换[]byte")
				d.Type = ESDataBinary
				d.Binary = value.Bytes()
			} else {
				//fmt.Println("转换Slice")
				d.Type = ESDataJson
				var json = jsoniter.ConfigCompatibleWithStandardLibrary
				d.Binary, _ = json.Marshal(&data)
			}
		}
	default:
		//fmt.Println("转换other:", value.Kind())
		d.Type = ESDataJson
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		d.Binary, _ = json.Marshal(&data)
	}
	return d
}

// Request 调用服务的时候发出去的数据
type Request struct {
	NodeName   string   // 节点名称
	Path       string   //服务包路径
	Method     string   //服务名
	Parameters [][]byte //调用参数
}

func (r *Request) SetParameters(parameter ...interface{}) error {
	count := len(parameter)
	r.Parameters = make([][]byte, count)
	for i := 0; i < count; i++ {
		b, e := EncodeData(parameter[i])
		if e != nil {
			log.Log.Error("参数转换失败！", parameter[i])
			return e
		}
		r.Parameters[i] = b
	}
	return nil
}
func (r *Request) AddParameter(parameter interface{}) error {
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
	//Error error //执行出错内容，此错误不是远程方法返回的错误，而是服务调用过程出错，或者远程方法执行报错（非正常执行错误）
	Returns []*ESData // 方法返回的数据
}

func (r *Result) GetReturnData(index ...int) *ESData {
	i := 0
	if len(index) > 0 {
		i = index[0]
	}
	count := len(r.Returns)
	if i >= count {
		i = count - 1
	}
	return r.Returns[i]
}

func (r *Result) GetNumber(index ...int) float64 {
	return r.GetReturnData(index...).Number
}

//
//func (r *Result) GetResult(funcInstance interface{}) error {
//	ft := reflect.TypeOf(funcInstance)
//	if ft.Kind() != reflect.Func {
//		return errors.New("必须传入一个方法")
//	}
//
//	returnLen := len(r.Returns)
//	funcInLen := ft.NumIn()
//	if returnLen != funcInLen {
//		return errors.New("方法参数个数不同")
//	}
//	fo := reflect.ValueOf(funcInstance)
//	in := make([]reflect.Value, returnLen)
//	for i := 0; i < ft.NumIn(); i++ {
//		data := r.Returns[i]
//		fmt.Println(data.Type)
//		fmt.Println(data.Data)
//		fmt.Println(data.JsonStr)
//		dValue := reflect.ValueOf(data.Data)
//		//dType := reflect.TypeOf(data.Data)
//
//		inType := ft.In(i)
//		switch inType.Kind() {
//		// TODO:不同类型的转换
//		case reflect.Struct:
//			in[i] = dValue.Convert(inType)
//		case reflect.Interface:
//			//it := dValue
//			var e error = errors.New("")
//			in[i] = reflect.ValueError{}
//		default:
//			in[i] = dValue.Convert(inType)
//		}
//
//	}
//	fo.Call(in)
//	return nil
//}
