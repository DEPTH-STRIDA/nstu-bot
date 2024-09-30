package tg

import (
	"errors"
	"fmt"
	"nstu/internal/logger"
	"regexp"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseDate(input string) (time.Time, error) {
	if len(input) > 15 {
		return time.Time{}, errors.New("входная строка слишком длинная")
	}

	// Извлекаем только цифры
	reg := regexp.MustCompile(`\d+`)
	numbers := reg.FindAllString(input, -1)

	if len(numbers) < 2 || len(numbers) > 3 {
		return time.Time{}, errors.New("неверный формат даты")
	}

	day := numbers[0]
	month := numbers[1]
	year := ""

	if len(numbers) == 3 {
		year = numbers[2]
	}

	// Проверяем, что день и месяц двузначные
	if len(day) != 2 || len(month) != 2 {
		return time.Time{}, errors.New("день и месяц должны быть двузначными")
	}

	currentYear := time.Now().Year()

	// Если год не указан или указан не полностью, используем текущий
	if year == "" || len(year) != 4 {
		year = fmt.Sprintf("%04d", currentYear)
	}

	// Собираем дату в формате "02-01-2006"
	dateStr := fmt.Sprintf("%s-%s-%s", day, month, year)

	// Парсим дату
	date, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return time.Time{}, errors.New("неверная дата")
	}

	// Если указан только день и месяц, и получившаяся дата в прошлом,
	// добавляем один год
	if len(numbers) == 2 && date.Before(time.Now()) {
		date = date.AddDate(1, 0, 0)
	}

	return date, nil
}

// SendMessageRepet делает несколько попыток отправки сообщений.
// Останавливает попытки после первой успешной.
func (app *BotStruct) SendMessageRepet(msg tgbotapi.MessageConfig, numberRepetion int) (tgbotapi.Message, error) {
	for i := 0; i < numberRepetion; i++ {
		sendedMsg, err := app.SendMessage(msg)
		if err != nil {
			logger.Log.Info("Ошибка при отправке сообщения с повтором (", i, "):  ", err)
		} else {
			return sendedMsg, nil
		}
	}
	return tgbotapi.Message{}, fmt.Errorf("ни одна попытка не оказалось результативной")
}

// SendMessage синхронная функция для отправки сообщения
func (app *BotStruct) SendMessage(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	sendedMsg, err := app.sendMessage(msg)
	if err != nil {
		return sendedMsg, err
	}
	return sendedMsg, nil
}

// sendMessage асинхронная функция, которая с помощью waitgroup дожидается результатов от отправки сообщения
func (app *BotStruct) sendMessage(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	var sendedMsg tgbotapi.Message
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	// Отправляем функцию в канал
	app.msgRequestHandler.HandleRequest(func() error {
		defer wg.Done()

		// Устанавливаем глобальные параметры
		sendedMsg, err = app.Send(msg)
		if err != nil {
			return err
		}

		return nil
	})

	wg.Wait()
	return sendedMsg, err
}

func (app *BotStruct) EditMessageRepet(editMsg tgbotapi.EditMessageTextConfig, numberRepetion int) (*tgbotapi.APIResponse, error) {
	for i := 0; i < numberRepetion; i++ {
		response, err := app.editMessage(editMsg)
		if err != nil {
			logger.Log.Info("Ошибка при редактировании сообщения с повтором (", i, "):  ", err)
		} else {
			return response, nil
		}
	}
	return nil, fmt.Errorf("ни одна попытка не стала результативной")
}

// EditMessage синхронно редактирует сообщение
func (app *BotStruct) EditMessage(editMsg tgbotapi.EditMessageTextConfig) (*tgbotapi.APIResponse, error) {
	response, err := app.editMessage(editMsg)
	if err != nil {
		return response, err
	}

	return response, nil
}

// editMessage редактирует сообщение в чате, отправив функцию редактирования в запросы
func (app *BotStruct) editMessage(editMsg tgbotapi.EditMessageTextConfig) (*tgbotapi.APIResponse, error) {
	var response *tgbotapi.APIResponse
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	// Отправляем функцию в канал
	app.msgRequestHandler.HandleRequest(func() error {
		defer wg.Done()

		// Устанавливаем глобальные параметры
		response, err = app.Request(editMsg)
		if err != nil {
			return err
		}

		return nil
	})

	wg.Wait()
	return response, err
}

func (app *BotStruct) DeleteMessageRepet(msgToDelete tgbotapi.DeleteMessageConfig, numberRepetion int) error {
	for i := 0; i < numberRepetion; i++ {
		err := app.deleteMessage(msgToDelete)
		if err != nil {
			logger.Log.Info("Не удалось удалить сообщение из чата. Попытка: ", i, " err: ", err)
		} else {
			return nil
		}
	}

	return fmt.Errorf("ни одна попытка не стала результативной")
}

// DeleteMessage удаляет сообщение
func (app *BotStruct) DeleteMessage(msgToDelete tgbotapi.DeleteMessageConfig) error {
	err := app.deleteMessage(msgToDelete)
	if err != nil {
		return err
	}

	return nil
}

func (app *BotStruct) deleteMessage(deleteMsg tgbotapi.DeleteMessageConfig) error {
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	// Отправляем функцию в канал
	app.msgRequestHandler.HandleRequest(func() error {
		defer wg.Done()

		_, err = app.Request(deleteMsg)
		if err != nil {
			return err
		}

		return nil
	})

	wg.Wait()
	return err
}

// ShowAlert показывает пользователю предупреждение. alert по типу браузерного.
// Для закрытия такого уведомления потребуется нажать "ок"
func (app *BotStruct) ShowAlert(CallbackQueryID string, alertText string) {
	callback := tgbotapi.NewCallback(CallbackQueryID, alertText)
	// Это заставит текст появиться во всплывающем окне
	callback.ShowAlert = true
	_, err := app.BotAPI.Request(callback)
	if err != nil {
		logger.Log.Info("Не удалось показать alert после CallbackQuery: ", err)
	}
}

