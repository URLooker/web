package test

import (
	"testing"
	"github.com/toolkits/net"
	"time"
	"fmt"
	"github.com/urlooker/web/api"
)

func Test_RRpc(t * testing.T) {
	client, err := net.JsonRpcClient("tcp", "127.0.0.1:1985", time.Second)
	if err != nil {
		fmt.Println("cannot connect to")
	}
	var reply interface{}
	var resp api.GetPortItemResponse
	client.Call("Web.GetPortItem", "hostname.1", &resp)
	client.Close()
	fmt.Println("reply: ", reply)
}
