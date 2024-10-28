package tg

import (
	"app/internal/db"
	"app/internal/logger"
	"app/internal/model"
	"app/internal/smtp"
	"app/internal/utils"
	u "app/internal/utils"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/rand"
)

type HandlerFunc func(app *Bot, update tgbotapi.Update) error

// handleStartMessage обработчик для команды /start
func HandleStartMessage(app *Bot, update tgbotapi.Update) error {

	// Отправка нового сообщения
	_, err := app.sendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Привет!"))
	if err != nil {
		return u.HandleError(err)
	}

	app.DeleteState(update.Message.From.ID)

	return nil
}

// handleStartMessage обработчик для команды /feedback
func HandleFeedback(app *Bot, update tgbotapi.Update) error {

	chatId := int64(0)
	switch {
	case update.Message != nil:
		chatId = update.Message.From.ID
	case update.CallbackQuery != nil:
		chatId = update.CallbackQuery.From.ID
	}

	// Отправка нового сообщения
	_, err := app.sendMessage(tgbotapi.NewMessage(chatId, "Ссылка на гугл  форму для заполнения ОС по методичкам"))
	if err != nil {
		return u.HandleError(err)
	}

	return nil
}

// handleStartMessage обработчик для команды /feedback
func HandleLinks(app *Bot, update tgbotapi.Update) error {

	chatId := int64(0)
	switch {
	case update.Message != nil:
		chatId = update.Message.From.ID
	case update.CallbackQuery != nil:
		chatId = update.CallbackQuery.From.ID
	}

	// Отправка нового сообщения
	_, err := app.sendMessage(tgbotapi.NewMessage(chatId, "сообщение с ссылками"))
	if err != nil {
		return u.HandleError(err)
	}

	return nil
}

