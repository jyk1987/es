package main

import (
	"gitee.com/jyk1987/es"
	"gitee.com/jyk1987/es/log"
	"github.com/gogf/gf/os/gfile"
	"time"
)

type ServerDemo struct {
	Name  string
	Value *ServerDemo
}

func (*ServerDemo) Service1(a, b string, c []byte, sd *ServerDemo) (string, *ServerDemo, []int, error) {
	sd.Value = &ServerDemo{Name: sd.Name + "儿子"}
	//fmt.Println(c)
	//println("input args:", a, b, c, sd)
	return a + b, sd, []int{1, 2, 3}, nil
}

func (*ServerDemo) UploadFile(fileName string, data []byte) (string, error) {
	f, e := gfile.Create(fileName)
	if e != nil {
		return "", e
	}
	_, e = f.Write(data)
	if e != nil {
		return "", e
	}
	log.Log.Debug(data)
	f.Close()
	return f.Name(), nil
}

func init() {
	es.Reg(new(ServerDemo))
}
func main() {
	e := es.InitNode()
	if e != nil {
		log.Log.Error(e)
		return
	}
	go es.StartNode()
	for {
		sd := &ServerDemo{Name: "张三"}
		// 调用服务
		result, err := es.Call("node1", "main.ServerDemo", "Service1", "你好", "世界", nil, sd)
		// 判断调用过程中是否有出错
		if err != nil {
			log.Log.Error("调用服务方法出错", err)
			continue
		}
		e = result.GetResult(func(s string, sd *ServerDemo, is []int, e error) {
			//log.Log.Debug(s, sd, is, e)
		})
		if e != nil {
			log.Log.Error(e)
		}
		time.Sleep(time.Second * 3)
	}
	//c := []byte("sfasf")

}
