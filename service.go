package es

import (
	"errors"
	"github.com/jyk1987/es/data"
	"github.com/jyk1987/es/log"
	"github.com/jyk1987/es/node"
)

// Reg 注册本地服务
func Reg(servicePath string, serviceInstance interface{}) {
	node.Reg(servicePath, serviceInstance)
}

// Call 调用服务
func Call(nodeName, path, method string, params ...interface{}) (*data.Result, error) {
	if node.GetNodeConfig() == nil {
		return nil, errors.New("conf not found")
	}
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

// GetNodeInfo 获取节点信息
func GetNodeInfo(nodeName string) (*data.Result, error) {
	if node.GetNodeConfig() == nil {
		return nil, errors.New("conf not found")
	}
	return callServiceGetInfo(nodeName)
}

// InitES 初始化es
func InitES(configFile ...string) error {
	e := node.InitESConfig(configFile...)
	if e != nil {
		return e
	}
	return nil
}

// StartNode 启动服务节点，此方法为阻塞方法，地用后服务会启动，不会有返回
func StartNode() {
	cfg := node.GetNodeConfig()
	if cfg == nil || cfg.Port == 0 {
		return
	}
	log.Log.Info("Start ESNode")
	node.StartNodeServer()
}
