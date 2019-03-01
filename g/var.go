package g

import (
	"sync"
	log "github.com/sirupsen/logrus"
	"time"
)

type IpIdc struct {
	Ip  string
	Idc string
}

//下发给agent的数据结构
type DetectedItem struct {
	Sid        int64  `json:"sid"`
	Domain     string `json:"domain"`
	Target     string `json:"target"`
	Ip         string `json:"ip"`
	Keywords   string `json:"keywords"`
	Timeout    int    `json:"timeout"`
	Creator    string `json:"creator"`
	Data       string `json:"data"`
	Tag        string `json:"tag"`
	ExpectCode string `json:"expect_code"`
	Idc        string `json:"idc"`
}

type DetectedPortItem struct {
	Sid        int64  `json:"sid"`
	Host       string `json:"target"`
	Port       string `json:"port"`
	Ip         string `json:"ip"`
	Keywords   string `json:"keywords"`
	Timeout    int    `json:"timeout"`
	Creator    string `json:"creator"`
	Tag        string `json:"tag"`
	Idc        string `json:"idc"`
}

//agent上报的数据结构
type CheckResult struct {
	Sid      int64  `json:"sid"`
	Domain   string `json:"domain"`
	Target   string `json:"target"`
	Creator  string `json:"creator"`
	Tag      string `json:"tag"`
	RespCode string `json:"resp_code"`
	RespTime int    `json:"resp_time"`
	Status   int64  `json:"status"`
	PushTime int64  `json:"push_time"`
	Ip       string `json:"ip"`
}

type CheckPortResult struct {
	Sid      int64  `json:"sid"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Creator  string `json:"creator"`
	Tag      string `json:"tag"`
	RespTime int    `json:"resp_time"`
	Result   int64  `json:"result"`
	Status   int64  `json:"status"`
	PushTime int64  `json:"push_time"`
	Ip       string `json:"ip"`
}
type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Tags      string      `json:"tags"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
	Type      string      `json:"counterType"`
	Step      int64       `json:"step"`
}

type DetectedItemSafeMap struct {
	sync.RWMutex
	M map[string][]*DetectedItem
}

type DetectedPortItemSafeMap struct {
	sync.RWMutex
	M map[string][]*DetectedPortItem
}
var (
	DetectedItemMap = &DetectedItemSafeMap{M: make(map[string][]*DetectedItem)}
	DetectedPortItemMap = &DetectedPortItemSafeMap{M: make(map[string][] *DetectedPortItem)}
)

func (this *DetectedItemSafeMap) Get(key string) ([]*DetectedItem, bool) {
	this.RLock()
	defer this.RUnlock()
	ipItem, exists := this.M[key]
	return ipItem, exists
}

func (this *DetectedItemSafeMap) Set(detectedItemMap map[string][]*DetectedItem) {
	this.Lock()
	defer this.Unlock()
	this.M = detectedItemMap
}

// 端口扫描方法
func (this *DetectedPortItemSafeMap) Get(key string) ([]*DetectedPortItem, bool) {
	this.RLock()
	defer this.RUnlock()
	ipItem, exists := this.M[key]
	return ipItem, exists
}

func (this *DetectedPortItemSafeMap) Set(detectedPortItemMap map[string][]*DetectedPortItem) {
	this.Lock()
	defer this.Unlock()
	this.M = detectedPortItemMap
}

func InitLog() {
	logPath := Config.Log.LogPath
	level:= Config.Log.LogLevel
	logFilename := Config.Log.FileName
	maxAge := Config.Log.MaxAge * time.Hour * 24
	rotationTime := Config.Log.RotationTime * time.Hour * 24
	// 配置日志
	ConfigLocalFilesystemLogger(logPath, logFilename, maxAge, rotationTime)
	// 配置级别
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.Fatal("log conf only allow [info, debug, warn], please check your confguire")
	}
}
