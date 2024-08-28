package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (b *Bot) HelpButtonHandler(c tele.Context) error {
	err := c.Send("Здесь будет краткое руководство по использованию бота")
	if err != nil {
		return err
	}
	return nil
}
