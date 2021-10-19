package main

import (
	"context"
	"flag"
	"fmt"
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/log"
	"github.com/smallnest/rpcx/client"
	"time"
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
	//c := map[string]string{"saf": "safas"}
	args := make([]interface{}, 4)
	args[0] = "你好"
	args[1] = "再见"
	args[2] = make([]byte, 0)
	type ServerDemo struct {
		Name  string
		Value *ServerDemo
	}
	sd := &ServerDemo{Name: "李四"}
	args[3] = sd
	reqs := &data.Request{
		NodeName: "",
		Path:     "main.ServerDemo",
		Method:   "Service1",
	}

	reqs.SetParameters(args...)

	begin := time.Now()
	count := 1 //10000 * 100
	log.Log.Info("开始测试", count)
	for i := 0; i < count; i++ {
		result := new(data.Result)
		err = xclient.Call(context.Background(), "Execute", reqs, result)
		if err != nil {
			fmt.Println("调用失败:", err)
		}
		for i := 0; i < len(result.Returns); i++ {
			r := result.Returns[i]
			fmt.Println(string(r.Binary))
		}
	}
	end := time.Now()
	log.Log.Info("测试结束，总耗时：", end.Sub(begin))
	//log.Log.Info("平均耗时：", end/time.Duration(count))

}
