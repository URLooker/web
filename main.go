package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"

	"github.com/peng19940915/urlooker/web/api"
	"github.com/peng19940915/urlooker/web/cron"
	"github.com/peng19940915/urlooker/web/g"
	"github.com/peng19940915/urlooker/web/http"
	"github.com/peng19940915/urlooker/web/http/cookie"
	"github.com/peng19940915/urlooker/web/sender"
	"github.com/peng19940915/urlooker/web/store"
	"github.com/gin-gonic/gin"
	"os/signal"
	"syscall"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func init() {
	prepare()
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleConfig(*cfg)

	store.InitMysql()
	cron.Init()
	sender.Init()
	cookie.Init()
}

func main() {
	go api.Start()
	router := gin.Default()
	go http.StartGin("0.0.0.0:1984", router)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()
	select {}
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	err := g.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
