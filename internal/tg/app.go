package tg

import (
	"app/internal/db"
	"app/internal/logger"
	"app/internal/model"
	"app/internal/request"
	"app/internal/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
)

var TelegramBot *Bot

type Request func() error

type State struct {
	Name    string
	UserMsg string
	MsgId   int
}

type Bot struct {
	botAPI *tgbotapi.BotAPI

	HandlCallbackLessonButton func(app *Bot, update tgbotapi.Update) error

	msgRequestHandler      *request.RequestHandler
	callbackRequestHandler *request.RequestHandler
	MsgRoutes              map[string]HandlerFunc
	CallbackRoutes         map[string]HandlerFunc

	userStates *cache.Cache // Хранит состояния пользователей
}

// Конструктор нового бота
func InitTgBot() error {
	conf := model.ConfigFile.TelegramConfig

	msgRequestHandler, err := request.NewRequestHandler(conf.MsgBufferSize)
	if err != nil {
		return err
	}
	callbackRequestHandler, err := request.NewRequestHandler(conf.CallBackBufferSize)
	if err != nil {
		return err
	}
	app := Bot{
		msgRequestHandler:      msgRequestHandler,
		callbackRequestHandler: callbackRequestHandler,
		userStates:             cache.New(1*time.Hour, 10*time.Minute),
	}

	go app.msgRequestHandler.ProcessRequests(time.Duration(conf.RequestUpdatePause) * time.Millisecond)
	go app.callbackRequestHandler.ProcessRequests(time.Duration(conf.RequestCallBackUpdatePause) * time.Millisecond)

	app.botAPI, err = tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		return fmt.Errorf("не удается инициализировать бота telegram: %v", err)
	}

	app.MsgRoutes = make(map[string]HandlerFunc)
	app.CallbackRoutes = make(map[string]HandlerFunc)

	err = app.SetMsgRoutes()
	if err != nil {
		return errors.New("установка обработчиков сообщений не удалась: " + err.Error())
	}
	err = app.SetCallBackRoutes()
	if err != nil {
		return errors.New("установка обработчиков callback событий не удалась: " + err.Error())
	}

	go app.HandleUpdates()

	TelegramBot = &app
	return nil
}

// HandleUpdates запускает обработку всех обновлений поступающих боту из телеграмма
func (app *Bot) HandleUpdates() {
	// Настройка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	// Получение канала обновлений
	updates := app.botAPI.GetUpdatesChan(u)
	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			go app.handleCallback(update)
		case update.Message != nil:
			go app.handleMessage(update)
		default:
			logger.Log.Info("Боту поступило неизвестное обновление от: ", update.SentFrom())
		}

	}
}

// handleMessage ищет команду в map'е и выполняет ее
func (app *Bot) handleMessage(update tgbotapi.Update) {
	if update.Message.Chat.ID == model.ConfigFile.TutorChatId {
		logger.Log.Info("Игнорирование команд из чата кураторов")
		return
	}

	// Добавление пользователя в БД
	db.DB.CreateUserIfNotExistsWithConditions(model.User{
		TelegramID: update.Message.From.ID,
		Username:   update.Message.From.UserName,
		FirstName:  update.Message.From.FirstName,
		LastName:   update.Message.From.LastName,
	})

	fondHandler := false

	// Поиск события в map'е. Если событие приведет к выходу из состояния, то отлично
	if currentAction, ok := app.MsgRoutes[strings.ToLower(update.Message.Text)]; ok {
		fondHandler = true
		if err := currentAction(app, update); err != nil {
			logger.Log.Error("Ошибка при обработки команды ", update.Message.Text, " от пользователя (", update.Message.Chat.ID, ":", update.Message.Chat.UserName)
		} else {
			logger.Log.Info("Успешно обработана команда: ", update.Message.Text, " от пользователя (", update.Message.Chat.ID, ":", update.Message.Chat.UserName)
		}
	} else {
		logger.Log.Info("Пользователь ( ", update.Message.Chat.ID, " : ", update.Message.Chat.UserName, " отправил команду ", update.Message.Text, ": в чат, которая не была найдена.")
	}

	// Реагируем как пооложено, если пользователь находится в состоянии
	state, ok := app.GetState(update.Message.From.ID)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!", state)
	if ok {
		switch state.Name {
		case "help-state":
			fondHandler = true
			utils.HandleError(HandleHelpState(app, update))

		case "enter-mail":
			fondHandler = true
			utils.HandleError(HandleEnterMail(app, update))

		case "enter-code":
			fondHandler = true
			utils.HandleError(HandleEnterCode(app, update))
		}
	}

	if !fondHandler {
		app.sendMessage(tgbotapi.NewMessage(update.Message.From.ID, "Неизвестная команда. Вы можете вызвать /menu или нажать на кнопку \"Меню\" слева от клавиатуры, чтобы увидеть доступные команды."))
	}
}

