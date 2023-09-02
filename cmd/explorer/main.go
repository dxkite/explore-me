package main

import (
	"log"
	"os"
	"time"

	"dxkite.cn/explorer/src"
	"dxkite.cn/explorer/src/core"
)

func Async(filename string, cfg *core.Config, ticker *time.Ticker) {
	for range ticker.C {
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

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	cfg, err := core.LoadConfig(filename)
	if err != nil {
		panic(err)
	}

	log.Println("init index", cfg.DataRoot, "scan", cfg.SrcRoot)
	if err := core.InitIndex(cfg); err != nil {
		log.Fatalln("InitIndexErr", err)
		return
	}

	go Async(filename, cfg, ticker)

	log.Println("start server")
	if err := src.Run(cfg); err != nil {
		log.Fatalln("StartServerErr", err)
	}
}
