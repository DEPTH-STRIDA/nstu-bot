package tg

import "time"

func (b *BotStruct) SetMsgRoutes() error {

	b.HandlMsgRoute("/start", HandleStartMessage())
	b.HandlMsgRoute("/help", HandleStartMessage())

	b.HandlMsgRoute("/open", HandleOpen())
	b.HandlMsgRoute("/close", HandleClose())
	b.HandlMsgRoute("/menu", HandleStartMessage())

	b.HandlMsgRoute("Сегодня", HandleSchedule(time.Now()))
	b.HandlMsgRoute("Завтра", HandleSchedule(time.Now().Add(24*time.Hour)))
	b.HandlMsgRoute("Все расписание", HandleAllSchedule())

	b.HandlMsgRoute("Войти в меню", HandleEnterMenu())

	return nil
}

func (b *BotStruct) SetCallBackRoutes() error {

	// b.HandlMsgRoute("/start", HandleStartMessage())

	return nil
}
