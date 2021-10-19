package main

import (
	"context"
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
	//share.Codecs[protocol.SerializeType(4)] = &data.GobCodec{}
	option := client.DefaultOption
	//option.SerializeType = protocol.SerializeType(4)
	xclient := client.NewXClient("ESNode", client.Failtry, client.RandomSelect, d, option)

	defer xclient.Close()
	//c := map[string]string{"saf": "safas"}
	args := make([]interface{}, 3)
	args[0] = "你好"
	args[1] = "再见"
	args[2] = make([]byte, 0)
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
		fmt.Println("调用失败:", err)
	}
	for i := 0; i < len(result.Returns); i++ {
		r := result.Returns[i]
		fmt.Println(string(r.Binary))
	}
	//r := ""
	//e := errors.New("")
	//err = result.GetResult(func(str string, er error) {
	//	r = str
	//	e = er
	//})
	//if err != nil {
	//	fmt.Printf("failed to convert result: %v", err)
	//}
	//if e != nil {
	//	println(e.Error())
	//}
	//fmt.Println(r)
}
