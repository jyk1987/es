package data

// IndexInfo 服务器索引信息，需要持久化
type IndexInfo struct {
	NodeName     string                  //节点名称
	Online       bool                    //服务是否在线
	ServicesCode uint32                  // 当前存储的Services的hashCode，用于判断服务是否发生变化
	Services     map[string]*ServiceInfo // 当前节点包含的服务
	Nodes        map[string]*NodeInfo    // 当前节点ming在线的实例
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
