package tool

import (
	"fmt"
	"net"
	"strings"
)

// GetOutBoundIP 获取本地出口网卡的本地ip地址
func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "114.114.114.114:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
