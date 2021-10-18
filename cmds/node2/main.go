package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"gitee.com/jyk1987/es/data"
	"github.com/smallnest/rpcx/client"
)

func main() {
	addr := flag.String("addr", "localhost:3456", "server address")
	flag.Parse()
	d, err := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	if err != nil {
		println(err.Error())
		return
	}
	xclient := client.NewXClient("ESNode", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := make([]interface{}, 2)
	args[0] = "你好"
	args[1] = "再见"
	reqs := &data.Request{
		NodeName: "",
		Path:     "main.ServerDemo",
		Method:   "Service1",
		Args:     args,
	}
	result := new(data.Result)
	println("开始调用")
	err = xclient.Call(context.Background(), "Execute", reqs, result)
	println("完成调用")
	if err != nil {
		fmt.Printf("failed to call: %v", err)
	}
	for i := 0; i < len(result.Returns); i++ {

	}
	r := ""
	e := errors.New("")
	err = result.GetResult(func(str string, er error) {
		r = str
		e = er
	})
	if err != nil {
		fmt.Printf("failed to convert result: %v", err)
	}
	if e != nil {
		println(e.Error())
	}
	fmt.Println(r)
}
