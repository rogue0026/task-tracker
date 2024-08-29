package telegram

import (
	"errors"
	"github.com/rogue0026/task-tracker/internal/config"
	"github.com/rogue0026/task-tracker/internal/models"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	api          *tele.Bot
	Logger       *logrus.Logger
	UserSessions *SessionsStorage
}

func NewBot(cfg config.BotCfg, env string) (*Bot, error) {
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

	// sessions хранит состояния сессий с конкретными пользователями
	sessions := NewSessionsStorage()

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
	b.api.Handle(&CreateTaskButton, b.CreateTaskHandler)
	b.api.Handle(tele.OnText, b.UserInputHandler)
	b.api.Handle(&ShowAllTasksButton, b.ShowAllTasksButtonHandler)
	b.api.Handle(&DeleteTaskButton, b.DeleteTaskButtonHandler)
}

func (b *Bot) StartCommandHandler(c tele.Context) error {
	const fn = "StartCommandHandler"
	keyboard := tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{TasksButton},
			{ContactsButton, HelpButton},
			{DonateButton},
		},
	}
	// Sending to user welcome message
	sentMsg, err := b.api.Send(c.Chat(), "Привет, этот бот поможет тебе отслеживать ежедневные задачи", &keyboard)
	if err != nil {
		return err
	}
	// Check if session with current user exists
	_, ok := b.UserSessions.GetSession(c.Chat().ID)
	// If no - create new session and add it to storage
	if !ok {
		b.Logger.Debugf("session created: chat_id=%v", c.Sender().ID)
		usrSession := NewSession(sentMsg)
		b.UserSessions.AddSession(c.Chat().ID, usrSession)
	} else { // If yes - we have to delete previous message sent by bot
		usrSession, _ := b.UserSessions.GetSession(c.Chat().ID)
		if err := b.api.Delete(usrSession.LastMessage); err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
		}
		usrSession.LastMessage = sentMsg
	}
	return nil
}

// HelpButtonHandler shows bot usage instruction
func (b *Bot) HelpButtonHandler(c tele.Context) error {
	err := c.Send("Здесь будет краткое руководство по использованию бота")
	if err != nil {
		return err
	}
	return nil
}

// ContactsButtonHandler shows tech support contacts
func (b *Bot) ContactsButtonHandler(c tele.Context) error {
	err := c.Send("Если возникли проблемы в работе бота, пиши на @paul35426")
	if err != nil {
		return err
	}
	return nil
}

// TasksButtonHandler вызывается при нажатии на кнопку TasksButton
func (b *Bot) TasksButtonHandler(c tele.Context) error {
	const fn = "TasksButtonHandler"
	keyboard := tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{CreateTaskButton},
			{DeleteTaskButton},
			{ShowAllTasksButton},
			{BackButton},
		},
	}
	// Sending message to user
	sentMsg, err := b.api.Send(c.Chat(), "Вы вошли в режим управления задачами", &keyboard)
	if err != nil {
		b.Logger.Errorf("func=%s error=%s", fn, err.Error())
		return err
	}
	if usrSession, ok := b.UserSessions.GetSession(c.Chat().ID); ok {
		b.Logger.Debugf("func=%s. Session found, chat_id=%v", fn, c.Chat().ID)
		// Deleting previous message sent by telegram bot
		if usrSession.LastMessage != nil {
			err := b.api.Delete(usrSession.LastMessage)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
			usrSession.LastMessage = sentMsg
		} else {
			b.Logger.Errorf("func=%s error: chat_id=%v, user session is empty", fn, c.Chat().ID)
		}
	}
	return nil
}

func (b *Bot) DonateButtonHandler(c tele.Context) error {
	return c.Send("Здесь будут реквизиты для пожертвований на развитие бота")
}

// CreateTaskHandler отправляет пользователю сообщение с просьбой отправить название задачи и переводи состояние сессии в режим ожидания сообщения от пользователя
func (b *Bot) CreateTaskHandler(c tele.Context) error {
	const fn = "CreateTaskHandler"
	usrSession, ok := b.UserSessions.GetSession(c.Chat().ID)
	if ok {
		err := c.Send("Отправь мне название задачи, которое ты хочешь сохранить")
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		// We change session state to wait for task name from user
		usrSession.CurrentBotState = WaitingTaskNameInputFromUser
		b.Logger.Debug("session state was changed to waiting for task name from user")
	} else {
		b.Logger.Errorf("func=%s error=%s", fn, "user session not found")
	}
	return nil
}

