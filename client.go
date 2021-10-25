package es

import (
	"context"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/util/guid"
	"github.com/jyk1987/es/data"
	"github.com/jyk1987/es/log"
	"github.com/jyk1987/es/node"
	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
)

var _UUID = guid.S()

var _RpcClientCache = make(map[string]client.XClient)

// getRpcClient 获取一个rpc客户端
func getRpcClient(nodeName string) (client.XClient, error) {
	gmlock.RLock(nodeName)
	c := _RpcClientCache[nodeName]
	if c != nil {
		gmlock.RUnlock(nodeName)
		return c, nil
	}
	gmlock.RUnlock(nodeName)
	gmlock.Lock(nodeName)
	etcurl := []string{node.GetNodeConfig().Etcd}
	log.Log.Debug(etcurl)
	d, e := etcd_client.NewEtcdV3Discovery(data.ETCDBasePath, nodeName, etcurl, false, nil)
	if e != nil {
		log.Log.Error(e)
		return nil, e
	}
	c = client.NewXClient(nodeName, client.Failover, client.RoundRobin, d, client.DefaultOption)
	_RpcClientCache[nodeName] = c
	gmlock.Unlock(nodeName)
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
	gmlock.RLock(nodeName)
	defer gmlock.RUnlock(nodeName)
	e = c.Call(context.Background(), node.RpcExecuteFuncName, req, result)
	if e != nil {
		return nil, e
	}
	return result, nil
}
