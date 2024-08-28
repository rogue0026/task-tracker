package main

import (
	"flag"
	"fmt"
	"github.com/rogue0026/task-tracker/internal/config"
)

var (
	cfgPath string
)

func main() {
	flag.StringVar(&cfgPath, "c", "", "path to bot config")
	flag.Parse()
	botCfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(botCfg.Token)
}
