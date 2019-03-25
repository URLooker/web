package g

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/toolkits/file"
	"time"
)



type MysqlConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type Log struct {
	LogLevel     string                  `json:"logLevel"`
	RotationTime time.Duration           `json:"rotationTime"`
	LogPath      string                  `json:"logPath"`
	MaxAge       time.Duration           `json:"maxAge"`
	FileName     string                  `json:"fileName"`
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
	Addrs    []string `json:"addrs"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
}


type InternalDnsConfig struct {
	Enable bool   `json:"enable"`
	Addr   string `json:"addr"`
}
type SSOConfig struct {
	ServerUrl  string `json:"serverUrl"`
	ServiceUrl string `json:"serviceUrl"`
}
type GlobalConfig struct {
	Log         *Log                `json:"log"`
	Admins      []string            `json:"admins"`
	Salt        string              `json:"salt"`
	Past        int                 `json:"past"` //查看最近几分钟内的报警历史和绘图，默认为30分钟
	Http        *HttpConfig         `json:"http"`
	Rpc         *RpcConfig          `json:"rpc"`
	SSO         *SSOConfig          `json:"sso"`
	Mysql       *MysqlConfig        `json:"mysql"`
	Alarm       *AlarmConfig        `json:"alarm"`
	DnsServer   string              `json:"dnsServer"`
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
