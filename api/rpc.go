package api

import (
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/peng19940915/urlooker/web/g"
)

type Web int

func Start() {
	addr := g.Config.Rpc.Listen

	server := rpc.NewServer()
	server.Register(new(Web))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen occur error", e)
	} else {
		log.Info("listening on", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorf("listener accept occur error, detail", err.Error())
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (this *Web) Ping(req interface{}, reply *string) error {
	*reply = "ok"
	return nil
}