// handleCallback ищет команду в map'е и выполняет ее
func (app *Bot) handleCallback(update tgbotapi.Update) {
	if update.CallbackQuery.Message.Chat.ID == model.ConfigFile.TutorChatId {
		logger.Log.Info("Игнорирование команд из чата кураторов")
		return
	}

	db.DB.CreateUserIfNotExistsWithConditions(model.User{
		TelegramID: update.CallbackQuery.From.ID,
		Username:   update.CallbackQuery.From.UserName,
		FirstName:  update.CallbackQuery.From.FirstName,
		LastName:   update.CallbackQuery.From.LastName,
	})

	// Поиск обработчика в map'е callback обработчиков
	if currentAction, ok := app.CallbackRoutes[update.CallbackQuery.Data]; ok {
		// Обработчик найден
		if err := currentAction(app, update); err != nil {
			logger.Log.Error("Ошибка при обработки Callback команды от пользователя (", update.CallbackQuery.From.ID, ":", update.CallbackQuery.From.UserName, err)
		} else {
			logger.Log.Info("Успешно обработана Callback команда: ", update.CallbackQuery.Data, " от пользователя (", update.CallbackQuery.From.ID, ":", update.CallbackQuery.From.UserName)
		}
	} else {

		if app.HandlCallbackLessonButton != nil {
			app.HandlCallbackLessonButton(app, update)
		} else {
			app.ShowAlert(update.CallbackQuery.ID, "Внутренняя ошибка. Пожалуйста подождите.")
		}

		logger.Log.Info("Вызван callback метод: ", update.CallbackQuery.Data, "-  для которого не установлен обработчик. От пользователя (", update.CallbackQuery.From.ID, ":", update.CallbackQuery.From.UserName)
	}

	// Реагируем как пооложено, если пользователь находится в состоянии
	state, ok := app.GetState(update.CallbackQuery.From.ID)
	if ok {
		switch state.Name {
		case "callback-yes-exit":
			HandleExit(app, update)
		case "callback-not-exit":
			HandleExit(app, update)
		}
	}

	// Отвечаем на callback query
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	app.botAPI.Request(callback)
}

// HandlRoute добавляет обработку комманды
func (app *Bot) HandlMsgRoute(command string, handler HandlerFunc) {
	app.MsgRoutes[command] = handler
	logger.Log.Info("Успешно установлена текстовая команда: ", command)
}

// HandlRoute добавляет обработку комманды
func (app *Bot) HandlCallbackRoute(command string, handler HandlerFunc) {
	app.CallbackRoutes[command] = handler
	logger.Log.Info("Успешно установлена callback команда: ", command)
}

func (app *Bot) DeleteCallbackRoute(command string) {
	delete(app.CallbackRoutes, "command")
}

func (app *Bot) DeleteMsgRoute(command string) {
	delete(app.MsgRoutes, "command")
}

// Сохранение в кеш
func (app *Bot) SetState(telegramID int64, stateName State) {
	tgId := strconv.FormatInt(telegramID, 10)

	app.userStates.Set(tgId, stateName, cache.DefaultExpiration)
}

// Получение из кеша
func (app *Bot) GetState(telegramID int64) (State, bool) {
	tgId := strconv.FormatInt(telegramID, 10)

	if value, found := app.userStates.Get(tgId); found {
		if state, ok := value.(State); ok {
			return state, true
		}
	}
	return State{}, false
}

// Удаление из кеша
func (app *Bot) DeleteState(telegramID int64) {
	tgId := strconv.FormatInt(telegramID, 10)

	app.userStates.Delete(tgId)
}
