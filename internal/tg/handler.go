package tg

import (
	"nstu/internal/logger"
	"nstu/internal/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(app *BotStruct, update tgbotapi.Update) error

const (
	helloMsgText = `
👋 Добро пожаловать!

🗓 Просмотр расписания:
- Используйте кнопки "Сегодня", "Завтра" или "Все расписание"
- Для конкретной даты отправьте: 14.03.2023 или 14-03-2023
- Краткий формат также работает: 14.03 или 14-03

👥 Управление группами:
- Вступите в группу через меню для просмотра её расписания
- Создавайте свои группы и настраивайте расписание
- Приглашайте одногруппников — любой пользователь может присоединиться

📋 Полезные команды:
- /help - показать это приветственное сообщение
- /open - открыть клавиатурное меню
- /close - закрыть клавиатурное меню
- /menu - сообщение с переходом в веб меню

⌨️ Клавиатурное меню:
- Сегодня - расписание на сегодня
- Завтра - расписание на завтра
- Все расписание - полное расписание группы
- Войти в меню - сообщение с переходом в веб меню

💡 Совет: Обязательно вступите в группу, прежде чем просматривать расписание!

🚀 Готовы начать? Используйте команду /open для доступа к функциям бота.
`
)

// handleStartMessage пересылает пользователю стикер в ответ на начало работы
func HandleStartMessage() HandlerFunc {
	return func(app *BotStruct, update tgbotapi.Update) error {

		// Отправка сообщения
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMsgText)
		msg.ParseMode = "html"
		sendedMsg, err := app.SendMessage(msg)
		if err != nil {
			logger.Log.Error("Ошибка при отправке сообщения: ", err)
			return err
		}

		// Открепление всех закрепов
		unpinConfig := tgbotapi.UnpinAllChatMessagesConfig{
			ChatID:          update.Message.Chat.ID,
			ChannelUsername: update.Message.From.UserName,
		}
		_, err = app.BotAPI.Request(unpinConfig)
		if err != nil {
			return err
		}

		if update.Message.Chat.IsGroup() {
			return nil
		}

		// Закрепление отправленного сообщения
		pinConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              update.Message.Chat.ID,
			MessageID:           sendedMsg.MessageID,
			DisableNotification: false, // Если true, уведомление о закреплении не будет отправлено
		}

		_, err = app.BotAPI.Request(pinConfig)
		if err != nil {
			return err
		}

		return nil
	}
}

func HandleSchedule(date time.Time) HandlerFunc {
	return func(app *BotStruct, update tgbotapi.Update) error {

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			logger.Log.Error("Не удалось конвертировать время : ", err)
		}

		date = date.In(loc)

		userId := update.Message.From.ID

		msgText := ""

		user, err := models.DataBase.GetUser(userId)
		if err != nil {
			logger.Log.Info("Пользователь не зарегистрирован в БД: ", err)

			_, err = app.SendMessage(
				tgbotapi.NewMessage(update.Message.Chat.ID,
					"Перед запросом расписания необходимо войти в группу",
				))
			if err != nil {
				logger.Log.Error("Не удалось отправить сообщение: ", err)
			}
			return err
		}

		if *user.ConsistsOf == 0 {
			logger.Log.Info("Пользователь не зарегистрирован в БД: ", err)

			_, err = app.SendMessage(
				tgbotapi.NewMessage(update.Message.Chat.ID,
					"Перед запросом расписания необходимо войти в группу",
				))
			if err != nil {
				logger.Log.Error("Не удалось отправить сообщение: ", err)
			}
			return err
		}

		group, err := models.DataBase.GetGroup(uint(*user.ConsistsOf))
		if err != nil {
			logger.Log.Info("Не удалось получить группу пользователя: ", err)
			msgText = "Группа не найдена. Проверьте существование группы."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		msg.ParseMode = "html"

		shedule, err := models.GetScheduleOnDate(date, group)
		if err != nil {
			logger.Log.Info("Не удалось получить расписание группы на определенный день: ", err)
			msgText = "Не удалось получить расписание группы на этот день."
		}

		if err == nil {
			msgText += "<strong>Расписание на " + date.Format("02.01.2006") + "(+" + models.GetDayOfWeekInRussian(date) + ")</strong>"
			msgText += models.FormatDaySchedule(shedule)
		}

		_, err = app.SendMessage(msg)
		if err != nil {
			logger.Log.Error("Не удалось отправить сообщение: ", err)
		}

		return nil
	}
}

func HandleAllSchedule() HandlerFunc {
	return func(app *BotStruct, update tgbotapi.Update) error {
		userId := update.Message.From.ID
		msgText := ""

		user, err := models.DataBase.GetUser(userId)
		if err != nil {
			logger.Log.Info("Пользователь не зарегистрирован в БД: ", err)

			_, err = app.SendMessage(
				tgbotapi.NewMessage(update.Message.Chat.ID,
					"Перед запросом расписания необходимо войти в группу",
				))
			if err != nil {
				logger.Log.Error("Не удалось отправить сообщение: ", err)
			}
			return err
		}
		if *user.ConsistsOf == 0 {
			logger.Log.Info("Пользователь не зарегистрирован в БД: ", err)

			_, err = app.SendMessage(
				tgbotapi.NewMessage(update.Message.Chat.ID,
					"Перед запросом расписания необходимо войти в группу",
				))
			if err != nil {
				logger.Log.Error("Не удалось отправить сообщение: ", err)
			}
			return err
		}

		group, err := models.DataBase.GetGroup(uint(*user.ConsistsOf))
		if err != nil {
			logger.Log.Info("Не удалось получить группу пользователя: ", err)
			msgText = "Группа не найдена. Проверьте существование группы."
		}

		if err == nil {
			msgText += "<strong>Все расписание</strong>"
			msgText += models.FormatAllSchedule(group)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

		_, err = app.SendMessage(msg)
		if err != nil {
			logger.Log.Error("Не удалось отправить сообщение: ", err)
		}

		return nil
	}
}

func HandleClose() HandlerFunc {
	return func(app *BotStruct, update tgbotapi.Update) error {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Закрыто")
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)

		// Устанавливаем его в качестве клавиатуры для сообщения
		msg.ReplyMarkup = removeKeyboard

		_, err := app.SendMessage(msg)
		if err != nil {
			logger.Log.Error("Не удалось отправить сообщение: ", err)
		}

		return nil
	}
}

func HandleOpen() HandlerFunc {
	return func(app *BotStruct, update tgbotapi.Update) error {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Открыто")

		removeKeyboard := CreateReplyKeyboard([][]ButtonData{
			{
				{Text: "Сегодня", Data: "Сегодня"},
				{Text: "Завтра", Data: "Завтра"},
			},
			{
				{Text: "Все расписание", Data: "Все расписание"},
				{Text: "Войти в меню", Data: "Войти в меню"},
			},
		})

		// Устанавливаем его в качестве клавиатуры для сообщения
		msg.ReplyMarkup = removeKeyboard

		_, err := app.SendMessage(msg)
		if err != nil {
			logger.Log.Error("Не удалось отправить сообщение: ", err)
		}

		return nil
	}
}

func HandleEnterMenu() HandlerFunc {
	return func(app *BotStruct, update tgbotapi.Update) error {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Чтобы, войти в меню, нажмите на кнопку")
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonURL("Открыть меню", "https://www.youtube.com/watch?v=2_bWhf5SdtI&t=721s"),
		}
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row)
		msg.ParseMode = "html"

		_, err := app.SendMessage(msg)
		if err != nil {
			return err
		}

		return nil
	}
}
