package al_dar

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"genieMap/cmd"
	_struct "genieMap/structures"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DoBookCSV(path string) (*bytes.Buffer, error) {
	filecsv, convertErr := ConvertCSVtoXLSX(path)
	if convertErr != nil {
		LogError("Ошибка при конвертации CSV в XLSX: %v", convertErr)
		return nil, convertErr
	}

	sheetName := filecsv.GetSheetName(0)

	cols, sheetErr := GetSheet(filecsv, sheetName)
	if sheetErr != nil {
		LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
		return nil, sheetErr
	}

	newXlsxFile := excelize.NewFile()

	defer func() {
		if err := newXlsxFile.Close(); err != nil {
			LogError("%v", err)
		}
	}()

	setColumnValues(newXlsxFile, cols[2], "A")    //number
	setColumnValues(newXlsxFile, cols[11], "B")   //price
	setColumnValues(newXlsxFile, cols[10], "C")   //Square
	setColumnValues(newXlsxFile, []string{}, "D") //height
	setColumnValues(newXlsxFile, cols[4], "E")    //type
	setColumnValues(newXlsxFile, cols[5], "F")    //layout
	setColumnValues(newXlsxFile, cols[6], "G")    //views

	ReplaceXLSXNumber(newXlsxFile, 0)
	ReplaceXLSXLayout(newXlsxFile, 5)
	ReplaceXLSXType(newXlsxFile, 4)
	replaceUnitViewsFieldInXLSX(newXlsxFile, 6)

	buf, err3 := convertXlsxToCsv(newXlsxFile)
	if err3 != nil {
		LogError("Ошибка при конвертации XLSX в CSV: %v", err3)

		return nil, err3
	}

	buffer, errFirstRow := cmd.UpdateFirstRowInCSV(buf, _struct.GetNameFirstRow())
	if errFirstRow != nil {
		LogError("Ошибка при замене первой строки", errFirstRow)

		return nil, errFirstRow
	}

	return buffer, nil
}

func ReplaceXLSXLayout(file *excelize.File, indexOfCell int) error {
	//не забудь поменять значение word!!!!
	// Получаем список имен листов из файла
	sheets := file.GetSheetList()

	// Обрабатываем каждый лист
	for _, sheet := range sheets {
		// Получаем все строки в листе
		rows, err := file.Rows(sheet)
		if err != nil {
			return err
		}

		// Инициализируем индекс строки
		rowIndex := 1

		// Перебираем каждую строку
		for rows.Next() {
			row, err2 := rows.Columns()
			if err2 != nil {
				return err2
			}

			colIndex := indexOfCell

			// Перебираем каждую ячейку в строке
			for _, cellValue := range row {
				if cellValue == row[colIndex] {
					// Заменяем значение ячейки на замену
					word := strings.Split(cellValue, "")
					row[colIndex] = word[0] + " BR"

					// Получаем имя столбца на основе индекса столбца
					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					// Обновляем значение ячейки в листе
					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), row[colIndex])
					if err1 != nil {
						return err1
					}
					break
				}
			}

			// Увеличиваем индекс строки
			rowIndex++
		}
	}

	return nil
}

func ReplaceXLSXNumber(file *excelize.File, indexOfCell int) error {
	//не забудь поменять значение word!!!!
	// Получаем список имен листов из файла
	sheets := file.GetSheetList()

	// Обрабатываем каждый лист
	for _, sheet := range sheets {
		// Получаем все строки в листе
		rows, err := file.Rows(sheet)
		if err != nil {
			return err
		}

		// Инициализируем индекс строки
		rowIndex := 1

		// Перебираем каждую строку
		for rows.Next() {
			row, err2 := rows.Columns()
			if err2 != nil {
				return err2
			}

			colIndex := indexOfCell

			// Перебираем каждую ячейку в строке
			for _, cellValue := range row {
				if cellValue == row[colIndex] && cellValue != "EndUnitCode" {

					word := strings.Split(row[colIndex], "-")
					divisionByNumber := strings.Split(word[len(word)-3], "_")
					row[colIndex] =
						divisionByNumber[len(divisionByNumber)-1] +
							word[len(word)-2] +
							word[len(word)-1]

					// Получаем имя столбца на основе индекса столбца
					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					// Обновляем значение ячейки в листе
					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), row[colIndex])
					if err1 != nil {
						return err1
					}
					break
				}
			}

			// Увеличиваем индекс строки
			rowIndex++
		}
	}

	return nil
}

