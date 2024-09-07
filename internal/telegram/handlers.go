package telegram

import (
	"errors"
	"fmt"
	"github.com/rogue0026/task-tracker/internal/models"
	"github.com/rogue0026/task-tracker/internal/storage"
	tele "gopkg.in/telebot.v3"
	"strings"
	"time"
)

const DateLayout string = "02.01.2006"

func (b *Bot) StartCommandHandler(c tele.Context) error {
	const fn = "StartCommandHandler"
	mainMenuKeyboard := tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{TasksButton},
			{ContactsButton, HelpButton},
			{DonateButton},
		},
	}

	usrSession, sessionExists := b.UserSessions.SessionByID(c.Chat().ID)

	if !sessionExists {
		// Sending to user welcome message
		builder := strings.Builder{}
		builder.WriteString("Привет! Этот бот поможет тебе отслеживать ежедневные задачи.")
		sentMsg, err := b.api.Send(c.Chat(), builder.String(), &mainMenuKeyboard)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}

		usrSession := NewSession(c.Sender().ID, sentMsg)
		usrSession.CurrentBotState = IdleInMainMenu
		b.Logger.Debugf("bot state changed to %s", IdleInMainMenu)
		b.UserSessions.AddSession(c.Chat().ID, usrSession)
		b.Logger.Debugf("session created: chat_id=%v", c.Sender().ID)
	} else {
		tasks, err := b.Tasks.UserTasks(usrSession.UserID)
		if err != nil && !errors.Is(err, storage.ErrNoTasksForUser) {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}

		builder := strings.Builder{}

		builder.WriteString("===============Workaholic==============\n")
		builder.WriteString("Дата и время: ")
		builder.WriteString(time.Now().Format("02.01.2006 15:04") + "\n")

		todayTasks := make([]models.Task, 0)

		for _, t := range tasks {
			if time.Now().Format(DateLayout) == t.Deadline.Format(DateLayout) {
				todayTasks = append(todayTasks, t)
			}
		}

		if len(todayTasks) == 0 {
			builder.WriteString("Задач на сегодня: задачи отсутствуют")
		} else {
			builder.WriteString(fmt.Sprintf("Задач на сегодня: %d\n", len(todayTasks)))
			builder.WriteString("======================================\n")
			for i := range todayTasks {
				name := todayTasks[i].Name
				remaining := todayTasks[i].Deadline.Sub(time.Now())
				h := int(remaining.Hours())
				m := int(remaining.Minutes()) % 60
				s := int(remaining.Seconds()) % 60
				builder.WriteString(fmt.Sprintf("%s %2d:%d:%d\n", name, h, m, s))
			}
			builder.WriteString("======================================\n")
		}

		sentMsg, err := b.api.Send(c.Chat(), builder.String(), &mainMenuKeyboard)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}

		usrSession, _ := b.UserSessions.SessionByID(c.Chat().ID)
		if usrSession.LastMessage != nil {
			err := b.api.Delete(usrSession.LastMessage)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
			usrSession.LastMessage = sentMsg
		}
		usrSession.CurrentBotState = IdleInMainMenu
		b.Logger.Debugf("bot state changed to %s", IdleInMainMenu)
	}

	return nil
}

// HelpButtonHandler shows bot usage instruction
func (b *Bot) HelpButtonHandler(c tele.Context) error {
	const fn = "HelpButtonHandler"
	sentMsg, err := b.api.Send(c.Chat(), "Здесь будет краткое руководство по использованию бота")
	if err != nil {
		b.Logger.Errorf("func=%s error=%s", fn, err.Error())
		return err
	}
	go func() {
		select {
		case <-time.After(time.Second * 45):
			err := b.api.Delete(sentMsg)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			} else {
				b.Logger.Debugf("func=%s msg=%s", fn, "help message successfully deleted")
			}
		}
	}()
	return nil
}

