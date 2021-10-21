package es

import (
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/is"
	"gitee.com/jyk1987/es/log"
	"gitee.com/jyk1987/es/node"
	_ "github.com/gogf/gf"
	"time"
)

var _IndexCache = make(map[string]*is.IndexInfo, 0)

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

func InitNode() error {
	e := node.InitNodeServer()
	if e != nil {
		return e
	}
	e = callServerRegNode()
	if e != nil {
		return e
	}
	return nil
}

func StartNode() {
	go node.StartNodeServer()

	for {
		e := callServerPing()
		if e != nil {
			log.Log.Error("ping 索引服务器失败:", e)
		}
		time.Sleep(time.Second)
	}
}
