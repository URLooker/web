package cron

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/model"
	"github.com/peng19940915/urlooker/web/utils"
)

func GetDetectedItem() {
	t1 := time.NewTicker(time.Duration(60) * time.Second)
	for {
		err1 := getDetectedItem()
		err2 := getPortDetectedItem()
		if err1 != nil || err2  != nil{
			time.Sleep(time.Second * 1)
			continue
		}
		<-t1.C
	}
}

func getDetectedItem() error {
	detectedItemMap := make(map[string][]*g.DetectedItem)
	stras, err := model.GetAllStrategyByCron()
	if err != nil {
		log.Println("get strategies error:", err)
		return err
	}

	for _, s := range stras {
		_, domain, _, _ := utils.ParseUrl(s.Url)
		var ipIdcArr []g.IpIdc
		if s.IP != "" {
			ips := strings.Split(s.IP, ",")
			for _, ip := range ips {
				var tmp g.IpIdc
				tmp.Ip = ip
				tmp.Idc = "default"
				ipIdcArr = append(ipIdcArr, tmp)
			}
		} else {
			ipIdcArr = getIpAndIdc(domain)
		}

		for _, tmp := range ipIdcArr {
			detectedItem := newDetectedItem(s, tmp.Ip, tmp.Idc)
			key := utils.Getkey(tmp.Idc, int(detectedItem.Sid))

			if _, exists := detectedItemMap[key]; exists {
				detectedItemMap[key] = append(detectedItemMap[key], &detectedItem)
			} else {
				detectedItemMap[key] = []*g.DetectedItem{&detectedItem}
			}
		}
	}

	for k, v := range detectedItemMap {
		log.Info(k)
		for _, i := range v {
			log.Info(i)
		}
	}
	g.DetectedItemMap.Set(detectedItemMap)
	return nil
}


func newDetectedItem(s *model.Strategy, ip string, idc string) g.DetectedItem {
	detectedItem := g.DetectedItem{
		Ip:         ip,
		Idc:        idc,
		Creator:    s.Creator,
		Sid:        s.Id,
		Keywords:   s.Keywords,
		Data:       s.Data,
		Tag:        s.Tag,
		ExpectCode: s.ExpectCode,
		Timeout:    s.Timeout,
	}
	schema, domain, port, path := utils.ParseUrl(s.Url)

	if port == "" {
		detectedItem.Target = schema + "//" + ip + path
	} else {
		detectedItem.Target = schema + "//" + ip + ":" + port + path
	}

	detectedItem.Domain = domain
	return detectedItem
}

func getIpAndIdc(domain string) []g.IpIdc {

	//公司内部提供接口，拿到域名解析的ip和机房列表
	if g.Config.InternalDns.Enable {
		ipIdcArr := utils.InternalDns(domain)
		return ipIdcArr
	}

	ipIdcArr := make([]g.IpIdc, 0)

	if utils.IsIP(domain) {
		var tmp g.IpIdc
		tmp.Ip = domain
		tmp.Idc = "default"
		ipIdcArr = append(ipIdcArr, tmp)
	} else {
		ips, _ := utils.LookupIP(domain, 5000)
		for _, ip := range ips {
			var tmp g.IpIdc
			tmp.Ip = ip
			tmp.Idc = "default"
			ipIdcArr = append(ipIdcArr, tmp)
		}
	}

	return ipIdcArr
}

// 获取Port 扫描数据
func getPortDetectedItem() error {
	detectedPortItemMap := make(map[string][]*g.DetectedPortItem)
	stras, err := model.GetAllPortStrategyByCron()
	if err != nil {
		log.Errorf("get port strategies error: %v", err)
		return err
	}
	for _, s := range stras {

		idc := utils.InternalHostIdc(s.Host)

		detectedPortItem := newPortDetectedItem(s, s.IP, idc)
		key := utils.Getkey(idc, int(detectedPortItem.Sid))

		if _, exists := detectedPortItemMap[key]; exists {
			detectedPortItemMap[key] = append(detectedPortItemMap[key], &detectedPortItem)
		} else {
			detectedPortItemMap[key] = []*g.DetectedPortItem{&detectedPortItem}
		}
	}
	for k, v := range detectedPortItemMap {
		log.Info(k)
		for _, i := range v {
			log.Info(i)
		}
	}

	g.DetectedPortItemMap.Set(detectedPortItemMap)
	return nil
}

func newPortDetectedItem(s *model.PortStrategy, ip string, idc string) g.DetectedPortItem {
	detectedPortItem := g.DetectedPortItem{
		Ip:         ip,
		Idc:        idc,
		Creator:    s.Creator,
		Sid:        s.Id,
		Keywords:   s.Keywords,
		Host:       s.Host,
		Tag:        s.Tag,
		Port:       s.Port,
		Timeout:    s.Timeout,
	}

	return detectedPortItem
}