// ContactsButtonHandler shows tech support contacts
func (b *Bot) ContactsButtonHandler(c tele.Context) error {
	const fn = "ContactsButtonHandler"
	sentMsg, err := b.api.Send(c.Chat(), "Если возникли проблемы в работе бота, пиши на @paul35426")
	if err != nil {
		b.Logger.Errorf("func=%s error=%s", fn, err.Error())
		return err
	}
	go func() {
		select {
		case <-time.After(time.Second * 20):
			err := b.api.Delete(sentMsg)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			} else {
				b.Logger.Debugf("func=%s msg=%s", fn, "contact message successfully deleted")
			}
		}
	}()
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

	usrSession, _ := b.UserSessions.SessionByID(c.Chat().ID)

	// Sending message to user
	sentMsg, err := b.api.Send(c.Chat(), "Управление задачами", &keyboard)
	if err != nil {
		b.Logger.Errorf("func=%s error=%s", fn, err.Error())
		return err
	}

	// Deleting previous message sent by telegram bot
	if usrSession.LastMessage != nil {
		err := b.api.Delete(usrSession.LastMessage)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
	}
	usrSession.LastMessage = sentMsg
	return nil
}

func (b *Bot) DonateButtonHandler(c tele.Context) error {
	const fn = "DonateButtonHandler"
	sentMsg, err := b.api.Send(c.Chat(), "Для пожертвований на развитие бота:\nBTC:\n `2N3oefVeg6stiTb5Kh3ozCSkaqmx91FDbsm`", tele.ModeMarkdownV2)
	if err != nil {
		b.Logger.Errorf("func=%s error=%s", fn, err.Error())
		return err
	}
	go func() {
		select {
		case <-time.After(time.Second * 20):
			err := b.api.Delete(sentMsg)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			} else {
				b.Logger.Debugf("func=%s donate message successfully deleted", fn)
			}
		}
	}()
	return nil
}

// CreateTaskHandler отправляет пользователю сообщение с просьбой отправить название задачи и переводи состояние сессии в режим ожидания сообщения от пользователя
func (b *Bot) CreateTaskHandler(c tele.Context) error {
	const fn = "CreateTaskHandler"
	usrSession, found := b.UserSessions.SessionByID(c.Chat().ID)
	if !found {
		b.Logger.Errorf("func=%s error=%s", fn, ErrUserSessionNotFound.Error())
		return ErrUserSessionNotFound
	}

	err := c.Send("Отправь мне название задачи, которое ты хочешь сохранить")
	if err != nil {
		b.Logger.Errorf("func=%s error=%s", fn, err.Error())
		return err
	}

	// We change session state to wait for task name from user
	usrSession.CurrentBotState = WaitingTaskNameInputFromUser
	b.Logger.Debug("session state was changed to waiting for task name from user")

	return nil
}

func (b *Bot) UserInputHandler(c tele.Context) error {
	const fn = "UserInputHandler"
	usrSession, ok := b.UserSessions.SessionByID(c.Chat().ID)
	if !ok {
		b.Logger.Errorf("func=%s error=%s", fn, ErrUserSessionNotFound.Error())
		return ErrUserSessionNotFound
	}

	switch usrSession.CurrentBotState {

	// This case must be executed if bot is in waiting for user task name
	case WaitingTaskNameInputFromUser:
		inMsg := c.Message()
		// If ok - save message text to tasks storage
		if inMsg != nil {
			b.Logger.Debugf("received message from user_id %v: %s", c.Sender().ID, c.Message().Text)
			usrSession.TempTask.Name = c.Message().Text
			usrSession.TempTask.UserID = c.Sender().ID
			err := c.Send("Отлично, теперь отправь мне время и дату в формате ЧЧ.ММ ДД.ММ.ГГГГ, до которого ты должен успеть выполнить поставленную задачу")
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			}
			usrSession.CurrentBotState = WaitingTaskDateInputFromUser
		} else {
			b.Logger.Errorf("func=%s error=input message is nil", fn)
			return errors.New("input message is nil")
		}
		//return b.onWaitingTaskNameInputFromUser(usrSession, c)

	case WaitingTaskDateInputFromUser:
		// parsing input message to time.Time struct
		deadline, err := time.Parse(TimeParseLayout, c.Message().Text)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			err := c.Send("Кажется ты допустил ошибку при вводе времени и даты. Попробуй еще раз.")
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			}
			return err
		}

		// adding deadline info to task
		usrSession.TempTask.Deadline = deadline
		err = b.Tasks.SaveTask(usrSession.TempTask)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		usrSession.TempTask = models.Task{}
		err = c.Send("Отлично, задача добавлена в список для отслеживания")
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		time.Sleep(time.Millisecond * 500)

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
		sentMsg, err := b.api.Send(c.Chat(), "Управление задачами", &keyboard)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		if usrSession.LastMessage != nil {
			err = b.api.Delete(usrSession.LastMessage)
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
			usrSession.LastMessage = sentMsg
			usrSession.CurrentBotState = IdleInMainMenu
		}

	case WaitingTaskNameForDelete:
		taskName := c.Message().Text
		err := b.Tasks.DeleteTask(taskName, usrSession.UserID)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		err = c.Send("Задача удалена")
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		err = b.StartCommandHandler(c)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		usrSession.CurrentBotState = IdleInMainMenu
	}
	return nil
}

