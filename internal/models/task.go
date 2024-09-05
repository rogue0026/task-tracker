package models

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID       int64
	Name     string
	Deadline time.Time
	UserID   int64
}

func (t *Task) String() string {
	firstLetter := strings.ToUpper(string([]rune(t.Name)[0]))
	suffix := string([]rune(t.Name)[1:])
	nameCorrected := firstLetter + suffix
	return fmt.Sprintf("%s %s", nameCorrected, t.Deadline.Format("02.01.2006 15:04:05"))
}
