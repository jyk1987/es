package es

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
)

func init() {
	// 初始化存放所有服务的map
	_Services = make(map[string]*Service, 0)
}

// _Services 存放所有的服务
var _Services map[string]*Service

// _ServicesLock 所有服务的操作锁
var _ServicesLock sync.RWMutex

// Service 服务,存储func
type Service struct {
	Path     string             //结构体的包路径
	Instance interface{}        //实例的指针
	Methods  map[string]*Method //实例的所有方法
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

// _NewService 创建一个服务
func _NewService(instRef interface{}) *Service {
	s := &Service{}
	instType := reflect.TypeOf(instRef) //获取实例的类型
	if instType.String()[0] != '*' {
		instType = reflect.TypeOf(&instRef)
	}
	instValue := reflect.ValueOf(instRef) //获取实例值
	// TODO: 包路径的不确定性处理
	s.Path = instType.String()[1:] //设置实例的包路径
	s.Instance = instValue         //设置实例引用
	eslog.Info("注册服务:", s.Path)
	methodCount := instType.NumMethod() //获取方法总数
	methods := make(map[string]*Method, methodCount)
	for i := 0; i < methodCount; i++ {
		m := &Method{}                                          //创建方法结构
		m.MethodName = instType.Method(i).Name                  //方法名称
		m.MethodInstance = instValue.MethodByName(m.MethodName) //方法实例
		m.MethodType = m.MethodInstance.Type()                  //方法类型
		eslog.Info("注册方法", i, ":", m.MethodName)
		//初始化方法的所有参数数据
		paramCount := m.MethodType.NumIn() //获取参数个数
		m.ParamCount = paramCount          //设置参数个数
		paramsType := make([]reflect.Type, paramCount)
		for j := 0; j < paramCount; j++ {
			paramsType[j] = m.MethodType.In(j) //获取参数的Type
			eslog.Info("参数", j, m.MethodType.In(j))
		}
		m.ParamsType = paramsType
		//初始化方法的所有返回数据
		returnCount := m.MethodType.NumOut() //获取返回数据个数
		m.ReturnCount = returnCount          //设置返回数据个数
		returnsType := make([]reflect.Type, returnCount)
		for j := 0; j < returnCount; j++ {
			returnsType[j] = m.MethodType.Out(j)
		}
		m.ReturnsType = returnsType
		methods[m.MethodName] = m
	}
	s.Methods = methods
	return s
}

// RegService 注册本地服务
func RegService(serviceInstance interface{}) {
	service := _NewService(serviceInstance)
	_ServicesLock.Lock()
	_Services[service.Path] = service
	_ServicesLock.Unlock()
	eslog.Info("注册完毕:", service.Path)
}

func CallService(path, methodName string, params ...interface{}) ([]reflect.Value, error) {
	// 获取服务
	_ServicesLock.RLock()
	s, ok := _Services[path]
	_ServicesLock.RUnlock()
	//未找到服务
	if !ok {
		msg := "没有找到服务:" + path
		eslog.Error(msg)
		return nil, errors.New(msg)
	}
	//服务不包含方法
	if s.Methods == nil {
		msg := "服务" + path + "的方法为初始化"
		eslog.Error(msg)
		return nil, errors.New(msg)
	}
	//获取需要调用的方法
	m, ok := s.Methods[methodName]
	if !ok {
		msg := "服务" + path + "中没有找到" + methodName + "方法"
		eslog.Error(msg)
		return nil, errors.New(msg)
	}
	if m.ParamCount != len(params) {
		msg := "参数不匹配:" + path + "." + methodName + "参数为" + strconv.Itoa(m.ParamCount) +
			"个,传入参数为" + strconv.Itoa(len(params)) + "个"
		eslog.Error(msg)
		return nil, errors.New(msg)
	}
	in := make([]reflect.Value, m.ParamCount) //创建参数集
	for i := 0; i < m.ParamCount; i++ {
		paramType := m.ParamsType[i]        //获取单个参数的类型
		param := reflect.ValueOf(params[i]) //实例化传入的参数位reflect.Value类型
		in[i] = param.Convert(paramType)    //转换参数为目标类型,增加兼容性
	}
	result := m.MethodInstance.Call(in)

	return result, nil
}
