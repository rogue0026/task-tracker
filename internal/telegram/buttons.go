package telegram

import (
	tele "gopkg.in/telebot.v3"
)

// Здесь объявлены все кнопки, используемые в интерфейсе бота
var (
	// HelpButton отображает вспомогательную информацию о пользовании ботом
	HelpButton = tele.InlineButton{
		Unique: "HelpButton",
		Text:   "ℹ️Инфоℹ️",
		Data:   "Помощь",
	}
	BackButton = tele.InlineButton{
		Unique: "BackButton",
		Text:   "🔙Назад🔙",
		Data:   "Назад",
	}
	// ContactsButton отображает контакты разработчика бота
	ContactsButton = tele.InlineButton{
		Unique: "ContactsButton",
		Text:   "📩Поддержка📩",
		Data:   "Поддержка",
	}
	// TasksButton отображает меню управления задачами пользователя
	TasksButton = tele.InlineButton{
		Unique: "TasksButton",
		Text:   "🗓Задачи🗓",
		Data:   "Задачи",
	}

	// DonateButton отображает платежные реквизиты для пожертвований, на которые будет осуществляться дальнейшее развитие бота
	DonateButton = tele.InlineButton{
		Unique: "DonateButton",
		Text:   "💰Donate💰",
		Data:   "Donate",
	}

	// CreateTaskButton запускает процесс создания новой задачи
	CreateTaskButton = tele.InlineButton{
		Unique: "CreateTaskButton",
		Text:   "Добавить задачу",
		Data:   "Добавить задачу",
	}

	// DeleteTaskButton запускает процесс удаления существующей задачи
	DeleteTaskButton = tele.InlineButton{
		Unique: "DeleteTaskButton",
		Text:   "Удалить задачу",
		Data:   "Удалить задачу",
	}

	// ShowAllTasksButton отображает все созданные для текущего пользователя задачи
	ShowAllTasksButton = tele.InlineButton{
		Unique: "ShowAllTasksButton",
		Text:   "Все задачи",
		Data:   "Все задачи",
	}

	StartTrackingButton = tele.InlineButton{
		Unique: "StartTrackingButton",
		Text:   "Отслеживать",
		Data:   "Отслеживать",
	}
)
