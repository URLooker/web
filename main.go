package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/urlooker/web/api"
	"github.com/urlooker/web/cron"
	"github.com/urlooker/web/g"
	"github.com/urlooker/web/http"
	"github.com/urlooker/web/http/cookie"
	"github.com/urlooker/web/sender"
	"github.com/urlooker/web/store"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
	http.Start()
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
