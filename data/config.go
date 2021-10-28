package data

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/jyk1987/es/log"
)

const ESVersion = 56
const DefaultPort = 8910
const ESConfigPath = "esconfig"
const ESConfigFileName = "es.json"
const ETCDBasePath = "/es_rpc"

// ESConfig 配置文件映射结构
type ESConfig struct {
	Port     int    `json:"port"`     //服务端口,默认端口8910
	Name     string `json:"name"`     //系统中的nodename用于区分不同服务
	Key      string `json:"key"`      //链接密钥，用于链接到整个系统中
	Etcd     string `json:"etcd"`     //发现服务地址
	Endpoint string `json:"endpoint"` //访问端点，如果配置，服务启动时会使用访问端点向etcd进行注册，其他服务会通过此访问端点来访问此服务
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

var _Configs = make(map[string]*ESConfig)
var _ConfigsLock sync.RWMutex

func GetConfig(configFile ...string) (*ESConfig, error) {
	config := &ESConfig{Port: DefaultPort}
	fileName := ESConfigFileName
	if len(configFile) > 0 {
		fileName = configFile[0]
	}
	_ConfigsLock.RLock()
	if c := _Configs[fileName]; c != nil {
		_ConfigsLock.RUnlock()
		return c, nil
	}
	_ConfigsLock.RUnlock()
	fullPath := filepath.Join(ESConfigPath, fileName)
	b, e := os.ReadFile(fullPath)
	if e != nil {
		return nil, e
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	e = json.Unmarshal(b, config)
	if e != nil {
		return nil, e
	}
	_ConfigsLock.Lock()
	_Configs[fileName] = config
	_ConfigsLock.Unlock()
	return config, nil
}

const ESKeyFileExt = ".eskey"

// ESKey 密钥文件映射结构，索引服务配置文件目录需要增加相应客户端的密码，客户端才可以连
type ESKey struct {
	Name string `json:"name"` // NodeName
	Key  string `json:"key"`  // 密钥
}