func (b *Bot) ShowAllTasksButtonHandler(c tele.Context) error {
	const fn = "ShowAllTasksButtonHandler"
	usrSession, ok := b.UserSessions.SessionByID(c.Chat().ID)
	if ok {
		tasks, err := b.Tasks.UserTasks(usrSession.UserID)
		if err != nil {
			b.Logger.Error(err.Error())
			if errors.Is(err, storage.ErrNoTasksForUser) {
				sentMsg, err := b.api.Send(c.Chat(), "У тебя нет задач")
				if err != nil {
					b.Logger.Errorf("func=%s error=%s", fn, err.Error())
					return err
				}
				go func() {
					time.Sleep(time.Second * 10)
					err := b.api.Delete(sentMsg)
					if err != nil {
						b.Logger.Errorf("func=%s error=%s", fn, err.Error())
					}
				}()
				return nil
			}
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		builder := strings.Builder{}
		for _, t := range tasks {
			builder.WriteString(t.String() + "\n")
		}
		err = c.Send(builder.String())
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
	} else {
		b.Logger.Errorf("func=%s error=%s", fn, ErrUserSessionNotFound)
	}
	return nil
}

func (b *Bot) DeleteTaskButtonHandler(c tele.Context) error {
	const fn = "DeleteTaskButtonHandler"
	usrSession, ok := b.UserSessions.SessionByID(c.Chat().ID)
	if ok {
		tasks, err := b.Tasks.UserTasks(usrSession.UserID)
		if err != nil {
			if errors.Is(err, storage.ErrNoTasksForUser) {
				err := c.Send("Нет задач для удаления")
				if err != nil {
					b.Logger.Errorf("func=%s error=%s", fn, err.Error())
					return err
				}
				return nil
			}
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		builder := strings.Builder{}
		builder.WriteString("Выбери задачу из списка и отправь мне ее название\n")
		for _, t := range tasks {
			builder.WriteString(fmt.Sprintf("`%s`\n\n", t.Name))
		}
		err = c.Send(builder.String(), tele.ModeMarkdownV2)
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		usrSession.CurrentBotState = WaitingTaskNameForDelete
	}
	return nil
}

func (b *Bot) StartTrackingButtonHandler(c tele.Context) error {
	const fn = "StartTrackingButtonHandler"
	usrSession, ok := b.UserSessions.SessionByID(c.Chat().ID)
	if ok {
		err := c.Send("Отправь мне номер задачи, которую ты хочешь отслеживать")
		if err != nil {
			b.Logger.Errorf("func=%s error=%s", fn, err.Error())
			return err
		}
		usrSession.CurrentBotState = WaitingTaskIDForSetupTracking
	} else {
		b.Logger.Errorf("func=%s error=user session not found", fn)
	}
	return nil
}
