package telegram

import (
	"github.com/rogue0026/task-tracker/internal/config"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"os"
)

// Здесь объявлены все кнопки, используемые в интерфейсе бота

type Bot struct {
	api          *tele.Bot
	Logger       *logrus.Logger
	UserSessions Sessions
}

func NewBot(cfg config.BotCfg, env string) (*Bot, error) {
	var botSettings tele.Settings
	var botLogger *logrus.Logger
	switch env {
	case "dev":
		botSettings = tele.Settings{
			Token:   cfg.Token,
			Poller:  &tele.LongPoller{},
			Verbose: true,
		}
		botLogger = &logrus.Logger{
			Level: logrus.DebugLevel,
			Out:   os.Stdout,
			Formatter: &logrus.TextFormatter{
				DisableLevelTruncation: true,
				ForceColors:            true,
			},
		}
	case "prod":
		botSettings = tele.Settings{
			Token:  cfg.Token,
			Poller: &tele.LongPoller{},
		}
		botLogger = &logrus.Logger{
			Level:     logrus.InfoLevel,
			Out:       os.Stdout,
			Formatter: &logrus.JSONFormatter{},
		}
	}

	// sessions хранит состояние сессии с конкретным пользователем
	sessions := make(map[int64]*Session)

	api, err := tele.NewBot(botSettings)
	if err != nil {
		return nil, err
	}
	b := Bot{
		api:          api,
		Logger:       botLogger,
		UserSessions: sessions,
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
	b.api.Handle("/start", b.StartCommandHandler)
	b.api.Handle(&HelpButton, b.HelpButtonHandler)
	b.api.Handle(&BackButton, b.StartCommandHandler)
	b.api.Handle(&ContactsButton, b.ContactsButtonHandler)
	b.api.Handle(&TasksButton, b.TasksButtonHandler)
	b.api.Handle(&DonateButton, b.DonateButtonHandler)
}

func (b *Bot) DeletePreviousMessage() {

}
