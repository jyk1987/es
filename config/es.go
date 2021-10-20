package config

import (
	"os"
	"path/filepath"
	"strings"

	"gitee.com/jyk1987/es/log"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gfile"
)

const DefaultPort = 8910
const ESConfigPath = "esconfig"
const ESConfigFileName = "es.json"

// ESConfig 配置文件映射结构
type ESConfig struct {
	Port int    `json:"port"` //服务端口,默认端口8910
	Name string `json:"name"` //系统中的nodename用于区分不同服务
	Key  string `json:"key"`  //链接密钥，用于链接到整个系统中
	//以下问索引服务配置
	// 访问端点指的是一个可以公共可访问的端点，ip+端口 或者域名+端口
	//IndexServer才需要访问端点，比如istest.kuaibang360.com:3456
	Endpoint string `json:"endpoint"` // 访问端点
}

// GetCurrentDirectory 获取程序运行路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Log.Panic("获取启动目录失败:", err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// GetRunDirectory 获取启动指令的执行目录
func GetRunDirectory() string {
	path, _ := os.Getwd()
	return path
}

// func GetConfigPath() string {
// 	basePath := GetRunDirectory()
// 	findPath := func(p string) string {
// 		//TODO:需要添加目录搜索功能
// 		configPath := filepath.Join(p, ESConfigPath)
// 		return configPath
// 	}
// 	return findPath(basePath)
// }

func GetConfig(configFileName ...string) (*ESConfig, error) {
	config := &ESConfig{Port: DefaultPort}
	fileName := ESConfigFileName
	if len(configFileName) > 0 {
		fileName = configFileName[0]
	}
	fullPath := gfile.Join(ESConfigPath, fileName)
	json, err := gjson.Load(fullPath)
	if err != nil {
		log.Log.Errorf("加载配置文件%v出错:%v", fullPath, err)
		return nil, err
	}
	err = json.Scan(config)
	if err != nil {
		log.Log.Errorf("转换配置文件%v出错:%v", fullPath, err)
		return nil, err
	}
	return config, nil
}

const ESKeyFileExt = ".eskey"

// ESKey 密钥文件映射结构，索引服务配置文件目录需要增加相应客户端的密码，客户端才可以连
type ESKey struct {
	Name string `json:"name"` // NodeName
	Key  string `json:"key"`  // 密钥
}
