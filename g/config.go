package g

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type LogConfig struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Level    string `json:"level"`
}

type MysqlConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type HttpConfig struct {
	Listen string `json:"listen"`
	Secret string `json:"secret"`
}

type RpcConfig struct {
	Listen string `json:"listen"`
}

type AlarmConfig struct {
	Enable      bool              `json:"enable"`
	Batch       int               `json:"batch"`
	Replicas    int               `json:"replicas"`
	ConnTimeout int               `json:"connTimeout"`
	CallTimeout int               `json:"callTimeout"`
	MaxConns    int               `json:"maxConns"`
	MaxIdle     int               `json:"maxIdle"`
	SleepTime   int               `json:"sleepTime"`
	Cluster     map[string]string `json:"cluster"`
}

type FalconConfig struct {
	Enable   bool   `json:"enable"`
	Addr     string `json:"addr"`
	Interval int    `json:"interval"`
}

type InternalDnsConfig struct {
	Enable bool   `json:"enable"`
	Addr   string `json:"addr"`
}

type GlobalConfig struct {
	Debug       bool                `json:"debug"`
	Admins      []string            `json:"admins"`
	Salt        string              `json:"salt"`
	Past        int                 `json:"past"`
	Http        *HttpConfig         `json:"http"`
	Rpc         *RpcConfig          `json:"rpc"`
	Log         *LogConfig          `json:"log"`
	Mysql       *MysqlConfig        `json:"mysql"`
	Alarm       *AlarmConfig        `json:"alarm"`
	Falcon      *FalconConfig       `json:"falcon"`
	InternalDns *InternalDnsConfig  `json:"internalDns"`
	MonitorMap  map[string][]string `json:"monitorMap"`
}

var (
	Config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	Config = &c

	log.Println("load configuration file", cfg, "successfully")
	return nil
}
