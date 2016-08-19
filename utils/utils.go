package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/toolkits/str"

	"github.com/urlooker/web/g"
)

func Getkey(idc string) string {
	keys := g.Config.MonitorMap[idc]
	count := len(keys)
	now := int(time.Now().Unix())
	return keys[now%count]
}

func IsIP(ip string) bool {
	if ip != "" {
		isOk, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, ip)
		if isOk {
			return isOk
		}
	}
	return false
}

func ParseUrl(target string) (schema, host, port, path string) {
	targetArr := strings.Split(target, "//")
	schema = targetArr[0]
	url := strings.Split(targetArr[1], "/")
	addrArr := strings.Split(url[0], ":")
	if len(addrArr) == 2 {
		host = addrArr[0]
		port = addrArr[1]
	} else {
		host = url[0]
	}

	for _, seg := range url[1:] {
		path += ("/" + seg)
	}
	return
}

func KeysOfMap(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for key, _ := range m {
		keys[i] = key
		i++
	}

	return keys
}

func EncryptPassword(raw string) string {
	return str.Md5Encode(g.Config.Salt + raw)
}

func CheckUrl(url string) error {
	if !strings.Contains(url, "https://") && !strings.Contains(url, "http://") {
		return fmt.Errorf("http or https is necessary")
	}
	if len(url) > 1024 {
		return fmt.Errorf("url is too long over 1024")
	}
	return nil
}

func TimeFormat(ts int64) string {
	t := time.Unix(ts, 0).Format("2006-01-02 15:04:05")
	arr := strings.Split(t, " ")
	t = arr[1]
	arr = strings.Split(t, ":")

	return fmt.Sprintf("%s:%s", arr[0], arr[1])
}
