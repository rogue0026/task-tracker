package telegram

import (
	"errors"
	"fmt"
	"github.com/rogue0026/task-tracker/internal/config"
	"github.com/rogue0026/task-tracker/internal/storage/postgres"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"os"
)

const (
	// Bot states

	IdleInMainMenu                string = "IdleInMainMenu"
	WaitingTaskNameInputFromUser  string = "WaitingTaskNameInputFromUser"
	WaitingTaskDateInputFromUser  string = "WaitingTaskDateInputFromUser"
	WaitingTaskIDForDeleteTask    string = "WaitingTaskIDForDeleteTask"
	WaitingTaskIDForSetupTracking string = "WaitingTaskIDForSetupTracking"

	TimeParseLayout string = "15.04 02.01.2006"
)

var (
	ErrUserSessionNotFound = errors.New("user session not found")
)

type Bot struct {
	api          *tele.Bot
	Logger       *logrus.Logger
	UserSessions *SessionsStorage
	Tasks        *postgres.TasksStorage
}

func NewBot(cfg config.BotCfg, env string) (*Bot, error) {
	const fn = "internal.telegram.NewBot"
	var botSettings tele.Settings
	var botLogger *logrus.Logger

	switch env {
	case "dev":
		botSettings = tele.Settings{
			Token:   cfg.Token,
			Poller:  &tele.LongPoller{},
			Verbose: false,
		}

		botLogger = &logrus.Logger{
			Level: logrus.DebugLevel,
			Out:   os.Stdout,
			Formatter: &logrus.TextFormatter{
				DisableLevelTruncation: true,
				ForceColors:            true,
			},
			ReportCaller: true,
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

	sessions := NewSessionsStorage()

	api, err := tele.NewBot(botSettings)
	if err != nil {
		return nil, err
	}

	tasks, err := postgres.New(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("func=%s, error=%w", fn, err)
	}
	b := Bot{
		api:          api,
		Logger:       botLogger,
		UserSessions: sessions,
		Tasks:        &tasks,
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
	b.api.Handle("/menu", b.StartCommandHandler)
	b.api.Handle("/start", b.StartCommandHandler)
	b.api.Handle(&HelpButton, b.HelpButtonHandler, b.CheckRegistration)
	b.api.Handle(&BackButton, b.StartCommandHandler, b.CheckRegistration)
	b.api.Handle(&ContactsButton, b.ContactsButtonHandler, b.CheckRegistration)
	b.api.Handle(&TasksButton, b.TasksButtonHandler, b.CheckRegistration)
	b.api.Handle(&DonateButton, b.DonateButtonHandler, b.CheckRegistration)
	b.api.Handle(&CreateTaskButton, b.CreateTaskHandler, b.CheckRegistration)
	b.api.Handle(tele.OnText, b.UserInputHandler, b.CheckRegistration)
	b.api.Handle(&ShowAllTasksButton, b.ShowAllTasksButtonHandler, b.CheckRegistration)
	//b.api.Handle(&DeleteTaskButton, b.DeleteTaskButtonHandler, b.CheckRegistration)
}
