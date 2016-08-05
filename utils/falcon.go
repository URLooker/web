package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/astaxie/beego/httplib"

	"github.com/urlooker/web/g"
)

var transport *http.Transport = &http.Transport{}

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Tags      string      `json:"tags"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
	Type      string      `json:"counterType"`
	Step      int64       `json:"step"`
}

func PushFalcon(itemCheckedArray []*g.CheckResult, hostname string) {

	pushDatas := make([]*MetricValue, 0)
	for _, itemChecked := range itemCheckedArray {
		var data MetricValue
		data.Metric = "url_status"
		data.Endpoint = fmt.Sprintf("url_%d", itemChecked.Sid)
		data.Timestamp = itemChecked.PushTime
		data.Type = "GAUGE"
		data.Step = int64(g.Config.Falcon.Interval)
		data.Value = itemChecked.Status

		if len(itemChecked.Tag) < 1 {
			data.Tags = fmt.Sprintf("ip=%s,domain=%s,creator=%s,from=%s", itemChecked.Ip, itemChecked.Domain, itemChecked.Creator, hostname)
		} else {
			data.Tags = fmt.Sprintf("ip=%s,domain=%s,creator=%s,%s,from=%s", itemChecked.Ip, itemChecked.Domain, itemChecked.Creator, itemChecked.Tag, hostname)
		}

		pushDatas = append(pushDatas, &data)
	}

	err := push(pushDatas)
	if err != nil {
		log.Println("push error", err)
	}
}

func push(data []*MetricValue) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = httplib.Post(g.Config.Falcon.Addr).Body(d).String()
	if err != nil {
		return err
	}

	return nil
}
