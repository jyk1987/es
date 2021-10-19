package data

import (
	"github.com/modern-go/reflect2"
	"github.com/smallnest/rpcx/codec"
	"reflect"
)

func EncodeData(data interface{}) ([]byte, error) {
	return new(codec.MsgpackCodec).Encode(data)
}

func DecodeData(data []byte, i interface{}) error {
	return new(codec.MsgpackCodec).Decode(data, i)
}
func DecodeByType(data []byte, t reflect.Type) (interface{}, error) {
	i := reflect2.Type2(t).New()
	err := DecodeData(data, i)
	if err != nil {
		return nil, err
	}
	i = reflect.ValueOf(i).Elem().Interface()
	return i, nil
}
