package node

import (
	"bytes"
	"errors"
	"github.com/jyk1987/es/data"
	"github.com/jyk1987/es/log"
	"github.com/jyk1987/es/tool"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// _Services 存放所有本地服务
var _Services map[string]*_Service

// _ServicesLock 所有本地服务的操作锁
var _ServicesLock sync.RWMutex

func init() {
	// 初始化存放所有服务的map
	_Services = make(map[string]*_Service, 0)
}

// Config 节点配置文件，在初始化后存储
var _Config *data.ESConfig

func GetNodeConfig() *data.ESConfig {
	if _Config == nil {
		cfg, e := data.GetConfig()
		if e != nil {
			log.Log.Error(e)
			return nil
		}
		_Config = cfg
	}
	return _Config
}

// _Service 服务,存储func
type _Service struct {
	Path        string              //结构体的包路径
	instance    interface{}         //实例的指针
	methods     map[string]*_Method //实例的所有方法
	methodsLock sync.RWMutex        // 服务设置锁
}

// GetMethod 获取方法
func (s *_Service) GetMethod(methodName string) *_Method {
	s.methodsLock.RLock()
	defer s.methodsLock.RUnlock()
	m, ok := s.methods[methodName]
	if ok {
		return m
	}
	return nil
}

// _Method 存储方法
type _Method struct {
	methodName  string         //方法名称
	instance    reflect.Value  //方法实例
	methodType  reflect.Type   //方法类型
	paramCount  int            //参数个数
	paramsType  []reflect.Type //全部参数的Type
	returnCount int            //返回数据个数
	returnsType []reflect.Type //全部返回数据的Type
}

// Execute 执行相应的请求
func (m *_Method) Execute(request *data.Request) (*data.Result, error) {
	// 获取参数的编码数据
	params := request.Parameters
	paramsLen := len(params)
	// 判断传入参数与方法接收参数数量是否一致
	if paramsLen != m.paramCount {
		log.Log.Error("方法参数个数不相符:paramcount:", paramsLen, "method count:", m.paramCount)
		return nil, errors.New("方法参数个数不相符！")
	}
	// 创建最终参数集用于执行调用方法
	inputArgs := make([]reflect.Value, m.paramCount)
	for i := 0; i < paramsLen; i++ {
		paramType := m.paramsType[i]
		// 解码参数为相应方法参数type的interface
		param, e := tool.DecodeDataByType(params[i], paramType)
		if e != nil {
			return nil, e
		}
		// 如果参数为nil，需要反射创建一个空的reflect.Value,防止反射执行方法报错
		if param == nil {
			inputArgs[i] = reflect.New(paramType).Elem()
		} else {
			inputArgs[i] = reflect.ValueOf(param)
		}
	}
	// 执行方法
	outs := m.instance.Call(inputArgs)
	// 将方法执行结果[]reflect.Value转换为es.data.Result结构用于最终的数据返回
	// 构建过程中会将reflect.Value类型通过编码器最终编码为[]byte数据
	r, e := data.NewResult(outs)
	if e != nil {
		return nil, e
	}
	return r, nil
}

// _NewService 创建一个服务
func _NewService(PkgPath string, instRef interface{}) *_Service {
	if len(PkgPath) == 0 {
		log.Log.Error("PkgPath is empty string")
		return nil
	}
	if instRef == nil {
		log.Log.Error("instRef is nil interface")
		return nil
	}
	s := &_Service{}
	instType := reflect.TypeOf(instRef) //获取实例的类型
	if instType.String()[0] != '*' {
		log.Log.Panic("请使用new方式创建服务，然后进行注册！", instType)
	}
	instValue := reflect.ValueOf(instRef) //获取实例值
	structName := strings.Split(instType.String(), ".")[1]
	s.Path = PkgPath + "." + structName //设置实例的包路径
	s.instance = instValue              //设置实例引用
	log.Log.Info("register:", s.Path)
	methodCount := instType.NumMethod() //获取方法总数
	methods := make(map[string]*_Method, methodCount)
	for i := 0; i < methodCount; i++ {
		m := &_Method{}                                   //创建方法结构
		m.methodName = instType.Method(i).Name            //方法名称
		m.instance = instValue.MethodByName(m.methodName) //方法实例
		m.methodType = m.instance.Type()                  //方法类型
		log.Log.Info("method", i, ":", m.methodName)
		//初始化方法的所有参数数据
		paramCount := m.methodType.NumIn() //获取参数个数
		m.paramCount = paramCount          //设置参数个数
		paramsType := make([]reflect.Type, paramCount)
		parameterContent := bytes.Buffer{}
		for j := 0; j < paramCount; j++ {
			paramsType[j] = m.methodType.In(j) //获取参数的Type
			parameterContent.WriteString("p")
			parameterContent.WriteString(strconv.Itoa(j))
			parameterContent.WriteString(":")
			parameterContent.WriteString(m.methodType.In(j).String())
			parameterContent.WriteString("\t")
		}
		log.Log.Info(parameterContent.String())
		m.paramsType = paramsType
		//初始化方法的所有返回数据
		returnCount := m.methodType.NumOut() //获取返回数据个数
		m.returnCount = returnCount          //设置返回数据个数
		returnsType := make([]reflect.Type, returnCount)
		for j := 0; j < returnCount; j++ {
			returnsType[j] = m.methodType.Out(j) //设置每个返回参数的类型
		}
		m.returnsType = returnsType
		s.methodsLock.Lock()
		methods[m.methodName] = m
		s.methodsLock.Unlock()
	}
	s.methods = methods
	return s
}

// Reg 注册本地服务
func Reg(servicePath string, serviceInstance interface{}) {
	service := _NewService(servicePath, serviceInstance)
	_ServicesLock.Lock()
	_Services[service.Path] = service
	_ServicesLock.Unlock()
	log.Log.Info("complete.")
}

// getService 获取一个服务
func getService(path string) *_Service {
	_ServicesLock.RLock()
	defer _ServicesLock.RUnlock()
	s, ok := _Services[path]
	if ok {
		return s
	}
	return nil
}

// GetLocalServiceIndex 获取本地服务器的索引信息
func GetLocalServiceIndex() map[string]*data.ServiceInfo {
	_ServicesLock.Lock()
	defer _ServicesLock.Unlock()
	localIndex := make(map[string]*data.ServiceInfo, len(_Services))
	for path, s := range _Services {
		serviceInfo := &data.ServiceInfo{
			Path:    path,
			Methods: make(map[string]*data.MethodInfo, len(s.methods)),
		}
		for name, m := range s.methods {
			methodInfo := &data.MethodInfo{
				MethodName:  m.methodName,
				ParamCount:  m.paramCount,
				ParamsType:  make([]string, m.paramCount),
				ReturnCount: m.returnCount,
				ReturnsType: make([]string, m.returnCount),
			}
			for i := 0; i < m.paramCount; i++ {
				methodInfo.ParamsType[i] = m.paramsType[i].String()
			}
			for i := 0; i < m.returnCount; i++ {
				methodInfo.ReturnsType[i] = m.returnsType[i].String()
			}
			serviceInfo.Methods[name] = methodInfo
		}
		localIndex[path] = serviceInfo
	}
	return localIndex
}
