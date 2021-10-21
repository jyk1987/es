package main

import (
	"gitee.com/jyk1987/es"
	"gitee.com/jyk1987/es/log"
	"sync"
	"time"
)

func main() {
	e := es.InitNode()
	if e != nil {
		log.Log.Error(e)
		return
	}
	go es.StartNode()

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

	begin := time.Now()
	count := 10000 * 10
	execCount := 0
	tcount := 100
	log.Log.Info("开始测试", count)
	wg := sync.WaitGroup{}
	for i := 0; i < tcount; i++ {
		wg.Add(1)
		go func() {
			for {
				if execCount++; execCount > count {
					break
				}
				result, e := es.Call("node1", "main.ServerDemo", "Service1", args...)
				if e != nil {
					log.Log.Debug(e)
				}
				e = result.GetResult(func(s string, sd *ServerDemo, is []int, e error) {
					//log.Log.Debug(s, sd, is, e)
				})
				if e != nil {
					log.Log.Error(e)
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
