package main

import (
	"github.com/jyk1987/es"
	"github.com/jyk1987/es/log"
	"time"
)

type ServerDemo struct {
	Name  string
	Value *ServerDemo
	//CallCount int64
}

var CallCount int

func (s *ServerDemo) Service1(a, b string, c []byte, sd *ServerDemo) (string, *ServerDemo, []int, error) {
	sd.Value = &ServerDemo{Name: sd.Name + "儿子"}
	CallCount++
	//fmt.Println(CallCount)
	//println("input args:", a, b, c, sd)
	return a + b, sd, []int{1, 2, 3}, nil
}

func init() {
	es.Reg("nana", new(ServerDemo))
}
func main() {
	e := es.InitES()
	if e != nil {
		log.Log.Error(e)
		return
	}
	go es.StartNode()
	for {
		time.Sleep(time.Second * 1)

	}

}
