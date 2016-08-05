package g

import (
	"sync"
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

type DetectedItemSafeMap struct {
	sync.RWMutex
	M map[string][]*DetectedItem
}

var (
	DetectedItemMap = &DetectedItemSafeMap{M: make(map[string][]*DetectedItem)}
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
