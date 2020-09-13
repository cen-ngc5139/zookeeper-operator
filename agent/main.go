package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ghostbaby/zk-agent/g"
	"github.com/ghostbaby/zk-agent/http"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	if g.Config().Debug {
		g.InitLog("debug")
	} else {
		g.InitLog("info")
	}

	go http.Start()

	select {}

}
