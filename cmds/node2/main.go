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
	}
	reqs.SetParameters(args...)

	begin := time.Now()
	count := 10000 * 10
	log.Log.Info("开始测试", count)
	for i := 0; i < count; i++ {
		result := new(data.Result)
		err = xclient.Call(context.Background(), "Execute", reqs, result)
		if err != nil {
			fmt.Println("调用失败:", err)
		}
	}
	end := time.Now()
	log.Log.Info("测试结束，总耗时：", end.Sub(begin))
	//log.Log.Info("平均耗时：", end/time.Duration(count))

	//for i := 0; i < len(result.Returns); i++ {
	//	r := result.Returns[i]
	//	fmt.Println(string(r.Binary))
	//}

}
