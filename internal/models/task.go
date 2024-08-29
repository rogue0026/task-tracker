package models

import (
	"fmt"
	"time"
)

const TimeParseLayout string = "02.01.2006"

type Task struct {
	ID          int64
	Name        string
	IsCompleted bool
	Deadline    time.Time
}

func (t *Task) String() string {
	if t.IsCompleted == true {
		return fmt.Sprintf("%s, %s, выполнено", t.Name, t.Deadline.Format(TimeParseLayout))
	} else {
		return fmt.Sprintf("%s, %s, не выполнено", t.Name, t.Deadline.Format(TimeParseLayout))
	}
}