func ReplaceXLSXType(file *excelize.File, indexOfCell int) error {
	//не забудь поменять значение word!!!!
	// Получаем список имен листов из файла
	sheets := file.GetSheetList()

	// Обрабатываем каждый лист
	for _, sheet := range sheets {
		// Получаем все строки в листе
		rows, err := file.Rows(sheet)
		if err != nil {
			return err
		}

		// Инициализируем индекс строки
		rowIndex := 1

		// Перебираем каждую строку
		for rows.Next() {
			row, err2 := rows.Columns()
			if err2 != nil {
				return err2
			}

			colIndex := indexOfCell

			// Перебираем каждую ячейку в строке
			for _, cellValue := range row {
				if cellValue == row[colIndex] {
					// Заменяем значение ячейки на замену
					word := strings.ToLower(cellValue)
					row[colIndex] = word

					// Получаем имя столбца на основе индекса столбца
					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					// Обновляем значение ячейки в листе
					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), word)
					if err1 != nil {
						return err1
					}
					break
				}
			}

			// Увеличиваем индекс строки
			rowIndex++
		}
	}

	return nil
}

func replaceUnitViewsFieldInXLSX(file *excelize.File, indexOfCell int) error {
	sheets := file.GetSheetList()

	for _, sheet := range sheets {
		rows, err := file.Rows(sheet)
		if err != nil {
			return err
		}

		rowIndex := 1

		for rows.Next() {
			row, err2 := rows.Columns()
			if err2 != nil {
				return err2
			}

			colIndex := indexOfCell

			for _, cellValue := range row {
				if cellValue == row[colIndex] {
					wordContains := strings.Contains(cellValue, "view")
					wordContains1 := strings.Contains(cellValue, "View")
					if wordContains == true || wordContains1 == true {
						replaceWord := strings.ReplaceAll(cellValue, "view", "")
						replaceWord = strings.ReplaceAll(replaceWord, "View", "")
						replaceWordArray := strings.Split(replaceWord, "Street")
						replaceWord = strings.Join(replaceWordArray, " Street")
						replaceWordArray = strings.Split(replaceWord, "park")
						replaceWord = strings.Join(replaceWordArray, " park")

						row[colIndex] = replaceWord

						columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
						if err1 != nil {
							return err1
						}

						err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), row[colIndex])
						if err1 != nil {
							return err1
						}
						break
					} else {
						row[colIndex] = ""

						columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
						if err1 != nil {
							return err1
						}

						// Обновляем значение ячейки в листе
						err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), row[colIndex])
						if err1 != nil {
							return err1
						}
						break
					}
				}
			}

			// Увеличиваем индекс строки
			rowIndex++
		}
	}

	return nil
}

func ReplaceNumberXLSX(file *excelize.File) error {
	sheets := file.GetSheetList()

	// Обрабатываем каждый лист
	for _, sheet := range sheets {
		rows, err := file.Rows(sheet)
		if err != nil {
			return err
		}

		rowIndex := 1

		for rows.Next() {
			row, err2 := rows.Columns()
			if err2 != nil {
				return err2
			}

			for colIndex, cellValue := range row {
				if cellValue == row[colIndex] {
					word := strings.Split(row[colIndex], "_")
					row[colIndex] = word[len(word)-1]

					// Получаем имя столбца на основе индекса столбца
					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					// Обновляем значение ячейки в листе
					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), word[len(word)-1])
					if err1 != nil {
						return err1
					}
					break
				}
			}

			// Увеличиваем индекс строки
			rowIndex++
		}
	}

	return nil
}

func downloadFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	fileContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func GetSheet(data *excelize.File, sheetName string) ([][]string, error) {
	cols, err := data.GetCols(sheetName)
	if err != nil {
		return nil, err
	}
	return cols, nil
}

