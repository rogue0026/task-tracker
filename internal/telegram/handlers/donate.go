package handlers

import tele "gopkg.in/telebot.v3"

var DonateButtonHandler = func(c tele.Context) error {
	return c.Send("Здесь будут реквизиты для пожертвований на развитие бота")
}
