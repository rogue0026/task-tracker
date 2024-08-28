package main

import (
	"flag"
	"fmt"
	"github.com/rogue0026/task-tracker/internal/config"
	"github.com/rogue0026/task-tracker/internal/telegram"
	"os"
	"os/signal"
	"syscall"
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
	tgBot, err := telegram.NewBot(botCfg, "dev")
	if err != nil {
		panic(err.Error())
	}

	// Создаем канал для прослушивания сигналов приложению от операционной системы
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем нашего бота
	go func() {
		fmt.Println("starting telegram bot")
		tgBot.Start()
	}()

	// Здесь начинаем слушать сигналы от операционной системы
	// Как только придет один из сигналов SIGINT или SIGTERM
	s := <-stopSignal
	fmt.Printf("Получен сигнал %s\n", s.String())
	err = tgBot.Shutdown()
	if err != nil {
		fmt.Printf("При отключении бота произошла ошибка: %s\n", err.Error())
	} else {
		fmt.Println("Бот успешно выключен")
	}
}
