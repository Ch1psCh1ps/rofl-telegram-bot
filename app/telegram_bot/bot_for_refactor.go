package telegram_bot

import (
	"bytes"
	"fmt"
	"genieMap/cmd"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
)

var serviceNum int

type BotState struct {
	ServiceNum int
}

func getServices() []string {
	services := []string{
		"Al Dar",
		"Town X",
		"Condor",
		"Siadah",
		"Binghatii",
		"Deyaar",
		"Emaar",
		"Ellington Properties",
		"Azizi",
		"Reportage Properties",
		"Condor 8 колонок",
		"Tiger",
	}

	return services
}

func getServiceName(serviceNumber int) string {
	services := getServices()
	return services[serviceNumber-1]
}

func isValidService(serviceNumber int) bool {
	services := getServices()
	return serviceNumber >= 1 && serviceNumber <= len(services)
}

func getServiceFileFormat(serviceNumber int) string {
	switch serviceNumber {
	case 1:
		return "CSV"
	case 2:
		return "XLSX"
	case 3:
		return "XLSX"
	case 4:
		return "XLSX"
	case 5:
		return "XLSX"
	case 6:
		return "XLSX"
	case 7:
		return "XLSX"
	case 8:
		return "XLSX"
	case 9:
		return "XLSX"
	case 10:
		return "XLSX"
	case 11:
		return "XLSX"
	case 12:
		return "XLSX"
	default:
		return "... Бот тупит и не может сказать в каком формате. Попробуйте скинуть как есть"
	}
}

func GetBot() {
	cmd.LoadEnv()
	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	//apiToken := os.Getenv("TELEGRAM_BOT_TOKEN_MINION")
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	state := &BotState{} // Состояние бота

	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := strings.ToLower(update.Message.Text)
		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(chatID, "Выбери сервис сучка (цифру)")
				bot.Send(msg)

				sendServiceList(bot, chatID)
			case "stop":
				msg := tgbotapi.NewMessage(chatID, "Бот прекратил ответы на команды")
				bot.Send(msg)
				continue //
			}
		}

		if isServiceNumber(message) {
			serviceNumber, _ := strconv.Atoi(message)

			if isValidService(serviceNumber) {
				state.ServiceNum = serviceNumber

				serviceName := getServiceName(serviceNumber)
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Вы выбрали сервис %s", serviceName))
				bot.Send(msg)

				fileFormat := getServiceFileFormat(serviceNumber)
				msg = tgbotapi.NewMessage(chatID, fmt.Sprintf("Пришлите файл в формате %s", fileFormat))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Некорректный выбор сервиса")
				bot.Send(msg)
			}
		}

		if update.Message.Document != nil {
			if state.ServiceNum == 0 {
				msg := tgbotapi.NewMessage(chatID, "Сначала напишите номер сервиса")
				bot.Send(msg)
			} else {
				switch state.ServiceNum {
				case 1:
					getServiceAlDar(update.Message, bot)
				case 2:
					getServiceTownX(update.Message, bot)
				case 3:
					getServiceCondor(update.Message, bot)
				case 4:
					getServiceSiadah(update.Message, bot)
				case 5:
					getServiceBinghatii(update.Message, bot)
				case 6:
					getServiceDeyaar(update.Message, bot)
				case 7:
					getServiceEmaar(update.Message, bot)
				case 8:
					getServiceEllingtonProperties(update.Message, bot)
				case 9:
					getServiceAzizi(update.Message, bot)
				case 10:
					getServiceReportageProperties(update.Message, bot)
				case 11:
					getServiceCondor8Cols(update.Message, bot)
				case 12:
					getServiceTiger(update.Message, bot)
				}
			}
		}
	}
}

func isServiceNumber(message string) bool {
	serviceNumber, err := strconv.Atoi(message)
	if err != nil {
		return false
	}

	correctNumber := serviceNumber >= 1 && serviceNumber <= len(getServices())

	return correctNumber
}

func sendProcessingMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Принял, начинаю работать 🕖")
	bot.Send(msg)
}

func sendUpdateMessage(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Готово 🤩")
	bot.Send(msg)
}

func sendCSVFile(bot *tgbotapi.BotAPI, chatID int64, xlsxBuffer *bytes.Buffer, fileName string) {
	xlsxConfig := tgbotapi.FileBytes{
		Name:  fileName + ".csv",
		Bytes: xlsxBuffer.Bytes(),
	}

	docMsg := tgbotapi.NewDocument(chatID, xlsxConfig)
	bot.Send(docMsg)
}

func sendAttentionMessage(bot *tgbotapi.BotAPI, chatID int64) {
	bot.Send(tgbotapi.NewMessage(chatID, "Обязательно проверь документ❗❗❗\nВозможно он не полностью заполнен"))
}

func errMsg(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Требуется отправить файл")
	bot.Send(msg)
}

func sendServiceList(bot *tgbotapi.BotAPI, chatID int64) {
	services := getServices()

	var numberedServices strings.Builder
	for i, service := range services {
		numberedServices.WriteString(fmt.Sprintf("%d. %s\n", i+1, service))
	}

	servicesMsg := tgbotapi.NewMessage(chatID, numberedServices.String())
	bot.Send(servicesMsg)
}
