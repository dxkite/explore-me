package main

import (
	"os"
	"time"

	"dxkite.cn/log"

	"dxkite.cn/explore-me/src"
	"dxkite.cn/explore-me/src/core"
	"dxkite.cn/explore-me/src/core/config"
)

func Async(filename string, ticker *time.Ticker) {
	for range ticker.C {
		config.LoadConfig(filename)
		cfg := config.GetConfig()
		if err := core.CreateIndex(cfg); err != nil {
			log.Fatalln("load index error", err)
			return
		}
	}
}

func main() {
	filename := "./config.yaml"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	err := config.InitConfig(filename)
	if err != nil {
		panic(err)
	}

	cfg := config.GetConfig()
	log.Println("init index", cfg.DataRoot, "scan", cfg.SrcRoot)
	if err := core.CreateIndex(cfg); err != nil {
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
