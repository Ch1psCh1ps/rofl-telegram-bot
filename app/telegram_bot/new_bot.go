package telegram_bot

//func StartNewBot() {
//	bot, err := tgbotapi.NewBotAPI("6094698233:AAFjy8J16Pd2Uv90W9VJ-UNqr9gZaBrWeqU")
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//
//	updates := bot.GetUpdatesChan(u)
//
//	myPhotoPath := "/Users/anastasyaplotnikova/Desktop/gachi.jpeg"
//
//	for update := range updates {
//		if update.Message != nil {
//			handleMessage(bot, update.Message, myPhotoPath)
//		}
//	}
//}
//
//func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, photoPath string) {
//	switch message.Text {
//	case "/start":
//		handleStartCommand(bot, message)
//	case "Рефактор xlsx":
//		handleRefactorCommand(bot, message)
//	default:
//		handleDefaultMessage(bot, message, photoPath)
//	}
//}
//
//func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	buttonMsg := tgbotapi.NewMessage(message.Chat.ID, "Чем могу помочь?")
//
//	// Создание кнопки
//	buttonText := "Рефактор xlsx"
//	button := tgbotapi.NewKeyboardButton(buttonText)
//
//	// Создание клавиатуры и добавление кнопки
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(button),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	buttonMsg.ReplyMarkup = keyboard
//
//	// Отправка сообщения с кнопкой
//	bot.Send(buttonMsg)
//}
//
//func handleRefactorCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Пришлите файл")
//	bot.Send(replyMsg)
//}
//
//func handleDefaultMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, photoPath string) {
//	if message.Document != nil {
//		if isXLSXFile(message.Document.FileName) {
//			fileID := message.Document.FileID
//			fileURL, err := bot.GetFileDirectURL(fileID)
//			if err != nil {
//				log.Printf("Ошибка при получении файла: %v", err)
//				sendErrorMessage(bot, message.Chat.ID)
//				return
//			}
//
//			sendProcessingMessage(bot, message.Chat.ID)
//
//			xlsxBuffer, err := refactor_xlsx.DoBookXlsx(fileURL)
//			if err != nil {
//				log.Printf("Ошибка при обработке файла: %v", err)
//				sendErrorMessage(bot, message.Chat.ID)
//				return
//			}
//
//			sendUpdateMessage(bot, message.Chat.ID)
//			sendXLSXFile(bot, message.Chat.ID, xlsxBuffer)
//
//			sendPhotoMessage(bot, message.Chat.ID, photoPath)
//
//		} else {
//			sendInvalidFileMessage(bot, message.Chat.ID)
//		}
//	} else {
//		sendDefaultMessage(bot, message.Chat.ID)
//	}
//}
//
//func sendProcessingMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Сейчас сделаю обновленный файл 🥰")
//	bot.Send(msg)
//}
//
//func sendUpdateMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Вот обновление 🤩")
//	bot.Send(msg)
//}
//
//func sendXLSXFile(bot *tgbotapi.BotAPI, chatID int64, xlsxBuffer *bytes.Buffer) {
//	xlsxConfig := tgbotapi.FileBytes{
//		Name:  "refactor_available_units.xlsx",
//		Bytes: xlsxBuffer.Bytes(),
//	}
//
//	docMsg := tgbotapi.NewDocument(chatID, xlsxConfig)
//	bot.Send(docMsg)
//}
//
//func sendPhotoMessage(bot *tgbotapi.BotAPI, chatID int64, photoPath string) {
//	photoBytes, err := ioutil.ReadFile(photoPath)
//	if err == nil {
//		msg := tgbotapi.NewMessage(chatID, "By the way: nice ass, Dude 😏\nAwesome balls, bro 😎")
//		bot.Send(msg)
//
//		photo := tgbotapi.FileBytes{
//			Name:  photoPath,
//			Bytes: photoBytes,
//		}
//
//		photoMsg := tgbotapi.NewPhoto(chatID, photo)
//		photoMsg.Caption = "Ну что? Разомнемся?"
//		bot.Send(photoMsg)
//	} else {
//		log.Printf("Failed to read photo: %v", err)
//	}
//}
//
//func sendInvalidFileMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Присланный документ не имеет расширение .xlsx 🤕\nПожалуйста, пришлите файл в формате .xlsx")
//	bot.Send(msg)
//}
//
//func sendErrorMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Упс... 🙊 Что-то пошло не так 😓\nПроизошла ошибка при обработке файла 📁 Уже бежим исправлять 🛠")
//	bot.Send(msg)
//}
//
//func sendDefaultMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Я могу принимать файлы, но не могу с вами общаться 😔\nМогу вам помочь с рефактором вашего файла? 🥺")
//	bot.Send(msg)
//}
