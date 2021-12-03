package es

import (
	"context"
	"github.com/jyk1987/es/tool"

	"github.com/jyk1987/es/data"
	"github.com/jyk1987/es/log"
	"github.com/jyk1987/es/node"
	"github.com/smallnest/rpcx/client"
)

var _RpcClientCache = make(map[string]client.XClient)

// getRpcClient 获取一个rpc客户端
func getRpcClient(nodeName string) (client.XClient, error) {
	tool.RLock(nodeName)
	c := _RpcClientCache[nodeName]
	if c != nil {
		tool.RUnlock(nodeName)
		return c, nil
	}
	tool.RUnlock(nodeName)
	tool.Lock(nodeName)
	consulServers := []string{node.GetNodeConfig().Consul}
	log.Log.Debug("consul server:", consulServers)
	d, e := client.NewConsulDiscovery(data.DiscoverBasePath, nodeName, consulServers, nil)
	if e != nil {
		log.Log.Error(e)
		return nil, e
	}
	c = client.NewXClient(nodeName, client.Failover, client.RoundRobin, d, client.DefaultOption)
	_RpcClientCache[nodeName] = c
	tool.Unlock(nodeName)
	return c, nil
}

func callServiceExecute(nodeName, path, method string, params ...interface{}) (*data.Result, error) {
	c, e := getRpcClient(nodeName)
	if e != nil {
		return nil, e
	}
	req := &data.Request{
		NodeName: nodeName,
		Path:     path,
		Method:   method,
	}
	req.SetParameters(params...)
	result := new(data.Result)
	tool.RLock(nodeName)
	defer tool.RUnlock(nodeName)
	e = c.Call(context.Background(), node.RpcExecuteFuncName, req, result)
	if e != nil {
		return nil, e
	}
	return result, nil
}

func callServiceGetInfo(nodeName string) (*data.Result, error) {
	c, e := getRpcClient(nodeName)
	if e != nil {
		return nil, e
	}
	result := new(data.Result)
	tool.RLock(nodeName)
	defer tool.RUnlock(nodeName)
	e = c.Call(context.Background(), node.RpcGetInfoFuncName, nil, result)
	if e != nil {
		return nil, e
	}
	return result, nil
}
