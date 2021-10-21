package main

import (
	"gitee.com/jyk1987/es"
	"gitee.com/jyk1987/es/log"
	"github.com/gogf/gf/os/gfile"
	"io/ioutil"
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
	testUpload()
	//log.Log.Info("平均耗时：", end/time.Duration(count))
}

func testUpload() {
	f, e := gfile.Open("README.md")
	if e != nil {
		log.Log.Error(e)
		return
	}
	data, e := ioutil.ReadAll(f)
	if e != nil {
		log.Log.Error(e)
		return
	}
	log.Log.Debug(len(data))
	r, e := es.Call("node1", "main.ServerDemo", "UploadFile", "r.md", data)

	if e != nil {
		log.Log.Error(e)
		return
	}
	log.Log.Debug(r.GetData())
}

func test() {
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
	count := 10000 * 100
	execCount := 0
	tcount := 1000
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
}
