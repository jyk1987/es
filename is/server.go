package is

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/log"
	"github.com/smallnest/rpcx/server"
	"time"
)

// Node 节点注册需要传入的信息
type Node struct {
	Services map[string]*ServiceInfo
	NodeInfo *NodeInfo
}

// Ping 节点ping数据包
type Ping struct {
	NodeName            string
	UUID                string
	LastActive          int64
	ServiceIndexVersion int64
}
type ReplyState int

const (
	ReplyOK ReplyState = iota
	ReplyFail
	ReplyServiceNofound
)

type Reply struct {
	State               ReplyState            //状态
	ServiceIndexVersion int64                 //服务信息版本
	ServiceIndex        map[string]*IndexInfo //所有服务的索引信息
	ESVersion           int                   //当前的版本
}

type ESIndexServer struct {
}

const RpcServiceName = "ESIndexServer"
const RpcRegNodeFuncName = "RegNode"
const RpcPingFuncName = "Ping"

func (is *ESIndexServer) RegNode(ctx context.Context, node *Node, reply *Reply) error {
	reply.ESVersion = data.ESVersion
	if node == nil {
		return errors.New("节点信息不能为空")
	}
	if len(node.NodeInfo.NodeName) == 0 {
		return errors.New("节点名不能为空")
	}
	if len(node.NodeInfo.UUID) == 0 {
		return errors.New("节点UUID不能为空")
	}
	if len(node.NodeInfo.IP) == 0 {
		return errors.New("节点ip不能为空")
	}
	if node.NodeInfo.Port == 0 {
		return errors.New("节点port不能为0")
	}
	//if node.Services == nil || len(node.Services) == 0 {
	//	log.Log.Error("节点服务不能为空", node)
	//	return errors.New("节点服务不能为空")
	//}
	regNode(node)
	reply.State = ReplyOK
	reply.ServiceIndex, reply.ServiceIndexVersion = getServiceIndex(0)
	return nil
}

func (is *ESIndexServer) Ping(ctx context.Context, ping *Ping, reply *Reply) error {
	reply.ESVersion = data.ESVersion
	if ping.LastActive == 0 {
		ping.LastActive = time.Now().UnixMilli()
	}
	if len(ping.NodeName) == 0 {
		return errors.New("节点名不能为空")
	}
	if len(ping.UUID) == 0 {
		return errors.New("节点UUID不能为空")
	}
	if !active(ping) {
		reply.State = ReplyServiceNofound
		return nil
	}
	reply.State = ReplyOK
	reply.ServiceIndex, reply.ServiceIndexVersion = getServiceIndex(ping.ServiceIndexVersion)
	return nil
}

func InitESIndexServer() error {
	cfg, err := data.GetConfig()
	if err != nil {
		log.Log.Error("加载配置文件出错:", err)
		return err
	}
	addr := flag.String("addr", fmt.Sprintf("0.0.0.0:%v", cfg.Port), "server address")
	flag.Parse()
	s := server.NewServer()
	s.Register(new(ESIndexServer), "")
	s.Serve("tcp", *addr)
	return nil
}
