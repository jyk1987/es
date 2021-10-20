package node

import (
	"errors"
	"flag"
	"fmt"
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/log"
	"github.com/smallnest/rpcx/server"
	"reflect"
	"sync"

	"context"
)

// _Services 存放所有的服务
var _Services map[string]*Service

// _ServicesLock 所有服务的操作锁
var _ServicesLock sync.RWMutex

func init() {
	// 初始化存放所有服务的map
	_Services = make(map[string]*Service, 0)
}

// Service 服务,存储func
type Service struct {
	Path          string             //结构体的包路径
	Instance      interface{}        //实例的指针
	Methods       map[string]*Method //实例的所有方法
	MethodSetLock sync.RWMutex       // 服务设置锁
}

// GetMethod 获取方法
func (s *Service) GetMethod(methodName string) *Method {
	s.MethodSetLock.RLock()
	defer s.MethodSetLock.RUnlock()
	m, ok := s.Methods[methodName]
	if ok {
		return m
	}
	return nil
}

// Method 存储方法
type Method struct {
	MethodName     string         //方法名称
	MethodInstance reflect.Value  //方法实例
	MethodType     reflect.Type   //方法类型
	ParamCount     int            //参数个数
	ParamsType     []reflect.Type //全部参数的Type
	ReturnCount    int            //返回数据个数
	ReturnsType    []reflect.Type //全部返回数据的Type
}

// Execute 执行方法
// TODO：此方法需要进行性能分析，并对性能做出优化
func (m *Method) Execute(params [][]byte) (*data.Result, error) {
	paramsLen := len(params)
	if paramsLen != m.ParamCount {
		log.Log.Error("方法参数个数不相符:paramcount:", paramsLen, "method count:", m.ParamCount)
		return nil, errors.New("方法参数个数不相符！")
	}
	inputArgs := make([]reflect.Value, m.ParamCount)
	for i := 0; i < paramsLen; i++ {
		paramType := m.ParamsType[i]
		param, e := data.DecodeDataByType(params[i], paramType)
		if e != nil {
			return nil, e
		}
		if param == nil {
			inputArgs[i] = reflect.New(paramType).Elem()
		} else {
			inputArgs[i] = reflect.ValueOf(param)
		}
	}
	outs := m.MethodInstance.Call(inputArgs)

	r, e := data.NewResult(outs)
	if e != nil {
		return nil, e
	}
	return r, nil
}

// _NewService 创建一个服务
func _NewService(instRef interface{}) *Service {
	s := &Service{}
	instType := reflect.TypeOf(instRef) //获取实例的类型
	if instType.String()[0] != '*' {
		log.Log.Panic("请使用new方式创建服务，然后进行注册！", instType)
	}
	instValue := reflect.ValueOf(instRef) //获取实例值
	s.Path = instType.String()[1:]        //设置实例的包路径
	s.Instance = instValue                //设置实例引用
	log.Log.Info("注册服务:", s.Path)
	methodCount := instType.NumMethod() //获取方法总数
	methods := make(map[string]*Method, methodCount)
	for i := 0; i < methodCount; i++ {
		m := &Method{}                                          //创建方法结构
		m.MethodName = instType.Method(i).Name                  //方法名称
		m.MethodInstance = instValue.MethodByName(m.MethodName) //方法实例
		m.MethodType = m.MethodInstance.Type()                  //方法类型
		log.Log.Info("注册方法", i, ":", m.MethodName)
		//初始化方法的所有参数数据
		paramCount := m.MethodType.NumIn() //获取参数个数
		m.ParamCount = paramCount          //设置参数个数
		paramsType := make([]reflect.Type, paramCount)
		for j := 0; j < paramCount; j++ {
			paramsType[j] = m.MethodType.In(j) //获取参数的Type
			log.Log.Info("参数", j, m.MethodType.In(j))
		}
		m.ParamsType = paramsType
		//初始化方法的所有返回数据
		returnCount := m.MethodType.NumOut() //获取返回数据个数
		m.ReturnCount = returnCount          //设置返回数据个数
		returnsType := make([]reflect.Type, returnCount)
		for j := 0; j < returnCount; j++ {
			returnsType[j] = m.MethodType.Out(j) //设置每个返回参数的类型
		}
		m.ReturnsType = returnsType
		s.MethodSetLock.Lock()
		methods[m.MethodName] = m
		s.MethodSetLock.Unlock()
	}
	s.Methods = methods
	return s
}

// Reg 注册本地服务
func Reg(serviceInstance interface{}) {
	service := _NewService(serviceInstance)
	_ServicesLock.Lock()
	_Services[service.Path] = service
	_ServicesLock.Unlock()
	log.Log.Info("注册完毕:", service.Path)
}

// ExecuteService 执行（本地）服务
func ExecuteService(request *data.Request) (*data.Result, error) {
	// 获取服务
	_ServicesLock.RLock()
	s, ok := _Services[request.Path]
	_ServicesLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("服务没有找到,path:%v", request.Path)
	}
	// 获取方法
	s.MethodSetLock.RLock()
	m := s.GetMethod(request.Method)
	s.MethodSetLock.RUnlock()
	if m == nil {
		return nil, fmt.Errorf("方法没有找到,path:%v,method:%v", request.Path, request.Method)
	}
	// 执行方法
	return m.Execute(request.Parameters)
}

// ESNode rpcx暴露的服务
type ESNode struct{}

// Execute 执行服务
func (*ESNode) Execute(ctx context.Context, request *data.Request, result *data.Result) error {
	r, err := ExecuteService(request)
	if err != nil {
		return err
	}
	result.SetData(r.GetData())
	return nil
}

// InitNodeServer 初始化rpc服务端
func InitNodeServer() {
	addr := flag.String("addr", "0.0.0.0:3456", "server address")
	flag.Parse()
	s := server.NewServer()
	s.Register(new(ESNode), "")
	s.Serve("tcp", *addr)
}
