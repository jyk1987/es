package main

import (
	"errors"
	"fmt"
	"gitee.com/jyk1987/es"
	"gitee.com/jyk1987/es/node"
)

type ServerDemo struct {
	Name  string
	Value *ServerDemo
}

func (*ServerDemo) Service1(a, b string) (string, *ServerDemo, error) {
	sd := &ServerDemo{Name: "张三", Value: &ServerDemo{Name: "张小三"}}
	println("input args:", a, b)
	return a + b, sd, errors.New("sdasdasd")
}

func init() {
	es.Reg(new(ServerDemo))
}
func main() {
	// 调用服务
	result, err := es.Call("", "main.ServerDemo", "Service1", "你好", "世界")
	// 判断调用过程中是否有出错
	if err != nil {
		println("调用服务方法出错", err)
		return
	}
	fmt.Println(result.Returns)
	//// 声明与服务方法相同的返回参数用于接受返回参数
	//var r string
	//var e error
	//re := result.GetResult(func(result string, err error) {
	//	r = result
	//	e = err
	//})
	//if re != nil {
	//	println(re.Error())
	//	return
	//}
	//// 方法执行完成后续操作
	//if e != nil {
	//	println(e.Error())
	//}
	//println(r)
	//println("服务测试完成")
	//println("开启启动node...")
	node.InitNodeServer()
}
