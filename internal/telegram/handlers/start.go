package handlers

import (
	tele "gopkg.in/telebot.v3"
)

var StartCommandHandler = func(c tele.Context) error {
	keyboard := tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{TasksButton},
			{ContactsButton, HelpButton},
			{DonateButton},
		},
	}
	err := c.Send("Привет, этот бот поможет тебе отслеживать ежедневные задачи", &keyboard)
	if err != nil {
		return err
	}
	return nil
}
