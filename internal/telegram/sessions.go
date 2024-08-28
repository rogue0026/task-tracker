package telegram

import "gopkg.in/telebot.v3"

// Sessions хранит данные сессий с пользователями где ключ id чата
type Sessions map[int64]*Session

type Session struct {
	CurrentBotState int64
	LastMessage     *telebot.Message
}

const (
	IdleInMainMenu int64 = 1<<60 + iota
	WaitingTaskNameInputFromUser
	WaitingTaskDateInputFromUser
)

func NewSession(botMsgID *telebot.Message) *Session {
	s := Session{
		CurrentBotState: IdleInMainMenu,
		LastMessage:     botMsgID,
	}
	return &s
}
