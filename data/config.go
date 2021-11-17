package data

import (
	"errors"
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/jyk1987/es/log"

	_ "github.com/jinzhu/configor"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const ESVersion = 70
const DefaultPort = 8910
const ESConfigPath = "esconfig"

var DefaultConfigs = []string{"es.yml", "es.json"}

const ETCDBasePath = "/es_rpc"

// ESConfig 配置文件映射结构
type ESConfig struct {
	Port     int    `default:"8910"` //服务端口,默认端口8910
	Name     string //系统中的nodename用于区分不同服务
	Key      string //链接密钥，用于链接到整个系统中
	Etcd     string //发现服务地址
	Endpoint string //访问端点，如果配置，服务启动时会使用访问端点向etcd进行注册，其他服务会通过此访问端点来访问此服务
}

// GetAppPath 获取程序运行路径
func GetAppPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Log.Panic("获取启动目录失败:", err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// GetPwdPath 获取启动指令的执行目录
func GetPwdPath() string {
	path, _ := os.Getwd()
	return path
}

var _Configs = make(map[string]*ESConfig)
var _ConfigsLock sync.RWMutex

func GetConfig(configFile ...string) (*ESConfig, error) {
	config := &ESConfig{Port: DefaultPort}
	configNames := make([]string, 0)
	var fileName string
	if len(configFile) > 0 {
		fileName = configFile[0]
		configNames = append(configNames, configFile[0])
	}
	_ConfigsLock.RLock()
	if c := _Configs[fileName]; c != nil {
		_ConfigsLock.RUnlock()
		return c, nil
	}
	_ConfigsLock.RUnlock()
	for i := 0; i < len(DefaultConfigs); i++ {
		configNames = append(configNames, DefaultConfigs[i])
	}
	log.Log.Debug(configNames)
	var fullPath string
	for _, name := range configNames {
		fullPath, _ = SearchFile(filepath.Join(ESConfigPath, name))
		if len(fullPath) > 0 {
			break
		}
	}
	if len(fullPath) == 0 {
		return nil, errors.New("config file not fount")
	}

	log.Log.Infof("load config file:%v", fullPath)
	err := configor.Load(config, fullPath)
	if err != nil {
		return nil, err
	}
	log.Log.Debug(config)
	//json := jsoniter.ConfigCompatibleWithStandardLibrary
	//e = json.Unmarshal(b, config)
	//if e != nil {
	//	return nil, e
	//}
	_ConfigsLock.Lock()
	_Configs[fileName] = config
	_ConfigsLock.Unlock()
	return config, nil
}

func SearchFile(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	fullPath := filepath.Join(GetAppPath(), path)
	if FileExist(fullPath) {
		return fullPath, nil
	}
	fullPath = filepath.Join(GetPwdPath(), path)
	if FileExist(fullPath) {
		return fullPath, nil
	}
	return "", fmt.Errorf("file %v not fount", path)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

const ESKeyFileExt = ".eskey"

// ESKey 密钥文件映射结构，索引服务配置文件目录需要增加相应客户端的密码，客户端才可以连
type ESKey struct {
	Name string `json:"name"` // NodeName
	Key  string `json:"key"`  // 密钥
}
