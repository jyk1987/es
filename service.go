package es

import (
	"gitee.com/jyk1987/es/data"
	server "gitee.com/jyk1987/es/server"
)

// Reg 注册本地服务
func Reg(serviceInstance interface{}) {
	server.Reg(serviceInstance)
}

func Call(node, path, method string, params ...interface{}) (*data.Result, error) {
	request := data.Request{Node: node, Path: path, Method: method, Args: params}
	return server.ExecuteService(request)
}

//func Call(path, methodName string, params ...interface{}) ([]reflect.Value, error, bool) {
//// 获取服务
//_ServicesLock.RLock()
//s, ok := _Services[path]
//_ServicesLock.RUnlock()
////未找到服务
//if !ok {
//	msg := "没有找到服务:" + path
//	log.Error(msg)
//	return nil, errors.New(msg), false
//}
////服务不包含方法
//if s.Methods == nil {
//	msg := "服务" + path + "的方法为初始化"
//	log.Error(msg)
//	return nil, errors.New(msg), false
//}
////获取需要调用的方法
//m, ok := s.Methods[methodName]
//if !ok {
//	msg := "服务" + path + "中没有找到" + methodName + "方法"
//	log.Error(msg)
//	return nil, errors.New(msg), false
//}
//if m.ParamCount != len(params) {
//	msg := "参数不匹配:" + path + "." + methodName + "参数为" + strconv.Itoa(m.ParamCount) +
//		"个,传入参数为" + strconv.Itoa(len(params)) + "个"
//	log.Error(msg)
//	return nil, errors.New(msg), false
//}
//in := make([]reflect.Value, m.ParamCount) //创建参数集
//for i := 0; i < m.ParamCount; i++ {
//	paramType := m.ParamsType[i]        //获取单个参数的类型
//	param := reflect.ValueOf(params[i]) //实例化传入的参数位reflect.Value类型
//	// TODO:实现不同类型的转换，防止报错
//	in[i] = param.Convert(paramType) //转换参数为目标类型,增加兼容性
//}
//defer func() {
//	err := recover()
//	if err != nil {
//		log.Error(path, ",", methodName, "调用失败")
//		log.Error(err)
//	}
//}()
//result := m.MethodInstance.Call(in)
//log.Debug("方法执行成功：", path, methodName, params)
//log.Debug("结果:")
//for _, r := range result {
//	log.Debug(r.String())
//}
//return result, nil, true
//}

//func Spread(serviceName string, path string, methodName string, params ...interface{}) error {
//	if len(params) > 0 {
//		cbf := params[len(params)-1]
//		cbfType := reflect.TypeOf(cbf)
//		//reflect.Method{cbf}
//		//最后一个参数是func
//		if cbfType.Kind().String() == "func" {
//			params = params[0 : len(params)-1]
//		} else {
//			cbf = nil
//		}
//		outParams := make([]reflect.Value, cbfType.NumIn())
//		result, e, ok := CallService(path, methodName, params...)
//		if !ok {
//			return errors.New("远程服务执行失败")
//		}
//		if e != nil {
//			return e
//		}
//		if len(result) > 1 && !result[1].IsNil() {
//			return result[1].Interface().(error)
//		}
//		if cbf != nil && len(result) > 0 {
//			outParams[0] = result[0]
//			executeResult := reflect.ValueOf(cbf).Call(outParams)
//			log.Debug(executeResult)
//		}
//		return nil
//	}
//	_, e, ok := CallService(path, methodName, params...)
//	if !ok {
//		return errors.New("远程服务执行失败")
//	}
//	if e != nil {
//		return e
//	}
//	return nil
//
//}