func setColumnValues(file *excelize.File, values []string, colPrefix string) {
	for i, value := range values {
		cell := colPrefix + strconv.Itoa(i+1)
		file.SetCellValue("Sheet1", cell, value)
	}
}

func LogError(format string, v ...interface{}) {
	logFile := "refactor.log"
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла логов: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Printf(format, v...)
}

func ConvertCSVtoXLSX(csvURL string) (*excelize.File, error) {
	// Получение данных из ссылки
	response, err := http.Get(csvURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить данные из ссылки: %v", err)
	}

	// Чтение CSV-файла
	reader := csv.NewReader(response.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать CSV-файл: %v", err)
	}

	defer response.Body.Close()

	//// Создание файла на диске
	//file1, err := os.Create("/Users/anastasyaplotnikova/GolandProjects/genieMap/telegram.csv")
	//if err != nil {
	//	return nil, fmt.Errorf("не удалось создать файл: %v", err)
	//}
	//defer file1.Close()
	//
	//// Копирование содержимого из response.Body в файл
	//_, err = io.Copy(file1, response.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("не удалось записать файл: %v", err)
	//}
	//
	//// Открытие существующего файла CSV
	//file, err := os.Open("/Users/anastasyaplotnikova/GolandProjects/genieMap/telegram.csv")
	//if err != nil {
	//	return nil, fmt.Errorf("не удалось открыть файл CSV: %v", err)
	//}
	//defer file.Close()

	// Создание нового файла XLSX
	xlsxFile := excelize.NewFile()

	// Создание нового листа
	sheetName := "Sheet1"

	// Запись данных из CSV-файла в XLSX
	for rowIndex, record := range records {
		for colIndex, cellValue := range record {
			cell, _ := excelize.ColumnNumberToName(colIndex + 1)
			cell += fmt.Sprintf("%d", rowIndex+1)
			xlsxFile.SetCellValue(sheetName, cell, cellValue)
		}
	}

	//xlsxFilePathForSave := "/Users/anastasyaplotnikova/GolandProjects/genieMap/refactor.xlsx"
	//
	////// Сброс указателя позиции чтения/записи в начало файла XLSX
	////xlsxFile.SetActiveSheet(0)
	//
	//// Сохранение файла XLSX на диске
	//if err1 := xlsxFile.SaveAs(xlsxFilePathForSave); err1 != nil {
	//	return nil, fmt.Errorf("не удалось сохранить файл XLSX: %v", err1)
	//}

	return xlsxFile, nil
}

func convertXlsxToCsv(inputFile *excelize.File) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Извлечение данных из каждого листа XLSX и запись в CSV
	sheetList := inputFile.GetSheetList()
	for _, sheetName := range sheetList {
		rows, err := inputFile.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении строк листа XLSX: %v", err)
		}

		for _, row := range rows {
			err1 := writer.Write(row)
			if err1 != nil {
				return nil, fmt.Errorf("ошибка при записи строки в CSV: %v", err1)
			}
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("ошибка при записи в буфер CSV: %v", err)
	}

	return buf, nil
}

func DownloadFileFromTelegramBot(token, fileID, filePath string) ([]byte, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", token, fileID)
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о файле: %v", err)
	}
	defer response.Body.Close()

	// Чтение ответа в формате JSON
	// Здесь вам потребуется ваша собственная логика для извлечения информации о файле
	// из ответа Telegram API, чтобы получить ссылку на файл

	fileURL := "https://api.telegram.org/file/bot" + token + "/" + filePath

	// Загрузка файла
	fileResponse, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить файл: %v", err)
	}
	defer fileResponse.Body.Close()

	// Чтение содержимого файла в []byte
	fileContent, err := ioutil.ReadAll(fileResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать содержимое файла: %v", err)
	}

	return fileContent, nil
}

func SaveXLSXFile(file *excelize.File, outputPath string) error {
	if err := file.SaveAs(outputPath); err != nil {
		return fmt.Errorf("не удалось сохранить файл XLSX: %v", err)
	}
	return nil
}

func GetDataFromBytes(fileContent []byte) *excelize.File {
	buffer := bytes.NewBuffer(fileContent)
	data, err := excelize.OpenReader(buffer)
	if err != nil {
		LogError("%v", err)
	}
	return data
}
