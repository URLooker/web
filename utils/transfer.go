package utils

import (
	"sync"
	"math/rand"
	"time"
	"github.com/peng19940915/urlooker/web/g"
	"math"
	"net/rpc"
	"github.com/toolkits/net"
	"log"
	"fmt"
)
type SingleConnRpcClient struct {
	sync.Mutex
	rpcClient *rpc.Client
	RpcServer string
	Timeout   time.Duration
}

var (
	TransferClientsLock *sync.RWMutex                   = new(sync.RWMutex)
	TransferClients     map[string]*SingleConnRpcClient = map[string]*SingleConnRpcClient{}
)

func (this *SingleConnRpcClient) close() {
	if this.rpcClient != nil {
		this.rpcClient.Close()
		this.rpcClient = nil
	}

}

func (this *SingleConnRpcClient) insureConn() (err error) {
	if this.rpcClient != nil {
		return
	}

	var retry int = 1
	for {
		if this.rpcClient != nil {
			return
		}

		this.rpcClient, err = net.JsonRpcClient("tcp", this.RpcServer, this.Timeout)
		if err == nil {
			return
		}

		log.Printf("ERROR: dial %s fail: %v", this.RpcServer, err)

		if retry > 3 {

			return fmt.Errorf("ERROR: init RpcClient Failed.")
		}

		time.Sleep(time.Duration(math.Pow(2.0, float64(retry))) * time.Second)

		retry++
	}
	return
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

	this.Lock()
	defer this.Unlock()

	err := this.insureConn()
	if err != nil {
		return err
	}
	timeout := time.Duration(50 * time.Second)
	done := make(chan error)

	go func() {
		err := this.rpcClient.Call(method, args, reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		log.Printf("WARN: rpc call timeout %v => %v", this.rpcClient, this.RpcServer)
		this.close()
	case err := <-done:
		if err != nil {
			this.close()
			return err
		}
	}

	return nil
}

func SendMetrics(metrics []*g.MetricValue) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(g.Config.Falcon.Addrs)) {
		addr := g.Config.Falcon.Addrs[i]
		log.Println("transfer is: %s", addr)
		c := getTransferClient(addr)
		var err error
		if c == nil {
			c, err = initTransferClient(addr)
			// 如果初始化失败，继续获取下一轮随机数
			if err != nil{
				fmt.Println("init rpc client failed, detail: %v", err.Error())
				//log.Errorf("init rpc client failed, detail: %v", err.Error())
				continue
			}
		}
		if updateMetrics(c, metrics) {
			break
		}
	}
}

func initTransferClient(addr string) (*SingleConnRpcClient, error) {
	var c *SingleConnRpcClient = &SingleConnRpcClient{
		RpcServer: addr,
		Timeout:   time.Duration(g.Config.Falcon.Timeout) * time.Millisecond,
	}

	TransferClientsLock.Lock()
	defer TransferClientsLock.Unlock()
	err := c.insureConn()
	if err == nil {
		TransferClients[addr] = c
		return c, nil
	}else {
		return nil, err
	}
}

func updateMetrics(c *SingleConnRpcClient, metrics []*g.MetricValue) bool {
	mvs := metrics
	//分批次传给transfer
	n := 700
	lenMvs := len(mvs)

	div := lenMvs / n
	mod := math.Mod(float64(lenMvs), float64(n))
	mvsSend := [] *g.MetricValue{}
	for i := 1; i <= div+1; i++ {
		if i < div+1 {
			mvsSend = mvs[n*(i-1) : n*i]
		} else {
			mvsSend = mvs[n*(i-1) : (n*(i-1))+int(mod)]
		}
		go callTransferUpdater(c, mvsSend)
		time.Sleep(100 * time.Millisecond)
	}
	return true
}


func getTransferClient(addr string) *SingleConnRpcClient {
	TransferClientsLock.RLock()
	defer TransferClientsLock.RUnlock()
	if c, ok := TransferClients[addr]; ok {
		var resp interface{}
		// 发送前先check下该连接是否存活。如果未存活则跳过，重新获取，同时删除该client
		err := c.rpcClient.Call("Transfer.Ping", nil, resp)
		if err != nil {
			fmt.Println("transfer: %s %s", addr, "is not alive")
			//log.Errorf("transfer: %s %s", addr, "is not alive")
			c.close()
			delete(TransferClients, addr)
			return nil
		}
		return c
	}
	return nil
}

func callTransferUpdater(c *SingleConnRpcClient, metrics []*g.MetricValue){
	//startTime := time.Now()

	//var resp model.TransferResponse
	var resp interface{}
	err := c.Call("Transfer.Update", metrics, resp)
	if err != nil {
		fmt.Println("call Transfer.Update fail: %v", err.Error())
		//log.Errorf( "call Transfer.Update fail: %v", err.Error())
	}
	//endTime := time.Now()
	//log.Debugf("Send metrics to transfer running in the background. Process time: %s, Send Metrics: %d",endTime.Sub(startTime).String(), len(metrics))
}