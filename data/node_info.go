package data

// NodeInfo 节点信息
type NodeInfo struct {
	NodeName  string                  //节点名称
	ESVersion int                     // 存储es的版本号
	Services  map[string]*ServiceInfo // 当前节点包含的服务
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Path    string                 //服务对象（结构）路径
	Methods map[string]*MethodInfo // 服务提供的方法
}

// MethodInfo 方法信息
type MethodInfo struct {
	MethodName  string   //方法名
	ParamsType  []string //每个参数类型
	ReturnsType []string // 方法返回参数的类型
}
