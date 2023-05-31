package Condor

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

	setColumnValues(newXlsxFile, cols[0], "A") //number
	setColumnValues(newXlsxFile, cols[8], "B") //price
	setColumnValues(newXlsxFile, cols[5], "C") //Square
	setColumnValues(newXlsxFile, cols[1], "D") //height
	setColumnValues(newXlsxFile, cols[3], "E") //type
	setColumnValues(newXlsxFile, cols[2], "F") //layout
	setColumnValues(newXlsxFile, cols[4], "G") //views

	//errRefactPrice := refactorUnitPriceFieldInXLSX(newXlsxFile, 1)
	//if errRefactPrice != nil {
	//	LogError("Рефактор ошибка %v", errRefactPrice)
	//}

	buffer, err3 := cmd.ConvertXlsxToCsv(newXlsxFile)

	if err3 != nil {
		LogError("Ошибка при конвертации XLSX в CSV: %v", err3)
		return nil, err3
	}

	//math.Round(x*100)/100

	buf, errFirstRow := cmd.UpdateFirstRowInCSV(buffer, _struct.GetNameFirstRow())
	if errFirstRow != nil {
		LogError("Ошибка при добавлении строки", errFirstRow)

		return nil, errFirstRow
	}

	//buffer1, err4 := refactorUnitPriceFieldInCSV(buf, 1)
	//
	//if err4 != nil {
	//	LogError("Ошибка при замене цены: %v", err3)
	//	return nil, err4
	//}

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

func refactorUnitPriceFieldInXLSX(file *excelize.File, indexOfCell int) error {
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
		rowIndex := 0

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
					fmt.Println(row[colIndex], "row[colIndex]")
					fmt.Println(row[colIndex], "row[colIndex]")
					fmt.Println(row[colIndex], "row[colIndex]")
					fmt.Println(row[colIndex], "row[colIndex]")
					fmt.Println(row[colIndex], "row[colIndex]")
					fmt.Println(row[colIndex], "row[colIndex]")
					fmt.Println(row[colIndex], "row[colIndex]")

					word, errRound := cmd.RoundString(row[colIndex])
					fmt.Println(word)
					fmt.Println(word)
					fmt.Println(word)
					fmt.Println(word)
					fmt.Println(word)

					if errRound != nil {
						return errRound
					}
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

func refactorUnitPriceFieldInCSV(data *bytes.Buffer, indexOfColumn int) (*bytes.Buffer, error) {
	// Преобразуем данные в формат CSV
	reader := csv.NewReader(data)

	// Читаем все записи из CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Создаем новый буфер для преобразованных данных
	transformedData := bytes.NewBuffer(nil)

	// Обрабатываем каждую запись
	for _, record := range records {
		// Проверяем, что индекс столбца находится в допустимом диапазоне
		if indexOfColumn >= 0 && indexOfColumn < len(record) {
			// Получаем значение ячейки
			cellValue := record[indexOfColumn]

			// Округляем значение ячейки до двух знаков после запятой
			word, errRound := cmd.RoundString(cellValue)
			if errRound != nil {
				return nil, errRound
			}

			// Заменяем значение ячейки на округленное
			record[indexOfColumn] = word
		}

		// Записываем преобразованную запись в буфер
		transformedData.WriteString(strings.Join(record, ","))
		transformedData.WriteByte('\n')
	}

	return transformedData, nil
}
