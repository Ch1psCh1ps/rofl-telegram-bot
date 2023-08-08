package telegram_bot

//
//import (
//	"bytes"
//	"fmt"
//	"genieMap/cmd"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"log"
//	"os"
//	"strconv"
//	"strings"
//)
//
//var serviceNum int
//
//type BotState struct {
//	ServiceNum int
//}
//
//func getServices() []string {
//	services := []string{
//		"Al Dar",
//		"Town X",
//		"Condor",
//		"Siadah",
//		"Binghatii",
//		"Deyaar",
//		"Emaar",
//		"Ellington Properties",
//		"Azizi",
//		"Reportage Properties",
//		"Condor 8 колонок",
//		"Tiger",
//		"Sothoby's",
//		"Object 1",
//		"Swiss",
//		"Vincitore",
//	}
//
//	return services
//}
//
//func getServiceName(serviceNumber int) string {
//	services := getServices()
//	return services[serviceNumber-1]
//}
//
//func isValidService(serviceNumber int) bool {
//	services := getServices()
//	return serviceNumber >= 1 && serviceNumber <= len(services)
//}
//
//func getServiceFileFormat(serviceNumber int) string {
//	switch serviceNumber {
//	case 1:
//		return "CSV"
//	case 2:
//		return "XLSX"
//	case 3:
//		return "XLSX"
//	case 4:
//		return "XLSX"
//	case 5:
//		return "XLSX"
//	case 6:
//		return "XLSX"
//	case 7:
//		return "XLSX"
//	case 8:
//		return "XLSX"
//	case 9:
//		return "XLSX"
//	case 10:
//		return "XLSX"
//	case 11:
//		return "XLSX"
//	case 12:
//		return "XLSX"
//	case 13:
//		return "XLSX"
//	case 14:
//		return "XLSX"
//	case 15:
//		return "XLSX"
//	case 16:
//		return "XLSX"
//	default:
//		return "... Бот тупит и не может сказать в каком формате. Попробуйте скинуть как есть"
//	}
//}
//
//func GetBot() {
//	cmd.LoadEnv()
//	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
//	//apiToken := os.Getenv("TELEGRAM_BOT_TOKEN_MINION")
//	bot, err := tgbotapi.NewBotAPI(apiToken)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	bot.Debug = true
//
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
//		Offset:  0,
//		Timeout: 60,
//	})
//	//u := tgbotapi.NewUpdate(0)
//	//u.Timeout = 60
//	//
//	//updates := bot.GetUpdatesChan(u)
//
//	state := &BotState{} // Состояние бота
//
//	for update := range updates {
//		if update.Message != nil {
//			go handleMessage(bot, update.Message, state)
//		}
//	}
//}
//
//func isServiceNumber(message string) bool {
//	serviceNumber, err := strconv.Atoi(message)
//	if err != nil {
//		return false
//	}
//
//	correctNumber := serviceNumber >= 1 && serviceNumber <= len(getServices())
//
//	return correctNumber
//}
//
//func sendProcessingMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Принял, начинаю работать 🕖")
//	bot.Send(msg)
//}
//
//func sendUpdateMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Готово 🤩")
//	bot.Send(msg)
//}
//
//func sendCSVFile(bot *tgbotapi.BotAPI, chatID int64, xlsxBuffer *bytes.Buffer, fileName string) {
//	xlsxConfig := tgbotapi.FileBytes{
//		Name:  fileName + ".csv",
//		Bytes: xlsxBuffer.Bytes(),
//	}
//
//	docMsg := tgbotapi.NewDocument(chatID, xlsxConfig)
//	bot.Send(docMsg)
//}
//
//func sendAttentionMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	bot.Send(tgbotapi.NewMessage(chatID, "Обязательно проверь документ❗❗❗\nВозможно он не полностью заполнен"))
//}
//
//func errMsg(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Требуется отправить файл")
//	bot.Send(msg)
//}
//
//func sendServiceList(bot *tgbotapi.BotAPI, chatID int64) {
//	services := getServices()
//
//	var numberedServices strings.Builder
//	for i, service := range services {
//		numberedServices.WriteString(fmt.Sprintf("%d. %s\n", i+1, service))
//	}
//
//	servicesMsg := tgbotapi.NewMessage(chatID, numberedServices.String())
//	bot.Send(servicesMsg)
//}
//
//func sendStartButton(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Выбери сервис (цифру)")
//	startButton := tgbotapi.NewKeyboardButton("Список застройщиков")
//
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(startButton),
//	)
//
//	msg.ReplyMarkup = keyboard
//	bot.Send(msg)
//}
//
//func handleAllClearCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Рота, ОТБОЙ! Играем в 3 скрипа")
//
//	startButton := tgbotapi.NewKeyboardButton("Список застройщиков")
//
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(startButton),
//	)
//
//	replyMsg.ReplyMarkup = keyboard
//	bot.Send(replyMsg)
//}
//
//func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, state *BotState) {
//	switch message.Text {
//	case "Список застройщиков":
//		sendStartButton(bot, message.Chat.ID)
//		sendServiceList(bot, message.Chat.ID)
//	case "Отбой":
//		handleAllClearCommand(bot, message)
//	default:
//		handleDefaultMessage(bot, message, state)
//	}
//}
//
//func handleDefaultMessage(bot *tgbotapi.BotAPI, updateMessage *tgbotapi.Message, state *BotState) {
//
//	messageString := strings.ToLower(updateMessage.Text)
//	chatID := updateMessage.Chat.ID
//
//	if isServiceNumber(messageString) {
//		serviceNumber, _ := strconv.Atoi(messageString)
//
//		if isValidService(serviceNumber) {
//			state.ServiceNum = serviceNumber
//
//			serviceName := getServiceName(serviceNumber)
//			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Вы выбрали сервис %s", serviceName))
//			bot.Send(msg)
//
//			fileFormat := getServiceFileFormat(serviceNumber)
//			msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Пришлите файл в формате %s", fileFormat))
//			bot.Send(msg)
//		} else {
//			msg := tgbotapi.NewMessage(chatID, "Некорректный выбор сервиса")
//			bot.Send(msg)
//		}
//	}
//
//	if updateMessage.Document != nil {
//		if state.ServiceNum == 0 {
//			msg := tgbotapi.NewMessage(chatID, "Сначала напишите номер сервиса")
//			bot.Send(msg)
//		}
//		switch state.ServiceNum {
//		case 1:
//			getServiceAlDar(updateMessage, bot)
//		case 2:
//			getServiceTownX(updateMessage, bot)
//		case 3:
//			getServiceCondor(updateMessage, bot)
//			handleNothingCommand(bot, updateMessage)
//		case 4:
//			getServiceSiadah(updateMessage, bot)
//		case 5:
//			getServiceBinghatii(updateMessage, bot)
//		case 6:
//			getServiceDeyaar(updateMessage, bot)
//		case 7:
//			getServiceEmaar(updateMessage, bot)
//		case 8:
//			getServiceEllingtonProperties(updateMessage, bot)
//		case 9:
//			getServiceAzizi(updateMessage, bot)
//		case 10:
//			getServiceReportageProperties(updateMessage, bot)
//		case 11:
//			getServiceCondor8Cols(updateMessage, bot)
//		case 12:
//			getServiceTiger(updateMessage, bot)
//		case 13:
//			getServiceSothobys(updateMessage, bot)
//		case 14:
//			getServiceObject1(updateMessage, bot)
//		case 15:
//			getServiceSwiss(updateMessage, bot)
//		case 16:
//			getServiceVincitore(updateMessage, bot)
//		}
//	} else if updateMessage.Text != "" {
//		//handleAllClearCommand(bot, updateMessage)
//		sendDefaultMessageTest(bot, updateMessage.Chat.ID)
//
//	} else {
//		//sendDefaultMessage(bot, updateMessage.Chat.ID)
//		sendDefaultMessageTest(bot, updateMessage.Chat.ID)
//	}
//}
//
//func sendDefaultMessage(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Я могу принимать файлы, но не могу с вами общаться 😔\nМогу вам помочь с рефактором вашего файла? 🥺")
//	startButton := tgbotapi.NewKeyboardButton("Список застройщиков")
//
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(startButton),
//	)
//	msg.ReplyMarkup = keyboard
//
//	bot.Send(msg)
//}
//
//func sendErrorMessageFile(bot *tgbotapi.BotAPI, chatID int64) {
//	msg := tgbotapi.NewMessage(chatID, "Требуется отправить файл")
//	bot.Send(msg)
//}
//
//func handleNothingCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	msg := tgbotapi.NewMessage(message.Chat.ID, "Чем еще могу помочь?")
//	startButton := tgbotapi.NewKeyboardButton("Список застройщиков")
//
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(startButton),
//	)
//	msg.ReplyMarkup = keyboard
//
//	bot.Send(msg)
//}
//
//func sendNoActionInlineButton(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
//	buttonText := "Нет действия"
//	buttonData := "no_action"
//	button := tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData)
//	row := tgbotapi.NewInlineKeyboardRow(button)
//	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
//
//	editMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, inlineKeyboard)
//	bot.Send(editMsg)
//}
//
//func sendDefaultMessageTest(bot *tgbotapi.BotAPI, chatID int64) {
//	startButton := tgbotapi.NewKeyboardButton("Список застройщиков")
//
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(startButton),
//	)
//
//	// Очистка предыдущих сообщений и клавиатур в чате
//	clearMsg := tgbotapi.NewMessage(chatID, "")
//	clearMsg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
//	bot.Send(clearMsg)
//
//	// Отправка клавиатуры с кнопкой "Список застройщиков"
//	msg := tgbotapi.NewMessage(chatID, "")
//	msg.ReplyMarkup = keyboard
//	bot.Send(msg)
//}

