package main

import "gitee.com/jyk1987/es"

type ServerDemo struct {
}

func (*ServerDemo) Service1(a, b string) (string, error) {
	println("input args:", a, b)
	return "ok", nil
}

func init() {
	es.Reg(new(ServerDemo))
}
func main() {

}
