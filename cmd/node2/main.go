package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/jyk1987/es"
	"github.com/jyk1987/es/log"
	"io/ioutil"
	"sync"
	"time"
)

func main() {

	e := es.InitES()
	if e != nil {
		log.Log.Error(e)
		return
	}
	//go es.StartNode()
	testGetInfo()
	//test()
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

func testGetInfo() {
	r, e := es.GetNodeInfo("node1")
	if e != nil {
		log.Log.Error(e)
	}
	log.Log.Debug(g.Export(r))
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
	es.Call("node1", "nana.ServerDemo", "Service1", args...)
	begin := time.Now()
	count := 10 //10000 * 10
	var execCount int
	tcount := 1
	log.Log.Info("开始测试", count)
	wg := sync.WaitGroup{}
	for i := 0; i < tcount; i++ {
		wg.Add(1)
		go func() {
			for {
				if execCount++; execCount > count {
					break
				}
				result, e := es.Call("node1", "nana.ServerDemo", "Service1", args...)
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
