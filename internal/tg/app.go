package tg

import (
	"errors"
	"fmt"
	"nstu/internal/logger"
	"nstu/internal/models"
	"nstu/internal/request"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var Bot *BotStruct

func InitTgBNot() (*BotStruct, error) {
	bot, err := NewBot()
	return bot, err
}

type Request func() error

type BotStruct struct {
	*tgbotapi.BotAPI
	msgRequestHandler      *request.RequestHandler
	callbackRequestHandler *request.RequestHandler
	MsgRoutes              map[string]HandlerFunc
	CallbackRoutes         map[string]HandlerFunc
}

// Конструктор нового бота
func NewBot() (*BotStruct, error) {
	msgRequestHandler, err := request.NewRequestHandler(int64(models.Config.TelegramConfig.MessageBufferSize))
	if err != nil {
		return nil, err
	}
	callbackRequestHandler, err := request.NewRequestHandler(int64(models.Config.TelegramConfig.CallBackBufferSize))
	if err != nil {
		return nil, err
	}
	app := BotStruct{
		msgRequestHandler:      msgRequestHandler,
		callbackRequestHandler: callbackRequestHandler,
	}

	go app.msgRequestHandler.ProcessRequests(time.Duration(models.Config.TelegramConfig.MessagePauseSec) * time.Second)
	go app.callbackRequestHandler.ProcessRequests(time.Duration(models.Config.TelegramConfig.CallBackPauseSec) * time.Second)

	app.BotAPI, err = tgbotapi.NewBotAPI(models.Config.TelegramConfig.NSTUToken)
	if err != nil {
		return nil, fmt.Errorf("не удается инициализировать бота telegram: %v", err)
	}

	app.MsgRoutes = make(map[string]HandlerFunc)
	app.CallbackRoutes = make(map[string]HandlerFunc)

	err = app.SetMsgRoutes()
	if err != nil {
		return nil, errors.New("установка обработчиков сообщений не удалась: " + err.Error())
	}
	err = app.SetCallBackRoutes()
	if err != nil {
		return nil, errors.New("установка обработчиков callback событий не удалась: " + err.Error())
	}
	return &app, nil
}

// HandleUpdates запускает обработку всех обновлений поступающих боту из телеграмма
func (app *BotStruct) HandleUpdates() {
	// Настройка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	// Получение канала обновлений
	updates := app.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			go app.handleCallback(update)
		case update.Message != nil:
			go app.handleMessage(update)
		default:
			logger.Log.Error("Боту поступило неизвестное обновление")
		}

	}
}

// handleMessage ищет команду в map'е и выполняет ее
func (app *BotStruct) handleMessage(update tgbotapi.Update) {
	if update.Message.From.IsBot {
		logger.Log.Info("Игнорирование сообщений от бота: ", update.Message.From)
		return
	}

	err := models.EnsureUserExists(models.DataBase.DB, initdata.InitData{User: initdata.User{
		ID:        update.Message.From.ID,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Username:  update.Message.From.UserName,
	}})
	if err != nil {
		logger.Log.Error("Не удалось проверить наличие пользвователя в БД: ", err)
		msgText := "Внутренняя ошибка: " + err.Error()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

		_, err := app.SendMessage(msg)
		if err != nil {
			logger.Log.Error("Ошибка при отправке сообщения: ", err.Error())
		}
		return
	}

	currentAction, ok := app.MsgRoutes[update.Message.Text]
	if ok {
		if currentAction == nil {
			logger.Log.Error("Обработчик команды " + update.Message.Text + " равен nil")
			return
		}

		err := currentAction(app, update)
		if err != nil {
			logger.Log.Error("Ошибка при обработки команды ", update.Message.Text, " от пользователя (", update.Message.Chat.ID, ":", update.Message.Chat.UserName)
		} else {
			logger.Log.Info("Успешно обработана команда: ", update.Message.Text, " от пользователя (", update.Message.Chat.ID, ":", update.Message.Chat.UserName)
		}
	} else {
		logger.Log.Info("Пользователь ( ", update.Message.Chat.ID, " : ", update.Message.Chat.UserName, " отправил команду ", update.Message.Text, ": в чат, которая не была найдена.")

		date, err := ParseDate(update.Message.Text)
		if err != nil {
			logger.Log.Info("Команда "+update.Message.Text+" не является датой: ", err)
			return
		}

		action := HandleSchedule(date)
		action(app, update)
	}
}

// handleCallback ищет команду в map'е и выполняет ее
func (app *BotStruct) handleCallback(update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}
	// Вызов стандартных функций из БД
	currentAction, ok := app.CallbackRoutes[update.CallbackQuery.Data]
	if !ok {
		if currentAction == nil {
			logger.Log.Error("Обработчик callback команды " + fmt.Sprint(update.CallbackQuery) + " равен nil")
			return
		}

		//  Запуск события на обработку реакции
		currentActionTemp, ok := app.CallbackRoutes["callbackDB"]
		if !ok {
			// Если метода для обработки реакции нет, значит явно есть ошибка
			logger.Log.Info("Неизвестная Callback команда")
			return
		} else {
			currentAction = currentActionTemp
		}
	}
	err := currentAction(app, update)
	if err != nil {
		logger.Log.Error("Ошибка при обработки Callback команды от пользователя (", update.CallbackQuery.From.ID, ":", update.CallbackQuery.From.UserName)
	} else {
		logger.Log.Info("Успешно обработана Callback команда: ", update.CallbackQuery.Data, " от пользователя (", update.CallbackQuery.From.ID, ":", update.CallbackQuery.From.UserName)
	}
}

// HandlRoute добавляет обработку комманды
func (app *BotStruct) HandlMsgRoute(command string, handler HandlerFunc) {
	app.MsgRoutes[command] = handler
}

// HandlRoute добавляет обработку комманды
func (app *BotStruct) HandlCallbackRoute(command string, handler HandlerFunc) {
	app.CallbackRoutes[command] = handler
	logger.Log.Info("Успешно установлена callback команда: ", command)
}

func (app *BotStruct) DeleteCallbackRoute(command string) {
	delete(app.CallbackRoutes, "command")
}

func (app *BotStruct) DeleteMsgRoute(command string) {
	delete(app.MsgRoutes, "command")
}
