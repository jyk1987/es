package es

import (
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/node"
)

const ESVersion = 1

// Reg 注册本地服务
func Reg(serviceInstance interface{}) {
	node.Reg(serviceInstance)
}

func Call(nodeName, path, method string, params ...interface{}) (*data.Result, error) {
	r := &data.Request{
		NodeName: nodeName,
		Path:     path,
		Method:   method,
	}
	e := r.SetParameters(params...)
	if e != nil {
		return nil, e
	}
	return node.ExecuteService(r)
}
