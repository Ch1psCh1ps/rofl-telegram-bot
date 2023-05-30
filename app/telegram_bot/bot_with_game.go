package telegram_bot

//
//import (
//	"bytes"
//	"genieMap/app/refactor_xlsx/al_dar"
//	"genieMap/cmd"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"io/ioutil"
//	"log"
//	"math/rand"
//	"os"
//	"strconv"
//)
//
//func StartNewNewBot() {
//	cmd.LoadEnv()
//	apiToken := os.Getenv("TELEGRAM_BOT_TOKEN")
//	bot, err := tgbotapi.NewBotAPI(apiToken)
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
//	case "Отбой":
//		handleAllClearCommand(bot, message)
//	case "СЮЮЮЮДАААА!!!111!!1!!":
//		handleYesCommand(bot, message)
//	case "Нет":
//		handleNoCommand(bot, message)
//	case "Игра":
//		sendPhotoMessage(bot, message.Chat.ID, photoPath)
//	case "Пойду работать":
//		handleWorkCommand(bot, message)
//	case "Тогда давай игру! 🕹":
//		handleGameCommand(bot, message)
//	case "Ничем":
//		handleNothingCommand(bot, message)
//	default:
//		handleDefaultMessage(bot, message, photoPath)
//	}
//}
//
//func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	buttonMsg := tgbotapi.NewMessage(message.Chat.ID, "Чем могу помочь, Педик? Сам то уже не справляешься")
//
//	buttonText := "Рефактор xlsx"
//	refactorButton := tgbotapi.NewKeyboardButton(buttonText)
//
//	buttonAllClear := "Отбой"
//	AllClearButton := tgbotapi.NewKeyboardButton(buttonAllClear)
//
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(refactorButton),
//		tgbotapi.NewKeyboardButtonRow(AllClearButton),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	buttonMsg.ReplyMarkup = keyboard
//
//	// Отправка сообщения с кнопками
//	bot.Send(buttonMsg)
//}
//
//func handleRefactorCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Пришлите файл с расширением .xlsx  или сиськи 😋")
//	bot.Send(replyMsg)
//}
//
//func handleAllClearCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Рота, ОТБОЙ! Играем в 3 скрипа")
//
//	// Завершение работы бота
//	// Создание кнопки "/start"
//	startButton := tgbotapi.NewKeyboardButton("/start")
//
//	// Создание клавиатуры и добавление кнопок
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(startButton),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	replyMsg.ReplyMarkup = keyboard
//	bot.Send(replyMsg)
//}
//
//func handleYesCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	urlMsg := tgbotapi.NewMessage(message.Chat.ID, "https://2048game.com/ru/")
//
//	// Создание сообщения
//	talkMsg := tgbotapi.NewMessage(message.Chat.ID, "Разминаем мозг 🧠")
//
//	// Создание сообщения "Еще?"
//	moreMsg := tgbotapi.NewMessage(message.Chat.ID, "Еще? (:")
//
//	// Создание клавиатуры и добавление кнопок
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("СЮЮЮЮДАААА!!!111!!1!!")),
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Пойду работать")),
//		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Нет")),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	urlMsg.ReplyMarkup = keyboard
//
//	bot.Send(urlMsg)
//	bot.Send(talkMsg)
//	bot.Send(moreMsg)
//}
//
//func handleNoCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Настоящего дотера ответ. Я уже хотел включить гачи ): Пиздуй работать.\nЧем еще могу помочь?")
//
//	// Создание кнопки "Рефактор xlsx"
//	buttonText := "Рефактор xlsx"
//	refactorButton := tgbotapi.NewKeyboardButton(buttonText)
//
//	// Создание кнопки "Ничем"
//	nothingButton := tgbotapi.NewKeyboardButton("Ничем")
//	gameButton := tgbotapi.NewKeyboardButton("Игра")
//
//	// Создание клавиатуры и добавление кнопок
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(refactorButton),
//		tgbotapi.NewKeyboardButtonRow(nothingButton),
//		tgbotapi.NewKeyboardButtonRow(gameButton),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	replyMsg.ReplyMarkup = keyboard
//
//	bot.Send(replyMsg)
//}
//
//func handleWorkCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Молодец, но зп все равно не повысят (:\nЧем еще могу помочь?")
//
//	// Создание кнопки "Рефактор xlsx"
//	buttonText := "Рефактор xlsx"
//	refactorButton := tgbotapi.NewKeyboardButton(buttonText)
//
//	// Создание кнопки "Ничем"
//	nothingButton := tgbotapi.NewKeyboardButton("Ничем")
//
//	orGameText := "Тогда давай игру! 🕹"
//	orGameButton := tgbotapi.NewKeyboardButton(orGameText)
//
//	// Создание клавиатуры и добавление кнопок
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(refactorButton),
//		tgbotapi.NewKeyboardButtonRow(nothingButton),
//		tgbotapi.NewKeyboardButtonRow(orGameButton),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	replyMsg.ReplyMarkup = keyboard
//
//	bot.Send(replyMsg)
//}
//
//func handleDefaultMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, photoPath string) {
//	if message.Document != nil {
//		//if isXLSXFile(message.Document.FileName) {
//		fileID := message.Document.FileID
//		fileURL, err := bot.GetFileDirectURL(fileID)
//		if err != nil {
//			log.Printf("Ошибка при получении файла: %v", err)
//			sendErrorMessage(bot, message.Chat.ID)
//			return
//		}
//
//		sendProcessingMessage(bot, message.Chat.ID)
//
//		xlsxBuffer, err := al_dar.DoBookCSV(fileURL)
//		if err != nil {
//			log.Printf("Ошибка при обработке файла: %v", err)
//			sendErrorMessage(bot, message.Chat.ID)
//			return
//		}
//
//		sendUpdateMessage(bot, message.Chat.ID)
//		sendXLSXFile(bot, message.Chat.ID, xlsxBuffer)
//
//		sendPhotoMessage(bot, message.Chat.ID, photoPath)
//
//		//} else {
//		//	sendInvalidFileMessage(bot, message.Chat.ID)
//		//}
//	} else if message.Text != "" {
//		// Обработка выбранного поля для игры
//		handleGameMove(bot, message)
//
//	} else {
//		// Отправка стандартного сообщения
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
//		Name:  "refactor_available_units.csv",
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
//
//		// Создание кнопки "ДА!"
//		yesButton := tgbotapi.NewKeyboardButton("СЮЮЮЮДАААА!!!111!!1!!")
//
//		// Создание кнопки "Нет"
//		noButton := tgbotapi.NewKeyboardButton("Нет")
//
//		// Создание клавиатуры и добавление кнопок
//		keyboard := tgbotapi.NewReplyKeyboard(
//			tgbotapi.NewKeyboardButtonRow(yesButton, noButton),
//		)
//
//		// Прикрепление клавиатуры к сообщению
//		photoMsg.ReplyMarkup = keyboard
//
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
//
//func handleNothingCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	msg := tgbotapi.NewMessage(message.Chat.ID, "Дружище, съеби нахуй .|.")
//
//	// Завершение работы бота
//	// Создание кнопки "/start"
//	startButton := tgbotapi.NewKeyboardButton("/start")
//
//	// Создание клавиатуры и добавление кнопок
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(startButton),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	msg.ReplyMarkup = keyboard
//
//	bot.Send(msg)
//}
//
//func isXLSXFile(filename string) bool {
//	return len(filename) >= 5 && filename[len(filename)-4:] == ".csv"
//}
//
//var (
//	gameState       [3][3]int
//	isPlayer1Turn   bool
//	isSinglePlayer  bool
//	lastMoveByHuman bool
//)
//
//func handleGameCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	clearGame()
//	createGame()
//
//	replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Давай сыграем в крестики-нолики! Чтобы сделать ход, отправь номер ячейки от 1 до 9.")
//
//	// Создание клавиатуры для игрового поля
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("1"),
//			tgbotapi.NewKeyboardButton("2"),
//			tgbotapi.NewKeyboardButton("3"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("4"),
//			tgbotapi.NewKeyboardButton("5"),
//			tgbotapi.NewKeyboardButton("6"),
//		),
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("7"),
//			tgbotapi.NewKeyboardButton("8"),
//			tgbotapi.NewKeyboardButton("9"),
//		),
//	)
//
//	// Прикрепление клавиатуры к сообщению
//	replyMsg.ReplyMarkup = keyboard
//
//	bot.Send(replyMsg)
//}
//
//func handleGameMove(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//	// Получить номер ячейки из сообщения пользователя
//	move, err := strconv.Atoi(message.Text)
//	if err != nil || move < 1 || move > 9 {
//		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Неверный номер ячейки. Попробуй еще раз.")
//		bot.Send(replyMsg)
//		return
//	}
//
//	// Перевести номер ячейки в индексы массива
//	row := (move - 1) / 3
//	col := (move - 1) % 3
//
//	// Проверить, является ли выбранная ячейка пустой
//	if gameState[row][col] != 0 {
//		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Эта ячейка уже занята. Попробуй другую.")
//		bot.Send(replyMsg)
//		return
//	}
//
//	// Сделать ход в выбранную ячейку
//	gameState[row][col] = 1 // Ход игрока всегда будет крестиком
//
//	// Проверить, есть ли победитель или ничья
//	if checkWin() {
//		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Поздравляю! Ты победил наитупейшего бота!\nЕсть чем гордится не так ли?")
//		bot.Send(replyMsg)
//		clearGame()
//		sendGameEndOptions(bot, message.Chat.ID)
//		return
//	} else if isDraw() {
//		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Ничья! Игра окончена.")
//		bot.Send(replyMsg)
//		clearGame()
//		sendGameEndOptions(bot, message.Chat.ID)
//		return
//	}
//
//	// Сменить ход
//	isPlayer1Turn = false
//	lastMoveByHuman = true
//
//	// Сделать ход компьютера
//	makeComputerMove()
//
//	// Проверить, есть ли победитель или ничья после хода компьютера
//	if checkWin() {
//		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Компьютер победил! А вот и наш засранец. Прими мои поздравления. Ты проиграл тупому боту...\n"+
//			"Ты полый неудачник. Теперь кайфуй в полном идиночестве, осозновая что теперь никто не хочет с тобой тусить из-за того какой ты огромный кусок говнины\n"+
//			"Расселся тут король говна на троне одиночества. Да славится его величество - самы жалкий чмошник во всем космосе\n"+
//			"Да здравствует кал-великий-царь-дрестун каждый день-перень на долгие века")
//		bot.Send(replyMsg)
//		clearGame()
//		sendGameEndOptions(bot, message.Chat.ID)
//		return
//	} else if isDraw() {
//		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Ничья! Игра окончена.")
//		bot.Send(replyMsg)
//		clearGame()
//		sendGameEndOptions(bot, message.Chat.ID)
//		return
//	}
//
//	// Отправить обновленное игровое поле
//	sendGameBoard(bot, message.Chat.ID)
//}
//
////2 работает с репитом но игра нормально не начинается
////func handleGameMove(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
////	// Получить номер ячейки из сообщения пользователя
////	move, err := strconv.Atoi(message.Text)
////	if err != nil || move < 1 || move > 9 {
////		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Неверный номер ячейки. Попробуй еще раз.")
////		bot.Send(replyMsg)
////		sendGameBoard(bot, message.Chat.ID) // Отправить обновленное игровое поле
////		return
////	}
////
////	// Перевести номер ячейки в индексы массива
////	row := (move - 1) / 3
////	col := (move - 1) % 3
////
////	// Проверить, является ли выбранная ячейка пустой
////	if gameState[row][col] != 0 {
////		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Эта ячейка уже занята. Попробуй другую.")
////		bot.Send(replyMsg)
////		sendGameBoard(bot, message.Chat.ID) // Отправить обновленное игровое поле
////		return
////	}
////
////	// Сделать ход в выбранную ячейку
////	gameState[row][col] = 1 // Ход игрока всегда будет крестиком
////
////	// Проверить, есть ли победитель или ничья
////	if checkWin() {
////		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Поздравляю! Вы победили!")
////		bot.Send(replyMsg)
////		clearGame()
////		sendGameEndOptions(bot, message.Chat.ID)
////		return
////	} else if isDraw() {
////		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Ничья! Игра окончена.")
////		bot.Send(replyMsg)
////		clearGame()
////		sendGameEndOptions(bot, message.Chat.ID)
////		return
////	}
////
////	// Сменить ход
////	isPlayer1Turn = false
////	lastMoveByHuman = true
////
////	// Сделать ход компьютера
////	makeComputerMove()
////
////	// Проверить, есть ли победитель или ничья после хода компьютера
////	if checkWin() {
////		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Компьютер победил! Вы проиграли.")
////		bot.Send(replyMsg)
////		clearGame()
////		sendGameEndOptions(bot, message.Chat.ID)
////		return
////	} else if isDraw() {
////		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "Ничья! Игра окончена.")
////		bot.Send(replyMsg)
////		clearGame()
////		sendGameEndOptions(bot, message.Chat.ID)
////		return
////	}
////
////	// Отправить обновленное игровое поле
////	sendGameBoard(bot, message.Chat.ID)
////}
//
//func makeComputerMove() {
//	// Выбрать случайную пустую ячейку для хода компьютера
//	emptyCells := getEmptyCells()
//
//	if len(emptyCells) == 0 {
//		return
//	}
//
//	index := rand.Intn(len(emptyCells))
//	cell := emptyCells[index]
//
//	gameState[cell.row][cell.col] = 2 // Ход компьютера всегда будет ноликом
//}
//
//func getEmptyCells() []Cell {
//	var emptyCells []Cell
//
//	for i := 0; i < 3; i++ {
//		for j := 0; j < 3; j++ {
//			if gameState[i][j] == 0 {
//				emptyCells = append(emptyCells, Cell{row: i, col: j})
//			}
//		}
//	}
//
//	return emptyCells
//}
//
//type Cell struct {
//	row int
//	col int
//}
//
//func checkWin() bool {
//	// Проверить все возможные комбинации для победы
//
//	// Горизонтальные линии
//	for i := 0; i < 3; i++ {
//		if gameState[i][0] != 0 && gameState[i][0] == gameState[i][1] && gameState[i][0] == gameState[i][2] {
//			return true
//		}
//	}
//
//	// Вертикальные линии
//	for i := 0; i < 3; i++ {
//		if gameState[0][i] != 0 && gameState[0][i] == gameState[1][i] && gameState[0][i] == gameState[2][i] {
//			return true
//		}
//	}
//
//	// Диагонали
//	if gameState[0][0] != 0 && gameState[0][0] == gameState[1][1] && gameState[0][0] == gameState[2][2] {
//		return true
//	}
//
//	if gameState[0][2] != 0 && gameState[0][2] == gameState[1][1] && gameState[0][2] == gameState[2][0] {
//		return true
//	}
//
//	return false
//}
//
//func isDraw() bool {
//	// Проверить, есть ли незанятые ячейки
//	for i := 0; i < 3; i++ {
//		for j := 0; j < 3; j++ {
//			if gameState[i][j] == 0 {
//				return false
//			}
//		}
//	}
//
//	return true
//}
//
//func sendGameBoard(bot *tgbotapi.BotAPI, chatID int64) {
//	// Отправить текущее игровое поле пользователю
//	board := ""
//	for i := 0; i < 3; i++ {
//		for j := 0; j < 3; j++ {
//			cell := ""
//			switch gameState[i][j] {
//			case 0:
//				cell = "-"
//			case 1:
//				cell = "X"
//			case 2:
//				cell = "O"
//			}
//			board += cell + " "
//		}
//		board += "\n"
//	}
//
//	replyMsg := tgbotapi.NewMessage(chatID, board)
//	bot.Send(replyMsg)
//}
//
//func clearGame() {
//	gameState = [3][3]int{}
//}
//
//func createGame() {
//	// Создать новое состояние игры
//	// Например, можно инициализировать массив или структуру данных, представляющую игровое поле.
//	// В данном примере, используем массив размером 3x3 для представления игрового поля.
//	// Пустая ячейка будет обозначаться значением "0", крестик - "1", нолик - "2".
//	gameState = [3][3]int{
//		{0, 0, 0},
//		{0, 0, 0},
//		{0, 0, 0},
//	}
//}
//
//func sendGameEndOptions(bot *tgbotapi.BotAPI, chatID int64) {
//	replyMsg := tgbotapi.NewMessage(chatID, "Игра окончена, тупица. Теперь все с начала)")
//	keyboard := tgbotapi.NewReplyKeyboard(
//		tgbotapi.NewKeyboardButtonRow(
//			tgbotapi.NewKeyboardButton("/start"),
//		),
//	)
//	replyMsg.ReplyMarkup = keyboard
//	bot.Send(replyMsg)
//}
