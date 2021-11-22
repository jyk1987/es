// Package node
// 文件包括本地服务（服务节点）的所有注册服务和远程调用本地服务的相应处理方法
package node

import (
	"context"
	"fmt"
	"github.com/jyk1987/es/data"
	"github.com/jyk1987/es/log"
	"github.com/jyk1987/es/tool"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"os"
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

// RpcGetInfoFuncName 获取节点信息
const RpcGetInfoFuncName = "GetInfo"

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
	result.NodeINfo = GetLocalNodeInfo()
	return nil
}

// InitESConfig 初始化rpc服务端
func InitESConfig(configFile ...string) error {
	cfg, e := data.GetConfig(configFile...)
	if e != nil {
		return e
	}
	_Config = cfg
	os.Setenv("ETCDCTL_API", "3")
	return nil
}

func StartNodeServer() {
	cfg := GetNodeConfig()
	s := server.NewServer()
	e := addRegistryPlugin(s)
	if e != nil {
		log.Log.Errorf("Etcd Register error:%v", e.Error())
		return
	}
	e = s.RegisterName(cfg.Name, new(ESNode), "")
	if e != nil {
		log.Log.Errorf("ESNode Register error:%v", e.Error())
		return
	}
	for {
		log.Log.Info("ESNode starting...")
		e = s.Serve("tcp", fmt.Sprintf("0.0.0.0:%v", cfg.Port))
		if e != nil {
			log.Log.Error("ESNode start fail:", e)
		}
		log.Log.Infof("wait 30 seconds reconnecting.")
		time.Sleep(time.Second * 30)
	}
}

func addRegistryPlugin(s *server.Server) error {
	cfg := GetNodeConfig()
	var endpoint string
	if len(cfg.Endpoint) > 0 {
		endpoint = fmt.Sprintf("tcp@%v", cfg.Endpoint)
	} else {
		local, _ := tool.GetOutBoundIP()
		endpoint = fmt.Sprintf("tcp@%v:%v", local, GetNodeConfig().Port)
	}
	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: endpoint,
		ConsulServers:  []string{GetNodeConfig().Consul},
		BasePath:       data.ETCDBasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Second * 10,
	}
CONNETCD:
	err := r.Start()
	if err != nil {
		log.Log.Errorf("Consul connect error:%v:%v", endpoint, err.Error())
		time.Sleep(time.Second * 5)
		goto CONNETCD
	}
	s.Plugins.Add(r)
	return nil
}