//package main
//
//import (
//	"genieMap/cmd"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"log"
//	"os"
//	"os/signal"
//	"syscall"
//)
//
//func main() {
//	bot, err := GetBot()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	updates := make(chan tgbotapi.Update, 100)
//
//	go func() {
//		for update := range bot.GetUpdatesChan(tgbotapi.NewUpdate(0)) {
//			updates <- update
//		}
//	}()
//
//	// Обработка входящих обновлений
//	for update := range updates {
//		go handleUpdate(bot, update)
//	}
//
//	// Перехватываем сигналы ОС для корректного завершения работы бота
//	stop := make(chan os.Signal, 1)
//	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
//	<-stop
//}
//
//// GetBot создает и возвращает экземпляр бота
//func GetBot() (*tgbotapi.BotAPI, error) {
//	cmd.LoadEnv()
//	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
//	//apiToken := os.Getenv("TELEGRAM_BOT_TOKEN_MINION")
//	bot, err := tgbotapi.NewBotAPI(apiToken)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	// Настройка параметров бота
//	bot.Debug = true
//
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	return bot, nil
//}
//
//func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
//	// Обработка входящего обновления
//	if update.Message != nil {
//		// Получение информации о пользователе и текста сообщения
//		userID := update.Message.From.ID
//		username := update.Message.From.UserName
//		messageText := update.Message.Text
//
//		log.Printf("Received message from %s (ID: %d): %s", username, userID, messageText)
//
//		// Отправка ответа пользователю
//		reply := "Привет, " + username + "! Ты сказал: " + messageText
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//
//		_, err := bot.Send(msg)
//		if err != nil {
//			log.Println("Error sending message:", err)
//		}
//	}
//}
