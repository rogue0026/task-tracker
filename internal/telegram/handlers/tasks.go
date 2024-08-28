package handlers

import tele "gopkg.in/telebot.v3"

// TasksButtonHandler вызывается при нажатии на кнопку TasksButton
var TasksButtonHandler = func(c tele.Context) error {
	keyboard := tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{CreateTaskButton},
			{DeleteTaskButton},
			{ShowAllTasksButton},
			{BackButton},
		},
	}
	err := c.Send("Вы вошли в режим управления задачами", &keyboard)
	if err != nil {
		return err
	}
	return nil
}
