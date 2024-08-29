package telegram

import (
	"github.com/rogue0026/task-tracker/internal/models"
	"gopkg.in/telebot.v3"
	"sync"
)

const (
	IdleInMainMenu int64 = 1<<60 + iota
	WaitingTaskNameInputFromUser
	WaitingTaskDateInputFromUser
)

type Session struct {
	CurrentBotState int64
	LastMessage     *telebot.Message
	UserTasksNames  []models.Task
}

func NewSession(botMsgID *telebot.Message) Session {
	s := Session{
		CurrentBotState: IdleInMainMenu,
		LastMessage:     botMsgID,
		UserTasksNames:  make([]models.Task, 0),
	}
	return s
}

// SessionsStorage хранит данные сессий с пользователями где ключ id чата
type SessionsStorage struct {
	mu       *sync.Mutex
	sessions map[int64]*Session
}

func NewSessionsStorage() *SessionsStorage {
	ss := &SessionsStorage{
		mu:       &sync.Mutex{},
		sessions: make(map[int64]*Session),
	}
	return ss
}

func (ss *SessionsStorage) AddSession(chatID int64, s Session) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	_, ok := ss.sessions[chatID]
	if !ok {
		ss.sessions[chatID] = &s
	}
}

func (ss *SessionsStorage) GetSession(chatID int64) (*Session, bool) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	session, ok := ss.sessions[chatID]
	return session, ok
}
