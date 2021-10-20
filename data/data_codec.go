package data

import (
	"errors"
	"gitee.com/jyk1987/es/log"
	"github.com/modern-go/reflect2"
	"github.com/smallnest/rpcx/codec"
	"reflect"
)

func EncodeData(data interface{}) ([]byte, error) {
	return new(codec.MsgpackCodec).Encode(data)
}
func EncodeDataByType(typeValue reflect.Value) ([]byte, error) {
	i := typeValue.Interface()
	if typeValue.Type().String() == "error" && !typeValue.IsNil() {
		e, ok := i.(error)
		if !ok {
			return nil, errors.New("转换错误类型失败")
		}
		ec := e.Error()
		i = ec
		log.Log.Debug("转换错误类型，内容：", ec)
	}
	return new(codec.MsgpackCodec).Encode(i)
}

func DecodeData(data []byte, i interface{}) error {
	return new(codec.MsgpackCodec).Decode(data, i)
}
func DecodeDataByType(data []byte, t reflect.Type) (interface{}, error) {
	i := reflect2.Type2(t).New()
	err := DecodeData(data, i)
	if err != nil {
		return nil, err
	}
	i = reflect.ValueOf(i).Elem().Interface()
	return i, nil
}
