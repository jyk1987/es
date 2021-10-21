package tool

import (
	"bytes"
	"errors"
	"gitee.com/jyk1987/es/log"
	"github.com/modern-go/reflect2"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

var _Codec = new(MsgpackCodec)

// MsgpackCodec uses messagepack marshaler and unmarshaler.
type MsgpackCodec struct{}

// Encode encodes an object into slice of bytes.
func (c MsgpackCodec) Encode(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	// enc.UseJSONTag(true)
	err := enc.Encode(i)
	return buf.Bytes(), err
}

// Decode decodes an object from slice of bytes.
func (c MsgpackCodec) Decode(data []byte, i interface{}) error {
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	// dec.UseJSONTag(true)
	err := dec.Decode(i)
	return err
}

func EncodeData(data interface{}) ([]byte, error) {
	return _Codec.Encode(data)
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
	return _Codec.Encode(i)
}

func DecodeData(data []byte, i interface{}) error {
	return _Codec.Decode(data, i)
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
