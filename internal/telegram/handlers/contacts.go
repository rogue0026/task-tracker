package handlers

import tele "gopkg.in/telebot.v3"

var ContactsButtonHandler = func(c tele.Context) error {
	err := c.Send("Если возникли проблемы в работе бота, пиши на @paul35426")
	if err != nil {
		return err
	}
	return nil
}
