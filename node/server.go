// Package node
// 文件包括本地服务（服务节点）的所有注册服务和远程调用本地服务的相应处理方法
package node

import (
	"context"
	"fmt"
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/log"
	"gitee.com/jyk1987/es/tool"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"time"
)

// _ExecuteService 执行（本地）服务
func _ExecuteService(request *data.Request) (*data.Result, error) {
	// 获取服务
	s := getService(request.Path)
	if s == nil {
		return nil, fmt.Errorf("服务没有找到,path:%v", request.Path)
	}
	// 获取方法
	m := s.GetMethod(request.Method)
	if m == nil {
		return nil, fmt.Errorf("方法没有找到,path:%v,method:%v", request.Path, request.Method)
	}
	// 执行方法
	return m.Execute(request)
}

// ESNode rpcx暴露的服务
type ESNode int

// RpcGetInfoName 获取节点信息
const RpcGetInfoName = "GetInfo"

// RpcExecuteFuncName 执行服务
const RpcExecuteFuncName = "Execute"

// Execute 执行服务
func (*ESNode) Execute(ctx context.Context, request *data.Request, result *data.Result) error {
	r, err := _ExecuteService(request)
	if err != nil {
		return err
	}
	result.SetData(r.GetData())
	return nil
}

// GetInfo 获取信息
func (*ESNode) GetInfo(ctx context.Context, request *data.Request, result *data.Result) error {
	r, err := _ExecuteService(request)
	if err != nil {
		return err
	}
	result.SetData(r.GetData())
	return nil
}

// InitESConfig 初始化rpc服务端
func InitESConfig(configFile ...string) error {
	cfg, e := data.GetConfig(configFile...)
	if e != nil {
		return e
	}
	_Config = cfg
	return nil
}

func StartNodeServer() {
	s := server.NewServer()
	addRegistryPlugin(s)
	s.RegisterName(GetNodeConfig().Name, new(ESNode), "")
	s.Serve("tcp", fmt.Sprintf("0.0.0.0:%v", GetNodeConfig().Port))
}

func addRegistryPlugin(s *server.Server) {
	ip, _ := tool.GetOutBoundIP()
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%v:%v", ip, GetNodeConfig().Port),
		EtcdServers:    []string{GetNodeConfig().Etcd},
		BasePath:       data.ETCDBasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Second,
	}
	err := r.Start()
	if err != nil {
		log.Log.Fatal(err)
	}
	s.Plugins.Add(r)
}