func (b *Bot) UserInputHandler(c tele.Context) error {
	const fn = "UserInputHandler"
	usrSession, ok := b.UserSessions.GetSession(c.Chat().ID)
	if ok {
		switch usrSession.CurrentBotState {
		// This case must be executed if bot is in waiting for user task name
		case WaitingTaskNameInputFromUser:
			inMsg := c.Message()
			b.Logger.Debugf("received message: %s", inMsg.Text)
			// If ok - save message text to tasks storage
			if inMsg != nil {
				t := models.Task{
					ID:   usrSession.FreeTaskID,
					Name: inMsg.Text,
				}
				usrSession.FreeTaskID++
				usrSession.UserTasks = append(usrSession.UserTasks, t)
				c.Send("Отлично, теперь отправь мне время и дату в формате ЧЧ.ММ ДД.ММ.ГГГГ, до которого ты должен успеть выполнить поставленную задачу")
				usrSession.CurrentBotState = WaitingTaskDateInputFromUser
			} else {
				b.Logger.Errorf("func=%s error=input message is nil", fn)
				return errors.New("input message is nil")
			}
		case WaitingTaskDateInputFromUser:
			// parsing input message to time.Time struct
			deadline, err := time.Parse(models.TimeParseLayout, c.Message().Text)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
			// adding deadline info to task
			usrSession.UserTasks[len(usrSession.UserTasks)-1].Deadline = deadline
			if deadline.After(time.Now()) {
				usrSession.UserTasks[len(usrSession.UserTasks)-1].Status = models.NotCompleted
			} else {
				usrSession.UserTasks[len(usrSession.UserTasks)-1].Status = models.Failed
			}
			c.Send("Отлично, задача добавлена в список для отслеживания")
			time.Sleep(time.Second * 1)

			// Sending message to user for further task management
			keyboard := tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{CreateTaskButton},
					{DeleteTaskButton},
					{ShowAllTasksButton},
					{BackButton},
				},
			}
			// Sending message to user
			sentMsg, err := b.api.Send(c.Chat(), "Режим управления задачами", &keyboard)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
			err = b.api.Delete(usrSession.LastMessage)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
			usrSession.LastMessage = sentMsg

			usrSession.CurrentBotState = IdleInMainMenu

		case WaitingTaskIDForDeleteTask:
			taskID, err := strconv.ParseInt(c.Message().Text, 10, 64)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				err = c.Send("Вы неверно указали номер задачи. Отправьте номер задачи повторно.")
				if err != nil {
					b.Logger.Errorf("func=%s error=%s", fn, err.Error())
					return err
				}
			} else {
				for i := range usrSession.UserTasks {
					if usrSession.UserTasks[i].ID == taskID {
						result := slices.Concat(usrSession.UserTasks[:i], usrSession.UserTasks[i+1:])
						usrSession.UserTasks = result
					}
				}

				err = c.Send("Задача удалена")
				if err != nil {
					b.Logger.Errorf("func=%s error=%s", fn, err.Error())
					return err
				}
				usrSession.CurrentBotState = IdleInMainMenu
				// Sending message to user for further task management
				keyboard := tele.ReplyMarkup{
					InlineKeyboard: [][]tele.InlineButton{
						{CreateTaskButton},
						{DeleteTaskButton},
						{ShowAllTasksButton},
						{BackButton},
					},
				}
				// Sending message to user
				sentMsg, err := b.api.Send(c.Chat(), "Режим управления задачами", &keyboard)
				if err != nil {
					b.Logger.Errorf("func=%s error=%s", fn, err.Error())
					return err
				}
				err = b.api.Delete(usrSession.LastMessage)
				if err != nil {
					b.Logger.Errorf("func=%s error=%s", fn, err.Error())
					return err
				}
				usrSession.LastMessage = sentMsg
			}
		}
	} else {
		b.Logger.Errorf("func=%s error=%s", fn, "user session not found")
	}
	return nil
}

func (b *Bot) ShowAllTasksButtonHandler(c tele.Context) error {
	const fn = "ShowAllTasksButtonHandler"
	usrSession, ok := b.UserSessions.GetSession(c.Chat().ID)
	if ok {
		builder := strings.Builder{}
		for _, t := range usrSession.UserTasks {
			builder.WriteString(t.String() + "\n")
		}
		if len(builder.String()) == 0 {
			err := c.Send("Пока что нет задач")
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
		} else {
			err := c.Send(builder.String())
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
		}
	}
	return nil
}

func (b *Bot) DeleteTaskButtonHandler(c tele.Context) error {
	const fn = "DeleteTaskButtonHandler"
	usrSession, ok := b.UserSessions.GetSession(c.Chat().ID)
	if ok {
		if len(usrSession.UserTasks) == 0 {
			err := c.Send("Нет задач для удаления")
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
		} else {
			err := c.Send("Отправь мне порядковый номер задачи, которую ты хочешь удалить")
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
			usrSession.CurrentBotState = WaitingTaskIDForDeleteTask
		}
	} else {
		b.Logger.Errorf("func=%s error=%s", fn, "user session not found")
	}
	return nil
}
