package handlers

import (
	tele "gopkg.in/telebot.v3"
)

var HelpButtonHandler = func(c tele.Context) error {
	keyboard := tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{BackButton},
		},
	}
	err := c.Send("Здесь будет краткое руководство по использованию бота", &keyboard)
	if err != nil {
		return err
	}
	return nil
}
