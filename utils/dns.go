package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/urlooker/web/g"
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

	ret, _, err := c.Exchange(m, "114.114.114.114:53")
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
	return ips, err
}

func InternalDns(domain string) []g.IpIdc {
	return []g.IpIdc{}
}
