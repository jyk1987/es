package is

import (
	"gitee.com/jyk1987/es/log"
	"gitee.com/jyk1987/es/tool"
	"github.com/gogf/gf/encoding/ghash"
	"sync"
	"time"
)

const _NodeActiveCheckInterval = time.Second * 1
const _NodeActiveOverTime = time.Second * 5

var _ServiceIndexVersion int64 // 所有服务状态的版本，服务索引服务发生变化后刷新此版本
func RefreshServiceIndexVersion() {
	_ServiceIndexVersion = time.Now().UnixMilli()
}

var _ServiceIndex map[string]*IndexInfo
var _ServiceIndexLock sync.RWMutex

func init() {
	_ServiceIndex = make(map[string]*IndexInfo, 0)
}

// IndexInfo 服务器索引信息，需要持久化
type IndexInfo struct {
	NodeName     string                  //节点名称
	Online       bool                    //服务是否在线
	ServicesCode uint32                  // 当前存储的Services的hashCode，用于判断服务是否发生变化
	Services     map[string]*ServiceInfo // 当前节点包含的服务
	nodes        map[string]*NodeInfo    // 当前节点ming在线的实例
}

func regNode(node *Node) error {
	d, e := tool.EncodeData(node.Services)
	if e != nil {
		log.Log.Error("节点注册失败:", node)
		return e
	}
	if node.NodeInfo.LastActive == 0 {
		node.NodeInfo.LastActive = time.Now().UnixMilli()
	}
	code := ghash.BKDRHash(d)
	_ServiceIndexLock.Lock()
	defer _ServiceIndexLock.Unlock()
	s := _ServiceIndex[node.NodeInfo.NodeName]
	if s == nil {
		indexInfo := &IndexInfo{
			NodeName:     node.NodeInfo.NodeName,
			Online:       true,
			ServicesCode: code,
			Services:     node.Services,
			nodes:        map[string]*NodeInfo{node.NodeInfo.UUID: node.NodeInfo},
		}
		_ServiceIndex[node.NodeInfo.NodeName] = indexInfo
		go node.NodeInfo.activeCheck()
		RefreshServiceIndexVersion()
	} else {
		if s.ServicesCode != code {
			s.Services = node.Services
			RefreshServiceIndexVersion()
		}
		nd := s.nodes[node.NodeInfo.UUID]
		if nd == nil {
			s.nodes[node.NodeInfo.UUID] = node.NodeInfo
			go node.NodeInfo.activeCheck()
			RefreshServiceIndexVersion()
		} else {
			nd.LastActive = node.NodeInfo.LastActive
		}
	}
	return nil
}

func active(ping *Ping) {
	_ServiceIndexLock.RLock()
	defer _ServiceIndexLock.RUnlock()
	s := _ServiceIndex[ping.NodeName]
	if s == nil {
		log.Log.Warning("需要刷新的服务信息不存在:", ping)
		return
	}
	n := s.nodes[ping.UUID]
	if n == nil {
		log.Log.Warning("需要刷新的节点信息不存在:", ping)
		return
	}
	n.LastActive = ping.LastActive
}

// NodeInfo 节点信息
type NodeInfo struct {
	NodeName   string //节点名
	UUID       string //节点的唯一标识
	LastActive int64  //节点最后活跃时间,毫秒时间戳
	IP         string //节点ip
	Port       int    //服务端口
	ESVersion  int    //服务使用的es版本
}

func (node *NodeInfo) activeCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Log.Error("节点活跃检测协程出错", err)
			go node.activeCheck()
		}
	}()
	for {
		time.Sleep(_NodeActiveCheckInterval)
		if node.LastActive == -1 {
			break
		}
		lastActive := time.UnixMilli(node.LastActive)
		if time.Now().Sub(lastActive) > _NodeActiveOverTime {
			log.Log.Infof("节点:%v,UUID:%v,已超时为活动,上次活跃时间:%v", node.NodeName, node.UUID, lastActive)
			removeNode(node)
			break
		}
	}
}
func removeNode(node *NodeInfo) {
	_ServiceIndexLock.Lock()
	defer _ServiceIndexLock.Unlock()
	s := _ServiceIndex[node.NodeName]
	if s == nil {
		return
	}
	n := s.nodes[node.UUID]
	if n != nil {
		n.LastActive = -1
	}
	delete(s.nodes, node.UUID)
	log.Log.Infof("节点:%v,UUID:%v,被移除.", node.NodeName, node.UUID)
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Path    string                 //服务对象（结构）路径
	Methods map[string]*MethodInfo // 服务提供的方法
}

// MethodInfo 方法信息
type MethodInfo struct {
	MethodName  string   //方法名
	ParamCount  int      //方法参数个数
	ParamsType  []string //每个参数类型
	ReturnCount int      //方法返回参数和数
	ReturnsType []string // 方法返回参数的类型
}
