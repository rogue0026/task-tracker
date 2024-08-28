package telegram

import tele "gopkg.in/telebot.v3"

func (b *Bot) DonateButtonHandler(c tele.Context) error {
	return c.Send("Здесь будут реквизиты для пожертвований на развитие бота")
}
