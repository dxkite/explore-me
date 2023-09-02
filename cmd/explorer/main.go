package main

import (
	"log"
	"os"
	"time"

	"dxkite.cn/explorer/src"
	"dxkite.cn/explorer/src/core"
)

func Async(filename string, ticker *time.Ticker) {
	for range ticker.C {
		core.LoadConfig(filename)
		cfg := core.GetConfig()
		log.Println("init index", cfg.DataRoot, "scan", cfg.SrcRoot)
		if err := core.InitIndex(cfg); err != nil {
			log.Fatalln("InitIndexErr", err)
			return
		}
	}
}

func main() {
	filename := "./config.yaml"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	err := core.InitConfig(filename)
	if err != nil {
		panic(err)
	}

	cfg := core.GetConfig()
	log.Println("init index", cfg.DataRoot, "scan", cfg.SrcRoot)
	if err := core.InitIndex(cfg); err != nil {
		log.Fatalln("InitIndexErr", err)
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(cfg.AsyncLoad))
	defer ticker.Stop()

	go Async(filename, ticker)

	log.Println("start server")
	if err := src.Run(cfg); err != nil {
		log.Fatalln("StartServerErr", err)
	}
}
