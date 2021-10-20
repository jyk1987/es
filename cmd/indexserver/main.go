package main

import (
	"gitee.com/jyk1987/es/is"
	"gitee.com/jyk1987/es/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
	log.Log.Info("索引服务启动...")
	go func() {
		err := is.InitESIndexServer()
		if err != nil {
			log.Log.Panic(err)
		}
	}()
	<-done
	// TODO: 触发关闭代码，完成已经注册数据的序列化，方便下次启动加载。
	log.Log.Info("索引服务关闭完成")
}
