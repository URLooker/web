package test

import (
	"testing"
	"crypto/tls"
	"time"
	"log"
	"github.com/astaxie/beego/httplib"
	"fmt"
)

func Test_302(t *testing.T) {
	req := httplib.Get("http://falcon.gridsum.com")
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(3*time.Second, 10*time.Second)
	req.Header("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.SetHost("falcon.gridsum.com")
	resp, err := req.Response()

	if err != nil {
		log.Println("[ERROR]:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)
}
