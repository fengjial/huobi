package main

import (
	conf "com/mc2/trade/config"
	"com/mc2/trade/log"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	help       *bool   = flag.Bool("h", false, "to show help")
	confPrefix *string = flag.String("c", "../config", "root path of config file")
)

// abnormal exit
func abnormalExit() {
	/* to overcome bug in log, sleep for a while */
	time.Sleep(1 * time.Second)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}
	conf.SetPrefix(confPrefix)
	err := log.Init(false)
	if err != nil {
		fmt.Printf("init err log.Init():%s\n", err.Error())
		abnormalExit()
	}
	defer log.Logger.Close()
	log.Logger.Error("%s", "this is a test")
	log.Logger.Info("hello,word")
	fmt.Printf("hello, world\n")
	fmt.Printf("%+v\n", log.Logger)
}
