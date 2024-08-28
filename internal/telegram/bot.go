package telegram

import (
	"github.com/rogue0026/task-tracker/internal/config"
	"github.com/rogue0026/task-tracker/internal/telegram/handlers"
	tele "gopkg.in/telebot.v3"
)

// Здесь объявлены все кнопки, используемые в интерфейсе бота

type Bot struct {
	api *tele.Bot
}

func NewBot(cfg config.BotCfg, env string) (*Bot, error) {
	var botSettings tele.Settings
	switch env {
	case "dev":
		botSettings = tele.Settings{
			Token:   cfg.Token,
			Poller:  &tele.LongPoller{},
			Verbose: true,
		}
	case "prod":
		botSettings = tele.Settings{
			Token:  cfg.Token,
			Poller: &tele.LongPoller{},
		}
	}

	api, err := tele.NewBot(botSettings)
	if err != nil {
		return nil, err
	}
	b := Bot{
		api: api,
	}
	return &b, nil
}

func (b *Bot) Start() {
	b.registerHandlers()
	b.api.Start()
}

func (b *Bot) Stop() {
	b.api.Stop()
}

func (b *Bot) Shutdown() error {
	ok, err := b.api.Close()
	if ok {
		return nil
	}
	return err
}

func (b *Bot) registerHandlers() {
	b.api.Handle("/start", handlers.StartCommandHandler)
	b.api.Handle(&handlers.HelpButton, handlers.HelpButtonHandler)
	b.api.Handle(&handlers.BackButton, handlers.StartCommandHandler)
	b.api.Handle(&handlers.ContactsButton, handlers.ContactsButtonHandler)
	b.api.Handle(&handlers.TasksButton, handlers.TasksButtonHandler)
	b.api.Handle(&handlers.DonateButton, handlers.DonateButtonHandler)
}
