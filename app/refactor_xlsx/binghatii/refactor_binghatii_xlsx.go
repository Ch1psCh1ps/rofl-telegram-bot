package binghatii

import (
	"bytes"
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
	fileContent, downloadFileErr := downloadFile(path)
	if downloadFileErr != nil {
		LogError("Ошибка при загрузке файла: %v", downloadFileErr)
		return nil, downloadFileErr
	}

	data := GetDataFromBytes(fileContent)

	sheetName := data.GetSheetName(0)

	cols, sheetErr := GetSheet(data, sheetName)
	if sheetErr != nil {
		LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
		cols, sheetErr = GetSheet(data, "Sheet1")
		if sheetErr != nil {
			LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
			return nil, sheetErr
		}
	}

	defer func() {
		if err := data.Close(); err != nil {
			LogError("%v", err)
		}
	}()

	newXlsxFile := excelize.NewFile()

	defer func() {
		if err := newXlsxFile.Close(); err != nil {
			LogError("%v", err)
		}
	}()

	setColumnValues(newXlsxFile, cols[1], "A")    //number
	setColumnValues(newXlsxFile, cols[4], "B")    //price
	setColumnValues(newXlsxFile, cols[7], "C")    //Square
	setColumnValues(newXlsxFile, []string{}, "D") //height
	setColumnValues(newXlsxFile, cols[2], "E")    //type
	setColumnValues(newXlsxFile, cols[3], "F")    //layout
	setColumnValues(newXlsxFile, cols[6], "G")    //views

	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
	replaceUnitViewsFieldInXLSX(newXlsxFile, 6)

	buffer, err3 := cmd.ConvertXlsxToCsv(newXlsxFile)

	if err3 != nil {
		LogError("Ошибка при конвертации XLSX в CSV: %v", err3)
		return nil, err3
	}

	refactorFile, errRefactor := cmd.RemoveAnyRowFromCSV(buffer, 2)

	if errRefactor != nil {
		LogError("%v", errRefactor)
		return nil, errRefactor
	}

	buf, errFirstRow := cmd.UpdateFirstRowInCSV(refactorFile, _struct.GetNameFirstRow())
	if errFirstRow != nil {
		LogError("Ошибка при добавлении строки", errFirstRow)
		return nil, errFirstRow
	}

	return buf, nil
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

func GetDataFromBytes(fileContent []byte) *excelize.File {
	buffer := bytes.NewBuffer(fileContent)
	data, err := excelize.OpenReader(buffer)
	if err != nil {
		LogError("%v", err)
	}
	return data
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

func replaceUnitNumberFieldInXLSX(file *excelize.File, indexOfCell int) error {
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
					word := strings.Replace(row[colIndex], "BUGA-", "", 1)
					//word := strings.Fields(row[colIndex])
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

func replaceUnitLayoutFieldInXLSX(file *excelize.File, indexOfCell int) error {
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
					words1 := strings.Split(row[colIndex], ",")

					for i, word1 := range words1 {
						if strings.Contains(word1, "&") || strings.Contains(word1, "/") {
							fmt.Println(words1[i], "words1")
							words1[i] = strings.Replace(word1, "/", " or ", -1)
						}
					}

					word := strings.Join(words1, ",")
					word = strings.Replace(word, "+", "/", -1)
					word = strings.Replace(word, ",", "/", -1)
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
					word := strings.Replace(row[colIndex], ",", "/", -1)
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
