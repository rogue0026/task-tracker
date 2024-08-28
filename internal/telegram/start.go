package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (b *Bot) StartCommandHandler(c tele.Context) error {
	keyboard := tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{TasksButton},
			{ContactsButton, HelpButton},
			{DonateButton},
		},
	}
	sentMsg, err := b.api.Send(c.Chat(), "Привет, этот бот поможет тебе отслеживать ежедневные задачи", &keyboard)
	if err != nil {
		return err
	}
	b.UserSessions[c.Chat().ID] = NewSession(sentMsg)
	b.Logger.Debugf("session added: user_id=%v", c.Sender().ID)
	return nil
}