// handleStartMessage обработчик для команды /feedback
func HandleMenu(app *Bot, update tgbotapi.Update) error {

	userIsValid, err := db.DB.UserIsValid(update.Message.From.ID)
	if err != nil {
		return err
	}

	buttons := [][]ButtonData{
		{
			{Text: "Помощь", Data: "callback-help"},
			{Text: "Обратная связь", Data: "callback-feedback"},
			{Text: "Ссылки", Data: "callback-links"},
		},
	}

	if userIsValid {
		buttons = append(buttons, []ButtonData{{Text: "Выйти", Data: "callback-exit"}})
	} else {
		buttons = append(buttons, []ButtonData{{Text: "Авторизоваться", Data: "callback-auth"}})
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
	msg.ReplyMarkup = CreateInlineKeyboard(buttons)
	msg.ParseMode = "html"

	// Отправка нового сообщения
	_, err = app.sendMessage(msg)
	if err != nil {
		return u.HandleError(err)
	}

	app.DeleteState(update.Message.From.ID)

	return nil
}

func HandleHelp(app *Bot, update tgbotapi.Update) error {
	chatId := int64(0)
	switch {
	case update.Message != nil:
		chatId = update.Message.From.ID
	case update.CallbackQuery != nil:
		chatId = update.CallbackQuery.From.ID
	}

	buttons := [][]ButtonData{
		{
			{Text: "Помощь куратора", Data: "callback-tutor"},
		},
		{
			{Text: "Помощь гугл форма", Data: "callback-google-form"},
		},
	}

	msg := tgbotapi.NewMessage(chatId, "Как вам помочь?")
	msg.ReplyMarkup = CreateInlineKeyboard(buttons)
	msg.ParseMode = "html"

	// Отправка нового сообщения
	_, err := app.sendMessage(msg)
	if err != nil {
		return u.HandleError(err)
	}

	app.DeleteState(chatId)

	return nil
}

func HandleGoogleFormHelp(app *Bot, update tgbotapi.Update) error {
	chatId := update.CallbackQuery.From.ID

	msg := tgbotapi.NewMessage(chatId, "Ссылка на помощь гугл форма")
	msg.ParseMode = "html"

	// Отправка нового сообщения
	_, err := app.sendMessage(msg)
	if err != nil {
		return u.HandleError(err)
	}

	app.DeleteState(chatId)

	return nil
}

func HandleTutorHelp(app *Bot, update tgbotapi.Update) error {
	chatId := update.CallbackQuery.From.ID

	msg := tgbotapi.NewMessage(chatId, "Введите запрос, который надо отправить куратору")
	msg.ReplyMarkup = CreateKeyboard([]string{"Отмена"}, 1)
	msg.ParseMode = "html"

	// Отправка нового сообщения
	_, err := app.sendMessage(msg)
	if err != nil {
		return u.HandleError(err)
	}

	app.SetState(chatId, State{
		Name: "help-state",
	})

	return nil
}

type MessageHandler struct {
	bot    *Bot
	update tgbotapi.Update
	state  State
	userID int64
}

func HandleHelpState(app *Bot, update tgbotapi.Update) error {
	state, ok := app.GetState(update.Message.From.ID)
	if !ok {
		return fmt.Errorf("пользователь не находится в состоянии")
	}

	if update.Message.Text != "Отмена" && update.Message.Text != "Отправить" {
		state.UserMsg = update.Message.Text
	}

	// Пользователь хочет отменить отправку
	if update.Message.Text == "Отмена" {
		app.DeleteState(update.Message.From.ID)
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Заполнение заявки отменено.")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		app.sendMessage(msg)
		return nil
	}
	// Пользователь хочет отправить сообщение
	if update.Message.Text == "Отправить" {
		// Сообщение не заполнено
		if state.UserMsg == "" {
			app.sendMessage(tgbotapi.NewMessage(update.Message.From.ID, "Невозможно отправить пустое сообщение"))
			return nil
			// Сообщение заполнено. Пересылка
		} else {
			msg := tgbotapi.NewMessage(model.ConfigFile.TelegramConfig.TutorChatId, prepareNewTutorRequest(model.User{
				Username:  update.Message.From.UserName,
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
			}, state.UserMsg))
			msg.ParseMode = "html"

			_, err := app.sendMessage(msg)
			if err != nil {
				logger.Log.Error("ошибка при отправке запроса: ", err)
				app.sendMessage(tgbotapi.NewMessage(update.Message.From.ID, "Ошибка при отправке запроса: "+err.Error()))
				return fmt.Errorf("ошибка при отправке запроса: %w", err)
			}
			app.DeleteState(update.Message.From.ID)
			msg = tgbotapi.NewMessage(update.Message.From.ID, "Ваш запрос успешно отправлен.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			app.sendMessage(msg)
			return nil
		}
	}

	var buttons []string
	if state.UserMsg != "" {
		buttons = []string{"Отправить", "Отмена"}
	} else {
		buttons = []string{"Отмена"}
	}

	newState := State{
		Name:    "help-state",
		UserMsg: update.Message.Text,
	}
	app.SetState(update.Message.From.ID, newState)

	msg := tgbotapi.NewMessage(update.Message.From.ID, "Сохранен следующий запрос:\n"+update.Message.Text+"\n\nОтправьте новое сообщение, для сохранения другого запроса.")
	msg.ReplyMarkup = CreateKeyboard(buttons, 2)
	msg.ParseMode = "html"
	_, err := app.sendMessage(msg)
	if err != nil {
		logger.Log.Error("Ошибка при отправке сообщения: ", err)
		return err
	}
	return nil
}

func prepareNewTutorRequest(user model.User, request string) string {
	msgText := fmt.Sprintf(`
	<strong>Новый запрос помощи</strong>
	@%s %s %s

	%s
	`,
		user.Username, user.FirstName, user.LastName, request)

	return msgText
}

func HandleEnterExit(app *Bot, update tgbotapi.Update) error {
	chatId := int64(0)
	switch {
	case update.Message != nil:
		chatId = update.Message.From.ID
	case update.CallbackQuery != nil:
		chatId = update.CallbackQuery.From.ID
	}

	ok, err := db.DB.UserIsValid(chatId)
	if err != nil {
		return fmt.Errorf("ошибка при проверке валидности пользователя: %w", err)
	}
	if !ok {
		app.sendMessage(tgbotapi.NewMessage(chatId, "Вы не авторизованы."))
		app.DeleteState(chatId)
		return nil
	}

	buttons := [][]ButtonData{
		{
			{Text: "Нет", Data: "callback-not-exit"},
			{Text: "Нет", Data: "callback-not-exit"},
		},
		{
			{Text: "Нет", Data: "callback-not-exit"},
			{Text: "Нет", Data: "callback-not-exit"},
		},
	}

	buttons[rand.Intn(2)][rand.Intn(2)] = ButtonData{Text: "Да", Data: "callback-yes-exit"}

	msg := tgbotapi.NewMessage(chatId, "Вы уверены?")
	msg.ReplyMarkup = CreateInlineKeyboard(buttons)

	// Отправка нового сообщения
	sendedMsg, err := app.sendMessage(msg)
	if err != nil {
		return u.HandleError(err)
	}

	app.SetState(chatId, State{
		Name:  "exit",
		MsgId: sendedMsg.MessageID,
	})

	return nil
}

// handleStartMessage обработчик для команды выхода из аккаунта
func HandleExit(app *Bot, update tgbotapi.Update) error {
	chatId := update.Message.From.ID

	state, ok := app.GetState(chatId)
	if !ok {
		return fmt.Errorf("пользователь не находится в состоянии")
	}
	if state.Name != "exit" {
		return fmt.Errorf("пользователь не находится в состоянии exit")
	}

	if state.Name == "callback-not-exit" {
		msg := tgbotapi.NewDeleteMessage(chatId, state.MsgId)
		err := app.deleteMessage(msg)
		if err != nil {
			return fmt.Errorf("не удалось удалить сообщение: %w", err)
		}
		app.sendMessage(tgbotapi.NewMessage(chatId, "Выход из аккаунта отменен"))
		app.DeleteState(chatId)
		app.deleteMessage(tgbotapi.NewDeleteMessage(chatId, state.MsgId))
		return nil
	}
	if state.Name == "callback-yes-exit" {
		msg := tgbotapi.NewDeleteMessage(chatId, state.MsgId)
		err := app.deleteMessage(msg)
		if err != nil {
			return fmt.Errorf("не удалось удалить сообщение: %w", err)
		}
		db.DB.Exit(chatId)
		app.sendMessage(tgbotapi.NewMessage(chatId, "Вы успешно вышли из аккаунта"))
		app.DeleteState(chatId)

		app.deleteMessage(tgbotapi.NewDeleteMessage(chatId, state.MsgId))
		return nil
	}

	return nil
}

func HandleEnterEnterMail(app *Bot, update tgbotapi.Update) error {
	chatId := int64(0)
	switch {
	case update.Message != nil:
		chatId = update.Message.From.ID
	case update.CallbackQuery != nil:
		chatId = update.CallbackQuery.From.ID
	}

	msg := tgbotapi.NewMessage(chatId, "Введите почту")
	msg.ReplyMarkup = CreateKeyboard([]string{"Отмена"}, 1)

	app.sendMessage(msg)

	app.SetState(chatId, State{
		Name: "enter-mail",
	})

	return nil
}

// enter-mail
func HandleEnterMail(app *Bot, update tgbotapi.Update) error {
	chatId := update.Message.From.ID

	// Проверка состояния
	state, ok := app.GetState(chatId)
	if !ok {
		return fmt.Errorf("пользователь не находится в состоянии")
	}

	if state.Name != "enter-mail" {
		return fmt.Errorf("пользователь не находится в состоянии")
	}

	// Пользователь хочет отменить отправку
	if update.Message.Text == "Отмена" {
		app.DeleteState(update.Message.From.ID)
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Вход в аккаунт отменен.")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		app.sendMessage(msg)
		return nil
	}

	// Сохраняем данные почты, если это не команда
	if update.Message.Text != "Отмена" && update.Message.Text != "Отправить код подтверждения" {
		err := db.DB.SetNewMail(strconv.FormatInt(chatId, 10), update.Message.Text)
		if err != nil {
			logger.Log.Error("Не удалось сохранить почту в БД: ", err)
			app.sendMessage(tgbotapi.NewMessage(chatId, "Не удалось сохранить почту: "+err.Error()))
			return fmt.Errorf("не удалось сохранить почту в БД: %w", err)
		}

		buttons := []string{"Отправить код подтверждения", "Отмена"}
		msg := tgbotapi.NewMessage(chatId, "Отправить код подтверждения на почту?\n"+update.Message.Text+"\nЕсли адрес неверный, просто введите сообщение с новым")
		msg.ReplyMarkup = CreateKeyboard(buttons, 1)
		msg.ParseMode = "html"
		_, err = app.sendMessage(msg)
		if err != nil {
			logger.Log.Error("Ошибка при отправке сообщения: ", err)
			return err
		}
	}

	// Пользователь хочет отправить сообщение
	if update.Message.Text == "Отправить код подтверждения" {
		mail, err, _ := db.DB.GetMail(chatId)
		if err != nil {
			logger.Log.Error("Ошибка при получении почты: %w", err)
			app.sendMessage(tgbotapi.NewMessage(chatId, "Ошибка при получении почты: "+err.Error()))
			return fmt.Errorf("ошибка при получении почты: %w", err)
		}

		code := utils.GenerateCode()

		err = smtp.SendAuthorizationCode(mail, code)
		if err != nil {
			logger.Log.Error("ошибка при отправке кода на почту: %w", err)
			app.sendMessage(tgbotapi.NewMessage(chatId, "Ошибка при получении почты: "+err.Error()))
			return fmt.Errorf("ошибка при отправке кода на почту: %w", err)
		}

		err = db.DB.SetNewCode(strconv.FormatInt(chatId, 10), code)
		if err != nil {
			logger.Log.Error("ошибка при сохранении кода: %w", err)
			app.sendMessage(tgbotapi.NewMessage(chatId, "ошибка при сохранении кода: "+err.Error()))
			return fmt.Errorf("ошибка при сохранении кода: %w", err)
		}

		buttons := []string{"Отмена"}
		msg := tgbotapi.NewMessage(chatId, "Код подтверждения отправлен на почту. Введите его.")
		msg.ReplyMarkup = CreateKeyboard(buttons, 1)
		msg.ParseMode = "html"
		_, err = app.sendMessage(msg)
		if err != nil {
			logger.Log.Error("Ошибка при отправке сообщения: ", err)
			return err
		}

		app.SetState(chatId, State{
			Name: "enter-code",
		})
		return nil
	}

	return nil
}

func HandleCancel(app *Bot, update tgbotapi.Update) error {
	_, ok := app.GetState(update.Message.From.ID)
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Меню закрыто")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		app.sendMessage(msg)
	}

	return nil
}

// enter-code
func HandleEnterCode(app *Bot, update tgbotapi.Update) error {
	chatId := update.Message.From.ID

	// Проверка состояния
	state, ok := app.GetState(chatId)
	if !ok {
		return fmt.Errorf("пользователь не находится в состоянии")
	}

	if state.Name != "enter-code" {
		return fmt.Errorf("пользователь не находится в состоянии кода")
	}

	// Пользователь хочет отменить отправку
	if update.Message.Text == "Отмена" {
		app.DeleteState(update.Message.From.ID)
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Вход в аккаунт отменен.")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		app.sendMessage(msg)
		return nil
	}

	// Сохраняем код в оперативку, если не команда, и уведомляет пользователя
	if update.Message.Text != "Отмена" && update.Message.Text != "Подтвердить" {
		state.UserMsg = update.Message.Text
		app.SetState(chatId, state)

		buttons := []string{"Подтвердить", "Отмена"}
		msg := tgbotapi.NewMessage(chatId, "Ваш код подтверждения: "+update.Message.Text+" ?")
		msg.ReplyMarkup = CreateKeyboard(buttons, 1)
		msg.ParseMode = "html"
		_, err := app.sendMessage(msg)
		if err != nil {
			logger.Log.Error("Ошибка при отправке сообщения: ", err)
			return err
		}
	}

	// Пользователь хочет отправить сообщение
	if update.Message.Text == "Подтвердить" {
		err := db.DB.CheckCode(strconv.FormatInt(chatId, 10), state.UserMsg)
		if err != nil {
			logger.Log.Error("не удалось подтвердить код: %w", err)
			app.sendMessage(tgbotapi.NewMessage(chatId, "Не удалось подтвердить код: "+err.Error()))
			return fmt.Errorf("не удалось подтвердить код: %w", err)
		}

		app.DeleteState(chatId)

		app.DeleteState(update.Message.From.ID)
		msg := tgbotapi.NewMessage(update.Message.From.ID, "Вход в аккаунт выполнен!")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		app.sendMessage(msg)

		return nil
	}

	return nil
}
