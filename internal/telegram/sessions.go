package telegram

type Sessions map[int64]*Session

type Session struct {
	LastMessage int
}
