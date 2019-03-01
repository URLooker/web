package utils

import (
	"fmt"
	"net/http"

	"github.com/peng19940915/urlooker/web/g"
)

var transport *http.Transport = &http.Transport{}


func PushFalcon(itemCheckedArray []*g.CheckResult, hostname string) {
	pushDatas := make([]*g.MetricValue, 0)
	for _, itemChecked := range itemCheckedArray {
		tags := fmt.Sprintf("ip=%s,domain=%s,creator=%s,from=%s", itemChecked.Ip, itemChecked.Domain, itemChecked.Creator, hostname)
		if len(itemChecked.Tag) > 0 { //补充用户自定义tag
			tags += "," + itemChecked.Tag
		}

		//url 状态
		data := getMetric(itemChecked, "url_status", tags, itemChecked.Status)
		pushDatas = append(pushDatas, &data)

		//url 响应时间
		data = getMetric(itemChecked, "url_resp_time", tags, int64(itemChecked.RespTime))
		pushDatas = append(pushDatas, &data)
	}
	SendMetrics(pushDatas)

}

func getMetric(item *g.CheckResult, metric, tags string, value int64) g.MetricValue {
	var data g.MetricValue
	data.Endpoint = fmt.Sprintf("url_%d", item.Sid)
	data.Timestamp = item.PushTime
	data.Metric = metric
	data.Type = "GAUGE"
	data.Step = int64(g.Config.Falcon.Interval)
	data.Tags = tags
	data.Value = value
	return data
}

// 推送port扫描数据到falcon
func PushPort2Falcon(itemCheckedArray []*g.CheckPortResult, hostname string) {
	pushDatas := make([]*g.MetricValue, 0)
	for _, itemChecked := range itemCheckedArray {
		tags := fmt.Sprintf("ip=%s,host=%s,creator=%s,from=%s", itemChecked.Ip, itemChecked.Host, itemChecked.Creator, hostname)
		if len(itemChecked.Tag) > 0 { //补充用户自定义tag
			tags += "," + itemChecked.Tag
		}

		//url 状态
		data := getPortMetric(itemChecked, "port_status", tags, int64(itemChecked.Result))
		pushDatas = append(pushDatas, data)
		//url 响应时间
		data = getPortMetric(itemChecked, "port_resp_time", tags, int64(itemChecked.RespTime))
		pushDatas = append(pushDatas, data)
	}

	for _, d := range pushDatas {
		fmt.Println(d)
	}
	SendMetrics(pushDatas)
}

func getPortMetric(item *g.CheckPortResult, metric, tags string, value int64) *g.MetricValue {
	var data g.MetricValue
	data.Endpoint = fmt.Sprintf("port_%d", item.Sid)
	data.Timestamp = item.PushTime
	data.Type = "GAUGE"
	data.Metric = metric
	data.Step = int64(g.Config.Falcon.Interval)
	data.Tags = tags
	data.Value = value
	return &data
}


