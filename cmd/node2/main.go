package main

import (
	"context"
	"flag"
	"fmt"
	"gitee.com/jyk1987/es/data"
	"gitee.com/jyk1987/es/log"
	"github.com/smallnest/rpcx/client"
	"sync"
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
	count := 10000 * 100
	execCount := 0
	tcount := 100000
	log.Log.Info("开始测试", count)
	wg := sync.WaitGroup{}
	for i := 0; i < tcount; i++ {
		wg.Add(1)
		go func() {
			for {
				if execCount++; execCount > count {
					break
				}
				result := new(data.Result)
				err = xclient.Call(context.Background(), "Execute", reqs, result)
				if err != nil {
					fmt.Println("调用失败:", err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	end := time.Now()
	log.Log.Info(tcount, "线程测试", count/10000, "万次，测试结束，总耗时：", end.Sub(begin))
	//log.Log.Info("平均耗时：", end/time.Duration(count))

}
