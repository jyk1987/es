package es

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/is"
	"gitee.com/jyk1987/es/log"
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

var _ServiceIndexCache = make(map[string]*is.IndexInfo, 0)
var _ServiceIndexVersion int64

var _ClientCache = make(map[string]*_ServiceClient)

var _IndexServerClient *_ServiceClient

type _ServiceClient struct {
	NodeName string                  // 节点名
	Type     ServiceClientType       // 客户端类型
	Info     map[string]*is.NodeInfo //服务节点的信息
	xclient  client.XClient          //通讯使用的client
}

func callServiceExecute(nodeName, path, method string, params ...interface{}) (*data.Result, error) {
	sc := getServiceClient(nodeName)
	if sc == nil {
		return nil, errors.New("未找到服务:" + nodeName)
	}
	c, e := sc.getClient(node.RpcServiceName)
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
	gmlock.RLock(sc.NodeName)
	defer gmlock.RUnlock(sc.NodeName)
	e = c.Call(context.Background(), node.RpcExecuteFuncName, req, result)
	if e != nil {
		return nil, e
	}
	return result, nil
}

func callServerRegNode() error {
	if _IndexServerClient == nil {
		_InitIndexServerClient()
	}
	c, e := _IndexServerClient.getClient(is.RpcServiceName)
	if e != nil {
		return e
	}
	localIndex := node.GetLocalServiceIndex()
	ip, e := tool.GetOutBoundIP()
	if e != nil {
		return e
	}
	node := &is.Node{
		Services: localIndex,
		NodeInfo: &is.NodeInfo{
			NodeName:   node.GetNodeConfig().Name,
			UUID:       _UUID,
			LastActive: time.Now().UnixMilli(),
			IP:         ip,
			Port:       node.GetNodeConfig().Port,
			ESVersion:  data.ESVersion,
		},
	}
	reply := new(is.Reply)
	gmlock.RLock(IndexServerNodeName)
	e = c.Call(context.Background(), is.RpcRegNodeFuncName, node, reply)
	gmlock.RUnlock(IndexServerNodeName)
	if e != nil {
		return e
	}
	_setServiceIndex(reply.ServiceIndex, reply.ServiceIndexVersion)
	return nil
}

func callServerPing() error {
	if _IndexServerClient == nil {
		_InitIndexServerClient()
	}
	c, e := _IndexServerClient.getClient(is.RpcServiceName)
	if e != nil {
		return e
	}
	ping := &is.Ping{
		NodeName:            node.GetNodeConfig().Name,
		UUID:                _UUID,
		LastActive:          time.Now().UnixMilli(),
		ServiceIndexVersion: _ServiceIndexVersion,
	}
	reply := new(is.Reply)
	gmlock.RLock(IndexServerNodeName)
	e = c.Call(context.Background(), is.RpcPingFuncName, ping, reply)
	gmlock.RUnlock(IndexServerNodeName)
	if reply.State == is.ReplyServiceNofound {
		go callServerRegNode()
	}
	if e != nil {
		return e
	}
	_setServiceIndex(reply.ServiceIndex, reply.ServiceIndexVersion)
	return nil
}

func _setServiceIndex(index map[string]*is.IndexInfo, version int64) {
	if len(index) == 0 && version == 0 {
		return
	}

	for nodeName, _ := range index {
		gmlock.Lock(nodeName)
		s := _ServiceIndexCache[nodeName]
		if s == nil {
			_ServiceIndexCache[nodeName] = index[nodeName]
		}
		cc := _ClientCache[nodeName]
		if cc != nil && cc.xclient != nil {
			//cc.xclient.Close()
			cc.xclient = nil
		}
		_ServiceIndexVersion = version
		gmlock.Unlock(nodeName)
	}
}

func (sc *_ServiceClient) getClient(rpcServiceName string) (client.XClient, error) {
	if sc.xclient != nil {
		return sc.xclient, nil
	}
	if len(sc.Info) == 0 {
		return nil, fmt.Errorf("service:%v,没有在线的节点", sc.NodeName)
	}
	gmlock.Lock(sc.NodeName)
	defer gmlock.Unlock(sc.NodeName)
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
		sc.xclient = client.NewXClient(rpcServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
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
		sc.xclient = client.NewXClient(rpcServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	}
	return sc.xclient, nil
}

const IndexServerNodeName = "ESIS"

func _InitIndexServerClient() {
	endpoint := node.GetNodeConfig().Server
	if len(endpoint) == 0 {
		log.Log.Panic("为配置索引服务地址:", node.GetNodeConfig())
	}
	addr := strings.Split(endpoint, ":")
	ip := addr[0]
	log.Log.Debug(addr)
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
		gmlock.RUnlock(nodeName)
		return c
	}
	gmlock.RUnlock(nodeName)
	indexInfo := _ServiceIndexCache[nodeName]
	if indexInfo == nil {
		return nil
	}
	gmlock.Lock(nodeName)
	defer gmlock.Unlock(nodeName)
	c := _NewServiceClient(indexInfo.NodeName, indexInfo.Nodes)
	_ClientCache[nodeName] = c
	return c
}

func _NewServiceClient(nodeName string, nodeInfo map[string]*is.NodeInfo) *_ServiceClient {
	sc := &_ServiceClient{
		NodeName: nodeName,
		Info:     nodeInfo,
	}
	return sc
}
