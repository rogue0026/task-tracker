package telegram

import tele "gopkg.in/telebot.v3"

func (b *Bot) ContactsButtonHandler(c tele.Context) error {
	err := c.Send("Если возникли проблемы в работе бота, пиши на @paul35426")
	if err != nil {
		return err
	}
	return nil
}