func CreateKeyboard(input []string, buttonsPerRow int) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	for i := 0; i < len(input); i += buttonsPerRow {
		var row []tgbotapi.KeyboardButton
		end := i + buttonsPerRow
		if end > len(input) {
			end = len(input)
		}
		for _, text := range input[i:end] {
			row = append(row, tgbotapi.NewKeyboardButton(text))
		}
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewReplyKeyboard(keyboard...)
}

type ButtonData struct {
	Text string
	Data string
}

//	buttons := [][]telegram.ButtonData{
//		{
//			{Text: "1.com", Data: "http://1.com"},
//			{Text: "2", Data: "2"},
//			{Text: "3", Data: "3"},
//		},
//		{
//			{Text: "4", Data: "4"},
//			{Text: "5", Data: "5"},
//			{Text: "6", Data: "6"},
//		},
//	}
func CreateInlineKeyboard(buttons [][]ButtonData) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, row := range buttons {
		var keyboardRow []tgbotapi.InlineKeyboardButton
		for _, btn := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Data))
		}
		keyboard = append(keyboard, keyboardRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func CreateReplyKeyboard(buttons [][]ButtonData) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	for _, row := range buttons {
		var keyboardRow []tgbotapi.KeyboardButton
		for _, button := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewKeyboardButton(button.Text))
		}
		keyboard = append(keyboard, keyboardRow)
	}

	return tgbotapi.NewReplyKeyboard(keyboard...)
}

func (app *BotStruct) SendMessageRepetLowPriority(msg tgbotapi.MessageConfig, numberRepetion int) (tgbotapi.Message, error) {
	for i := 0; i < numberRepetion; i++ {
		sendedMsg, err := app.SendMessageLowPriority(msg)
		if err != nil {
			logger.Log.Info("Ошибка при отправке сообщения с повтором (", i, "):  ", err)
		} else {
			return sendedMsg, nil
		}
	}
	return tgbotapi.Message{}, nil
}

// SendMessageLowPriority синхронная функция, которая отправляет сообщение в телеграм с низким приоритетом
func (app *BotStruct) SendMessageLowPriority(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	sendedMsg, err := app.sendMessageLowPriority(msg)
	if err != nil {
		return sendedMsg, err
	}
	return sendedMsg, nil
}

// sendMessage асинхронная функция, которая с помощью waitgroup дожидается результатов от отправки сообщения
func (app *BotStruct) sendMessageLowPriority(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	var sendedMsg tgbotapi.Message
	var err error

	var wg sync.WaitGroup
	wg.Add(1)

	// Отправляем функцию в канал
	app.msgRequestHandler.HandleLowPriorityRequest(func() error {
		defer wg.Done()

		// Устанавливаем глобальные параметры
		sendedMsg, err = app.Send(msg)
		if err != nil {
			return err
		}
		return nil

	})

	wg.Wait()
	return sendedMsg, err
}

func (app *BotStruct) EditMessageRepetLowPriority(editMsg tgbotapi.EditMessageTextConfig, numberRepetion int) (*tgbotapi.APIResponse, error) {
	for i := 0; i < numberRepetion; i++ {
		response, err := app.EditMessageLowPriority(editMsg)
		if err != nil {
			logger.Log.Info("Ошибка при редактировании сообщения с повтором (", i, "):  ", err)
		} else {
			return response, nil
		}
	}
	return nil, nil
}

// EditMessageLowPriority синхронная функция дли редактирования вообщения
func (app *BotStruct) EditMessageLowPriority(editMsg tgbotapi.EditMessageTextConfig) (*tgbotapi.APIResponse, error) {
	response, err := app.editMessageLowPriority(editMsg)
	if err != nil {
		return response, err
	}

	return response, nil
}

// editMessage редактирует сообщение в чате, отправив функцию редактирования в запросы
func (app *BotStruct) editMessageLowPriority(editMsg tgbotapi.EditMessageTextConfig) (*tgbotapi.APIResponse, error) {
	var err error
	var editedMsg *tgbotapi.APIResponse

	var wg sync.WaitGroup
	wg.Add(1)

	// Отправляем функцию в канал
	app.msgRequestHandler.HandleLowPriorityRequest(func() error {
		defer wg.Done()

		editedMsg, err = app.Request(editMsg)
		if err != nil {
			return err
		}
		return nil
	})

	wg.Wait()
	return editedMsg, err
}

func (app *BotStruct) SendMessageButtonLowPriorityRepet(chatID int64, msgText, buttonText, buttonCallbackText string, numberRepetion int) (tgbotapi.Message, error) {
	for i := 0; i < numberRepetion; i++ {
		sendedMsg, err := app.SendMessageButtonLowPriority(chatID, msgText, buttonText, buttonCallbackText)
		if err != nil {
			logger.Log.Info("Ошибка при отправке сообщения: ", err)
		} else {
			return sendedMsg, err
		}
	}

	return tgbotapi.Message{}, nil
}

// SendMessage синхронная функция, которая отправляет сообщение с кнопкой
func (app *BotStruct) SendMessageButtonLowPriority(chatID int64, msgText, buttonText, buttonCallbackText string) (tgbotapi.Message, error) {

	msg := tgbotapi.NewMessage(chatID, msgText)
	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCallbackText),
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row)
	msg.ParseMode = "html"
	sendedMsg, err := app.SendMessageLowPriority(msg)
	if err != nil {
		return sendedMsg, err
	}
	return sendedMsg, nil
}
