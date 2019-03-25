package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/peng19940915/urlooker/web/g"
)

func LookupIP(domain string, timeout int) ([]string, error) {
	var c dns.Client
	var err error
	ips := []string{}
	domain = strings.TrimRight(domain, ".") + "."
	c.DialTimeout = time.Duration(timeout) * time.Millisecond
	c.ReadTimeout = time.Duration(timeout) * time.Millisecond
	c.WriteTimeout = time.Duration(timeout) * time.Millisecond

	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeA)

	ret, _, err := c.Exchange(m, g.Config.DnsServer)
	if err != nil {
		domain = strings.TrimRight(domain, ".")
		e := fmt.Sprintf("lookup error: %s, %s", domain, err.Error())
		return ips, errors.New(e)
	}

	for _, i := range ret.Answer {
		result := strings.Split(i.String(), "\t")
		if result[3] == "A" && IsIP(result[4]) {
			ips = append(ips, result[4])
		}
	}

	log.Println("ips:", ips)
	return ips, err
}

func InternalDns(domain string) []g.IpIdc {
	return []g.IpIdc{}
}

// 获取主机所在的idc
func InternalHostIdc(hostname string) string{
	return "default"
}