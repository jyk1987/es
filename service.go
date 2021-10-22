package es

import (
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/node"
	_ "github.com/gogf/gf"
	"time"
)

var _IndexCache = make(map[string]*data.IndexInfo, 0)

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
	return callServiceExecute(nodeName, path, method, params...)
}

func InitES() error {
	e := node.InitESConfig()
	if e != nil {
		return e
	}
	return nil
}

func StartNode() {
	go node.StartNodeServer()

	for {
		time.Sleep(time.Minute)
	}
}
