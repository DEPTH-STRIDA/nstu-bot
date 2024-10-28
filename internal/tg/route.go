package tg

/*
menu - посмотреть меню
help -помощь куратора или через форму
feedback - оставить ОС по методичкам
links - посмотреть ссылки
auth - авторизоваться
*/

func (b *Bot) SetMsgRoutes() error {

	b.HandlMsgRoute("/start", HandleStartMessage)

	b.HandlMsgRoute("/help", HandleHelp)
	b.HandlMsgRoute("/menu", HandleMenu)
	b.HandlMsgRoute("/feedback", HandleFeedback)
	b.HandlMsgRoute("/links", HandleLinks)
	b.HandlMsgRoute("/menu", HandleMenu)

	b.HandlMsgRoute("отмена", HandleCancel)

	b.HandlMsgRoute("выйти", HandleExit)

	b.HandlMsgRoute("/auth", HandleEnterEnterMail)

	return nil
}

func (b *Bot) SetCallBackRoutes() error {
	b.HandlCallbackRoute("callback-help", HandleHelp)
	b.HandlCallbackRoute("callback-feedback", HandleFeedback)
	b.HandlCallbackRoute("callback-links", HandleLinks)
	b.HandlCallbackRoute("callback-google-form", HandleGoogleFormHelp)
	b.HandlCallbackRoute("callback-tutor", HandleTutorHelp)

	b.HandlCallbackRoute("callback-auth", HandleEnterEnterMail)
	return nil
}
