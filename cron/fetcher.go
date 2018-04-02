package cron

import (
	"log"
	"strings"
	"time"

	"github.com/urlooker/web/g"
	"github.com/urlooker/web/model"
	"github.com/urlooker/web/utils"
)

func GetDetectedItem() {
	t1 := time.NewTicker(time.Duration(60) * time.Second)
	for {
		err := getDetectedItem()
		if err != nil {
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
		log.Println(k)
		for _, i := range v {
			log.Println(i)
		}
	}

	g.DetectedItemMap.Set(detectedItemMap)
	return nil
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
