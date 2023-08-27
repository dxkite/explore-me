package main

import (
	"log"
	"os"

	"dxkite.cn/explorer/src"
	"dxkite.cn/explorer/src/core"
)

func main() {
	filename := "./config.yaml"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	cfg, err := core.LoadConfig(filename)
	if err != nil {
		panic(err)
	}

	log.Println("init index", cfg.DataRoot, "scan", cfg.SrcRoot)
	if err := core.InitIndex(cfg); err != nil {
		log.Fatalln("InitIndexErr", err)
		return
	}

	log.Println("start server")
	if err := src.Run(cfg); err != nil {
		log.Fatalln("StartServerErr", err)
	}
}
