package models

import (
	"fmt"
	"strings"
	"time"
)

type TaskStatus int64

func (ts TaskStatus) String() string {
	switch ts {
	case Completed:
		return "Завершено"
	case NotCompleted:
		return "Не завершено"
	case Failed:
		return "Провалено"
	}
	return ""
}

const (
	Completed TaskStatus = 1<<21 + iota
	NotCompleted
	Failed
)

const TimeParseLayout string = "15.04 02.01.2006"

type Task struct {
	ID       int64
	Name     string
	Status   TaskStatus
	Deadline time.Time
}

func NewTask(id int64, name string) Task {
	t := Task{
		ID:   id,
		Name: name,
	}
	return t
}

func (t *Task) String() string {
	firstLetter := strings.ToUpper(string([]rune(t.Name)[0]))
	suffix := string([]rune(t.Name)[1:])
	nameCorrected := firstLetter + suffix
	return fmt.Sprintf("%d. %s. Срок - до %s. Статус: %s", t.ID, nameCorrected, t.Deadline.Format("15:04 02.01.2006"), t.Status)
}
