package es

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/is"
	"gitee.com/jyk1987/es/node"
	"gitee.com/jyk1987/es/tool"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/util/guid"
	"github.com/smallnest/rpcx/client"
	"strconv"
	"strings"
	"time"
)

var _UUID = guid.S()

type ServiceClientType int

const (
	ServiceClientIS = iota
	ServiceClientNode
)

var _ServiceIndexCache map[string]*is.IndexInfo
var _ServiceIndexVersion int64

var _ClientCache = make(map[string]*_ServiceClient)

var _IndexServerClient *_ServiceClient

type _ServiceClient struct {
	NodeName string                  // 节点名
	Type     ServiceClientType       // 客户端类型
	Info     map[string]*is.NodeInfo //服务节点的信息
}

func callServiceExecute(nodeName, path, method string, params ...interface{}) (*data.Result, error) {
	sc := getServiceClient(nodeName)
	if sc == nil {
		return nil, errors.New("未找到服务:" + nodeName)
	}
	c, e := sc._CreateClient(node.RpcServiceName)
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
	e = c.Call(context.Background(), node.RpcExecuteFuncName, req, result)
	if e != nil {
		return nil, e
	}
	defer c.Close()
	return result, nil
}

func callServerRegNode() error {
	if _IndexServerClient == nil {
		_InitIndexServerClient()
	}
	c, e := _IndexServerClient._CreateClient(is.RpcServiceName)
	if e != nil {
		return e
	}
	defer c.Close()
	localIndex := node.GetLocalServiceIndex()
	ip, e := tool.GetOutBoundIP()
	if e != nil {
		return e
	}
	node := &is.Node{
		Services: localIndex,
		NodeInfo: &is.NodeInfo{
			NodeName:   node.Config.Name,
			UUID:       _UUID,
			LastActive: time.Now().UnixMilli(),
			IP:         ip,
			Port:       node.Config.Port,
			ESVersion:  data.ESVersion,
		},
	}
	reply := new(is.Reply)
	e = c.Call(context.Background(), is.RpcServiceName, node, reply)
	if e != nil {
		return e
	}
	_ServiceIndexCache = reply.ServiceIndex
	_ServiceIndexVersion = reply.ServiceIndexVersion
	return nil
}

func callServerPing() error {
	if _IndexServerClient == nil {
		_InitIndexServerClient()
	}
	c, e := _IndexServerClient._CreateClient(is.RpcServiceName)
	if e != nil {
		return e
	}
	defer c.Close()
	ping := &is.Ping{
		NodeName:            node.Config.Name,
		UUID:                _UUID,
		LastActive:          time.Now().UnixMilli(),
		ServiceIndexVersion: _ServiceIndexVersion,
	}
	reply := new(is.Reply)
	e = c.Call(context.Background(), is.RpcServiceName, ping, reply)
	if e != nil {
		return e
	}
	if len(reply.ServiceIndex) > 0 && reply.ServiceIndexVersion > 0 {
		_ServiceIndexCache = reply.ServiceIndex
		_ServiceIndexVersion = reply.ServiceIndexVersion
	}
	return nil
}

func (sc *_ServiceClient) _CreateClient(rpcServiceName string) (client.XClient, error) {
	if len(sc.Info) == 0 {
		return nil, fmt.Errorf("service:%v,没有在线的节点", sc.NodeName)
	}
	if len(sc.Info) == 1 {
		var node *is.NodeInfo
		for _, info := range sc.Info {
			node = info
			break
		}
		d, e := client.NewPeer2PeerDiscovery(fmt.Sprintf("tcp@%v:%v", node.IP, node.Port), "")
		if e != nil {
			return nil, e
		}
		c := client.NewXClient(rpcServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		return c, nil
	} else {
		address := make([]*client.KVPair, len(sc.Info))
		index := 0
		for _, node := range sc.Info {
			address[index] = &client.KVPair{Key: fmt.Sprintf("tcp@%v:%v", node.IP, node.Port)}
			index++
		}
		d, e := client.NewMultipleServersDiscovery(address)
		if e != nil {
			return nil, e
		}
		c := client.NewXClient(rpcServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		return c, nil
	}
}

const IndexServerNodeName = "ESIS"

func _InitIndexServerClient() {
	endpoint := node.Config.Server
	addr := strings.Split(endpoint, ":")
	ip := addr[0]
	port, _ := strconv.Atoi(addr[1])
	_IndexServerClient = _NewServiceClient(IndexServerNodeName, map[string]*is.NodeInfo{
		IndexServerNodeName: {
			NodeName: IndexServerNodeName,
			IP:       ip,
			Port:     port,
		},
	})
}

func getServiceClient(nodeName string) *_ServiceClient {
	gmlock.RLock(nodeName)
	if c := _ClientCache[nodeName]; c != nil {
		gmlock.Unlock(nodeName)
		return c
	}
	indexInfo := _ServiceIndexCache[nodeName]
	if indexInfo == nil {
		return nil
	}
	gmlock.Lock(nodeName)
	c := _NewServiceClient(indexInfo.NodeName, indexInfo.Nodes)
	_ClientCache[nodeName] = c
	gmlock.Unlock(nodeName)
	return c
}

func _NewServiceClient(nodeName string, nodeInfo map[string]*is.NodeInfo) *_ServiceClient {
	sc := &_ServiceClient{
		NodeName: nodeName,
		Info:     nodeInfo,
	}
	return sc
}
