package telegram

import tele "gopkg.in/telebot.v3"

func (b *Bot) CheckRegistration(next tele.HandlerFunc) tele.HandlerFunc {
	const fn = "CheckRegistration"
	return func(c tele.Context) error {
		_, ok := b.UserSessions.SessionByID(c.Chat().ID)
		if ok {
			next(c)
		} else {
			err := c.Send("Перезапусти бота")
			if err != nil {
				b.Logger.Errorf("func=%s error=%s", fn, err.Error())
				return err
			}
		}
		return nil
	}
}
