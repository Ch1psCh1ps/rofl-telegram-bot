package binghatii

import (
	"bytes"
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

//func DoBookCSV(path string) (*bytes.Buffer, error) {
//	fileContent, downloadFileErr := downloadFile(path)
//	if downloadFileErr != nil {
//		LogError("Ошибка при загрузке файла: %v", downloadFileErr)
//		return nil, downloadFileErr
//	}
//
//	data := GetDataFromBytes(fileContent)
//
//	sheetName := data.GetSheetName(0)
//
//	cols, sheetErr := GetSheet(data, sheetName)
//	if sheetErr != nil {
//		LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
//		cols, sheetErr = GetSheet(data, "Sheet1")
//		if sheetErr != nil {
//			LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
//			return nil, sheetErr
//		}
//	}
//
//	defer func() {
//		if err := data.Close(); err != nil {
//			LogError("%v", err)
//		}
//	}()
//
//	newXlsxFile := excelize.NewFile()
//
//	defer func() {
//		if err := newXlsxFile.Close(); err != nil {
//			LogError("%v", err)
//		}
//	}()
//
//	for i, col := range cols {
//		for _, colName := range col {
//			switch colName {
//			case "Unit Code":
//				setColumnValues(newXlsxFile, cols[i], "A") //number
//			case "Target Unit Price":
//				setColumnValues(newXlsxFile, cols[i], "B") //price
//			case "Total Area":
//				setColumnValues(newXlsxFile, cols[i], "C") //Square
//			case "Unit Type":
//				setColumnValues(newXlsxFile, cols[i], "F") //layout
//				//replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
//			case "Description":
//				setColumnValues(newXlsxFile, cols[i], "F") //layout
//				//replaceUnitLayoutDescriptionFieldInXLSX(newXlsxFile, 5)
//			case "View":
//				setColumnValues(newXlsxFile, cols[i], "G") //views
//			}
//		}
//	}
//	setColumnValues(newXlsxFile, []string{}, "E") //type
//	setColumnValues(newXlsxFile, []string{}, "D") //height
//
//	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
//	replaceUnitViewsFieldInXLSX(newXlsxFile, 6)
//
//	buffer, err3 := cmd.ConvertXlsxToCsv(newXlsxFile)
//
//	if err3 != nil {
//		LogError("Ошибка при конвертации XLSX в CSV: %v", err3)
//		return nil, err3
//	}
//
//	refactorFile, errRefactor := cmd.RemoveAnyRowFromCSV(buffer, 2)
//
//	if errRefactor != nil {
//		LogError("%v", errRefactor)
//		return nil, errRefactor
//	}
//
//	buf, errFirstRow := cmd.UpdateFirstRowInCSV(refactorFile, _struct.GetNameFirstRow())
//	if errFirstRow != nil {
//		LogError("Ошибка при добавлении строки", errFirstRow)
//
//		return nil, errFirstRow
//	}
//
//	return buf, nil
//}

func DoBookCSV(path string, sheetName string) (*bytes.Buffer, error) {
	fileContent, downloadFileErr := downloadFile(path)
	if downloadFileErr != nil {
		LogError("Ошибка при загрузке файла: %v", downloadFileErr)
		return nil, downloadFileErr
	}

	data := GetDataFromBytes(fileContent)

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

	err := ReplaceColumnsOnSheet(data, newXlsxFile, sheetName)
	if err != nil {
		LogError("Ошибка при замене колонок на странице %s: %v", sheetName, err)
		return nil, err
	}

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

func ReplaceColumnsOnSheet(data *excelize.File, newXlsxFile *excelize.File, sheetName string) error {
	cols, sheetErr := data.GetCols(sheetName)
	//rows, _ := data.GetRows(sheetName)

	if sheetErr != nil {
		return sheetErr
	}

	for i, col := range cols {
		for _, colName := range col {
			switch colName {
			case "Unit Code":
				setColumnValues(newXlsxFile, cols[i], "A") //number
			case "Target Unit Price":
				setColumnValues(newXlsxFile, cols[i], "B") //price
			case "Total Area":
				setColumnValues(newXlsxFile, cols[i], "C") //Square
			case "Total Area Sq.ft":
				setColumnValues(newXlsxFile, cols[i], "C") //Square
			case "Unit Type":
				setColumnValues(newXlsxFile, cols[i], "F") //layout
				//replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
			case "Description":
				setColumnValues(newXlsxFile, cols[i], "F") //layout
				//replaceUnitLayoutDescriptionFieldInXLSX(newXlsxFile, 5)
			case "View":
				setColumnValues(newXlsxFile, cols[i], "G") //views
			}
		}
	}
	setColumnValues(newXlsxFile, []string{}, "E") //type
	setColumnValues(newXlsxFile, []string{}, "D") //height

	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
	replaceUnitViewsFieldInXLSX(newXlsxFile, 6)
	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)

	cmd.AddLastRowWithEmptyWord(newXlsxFile)

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

func GetDataFromBytes(fileContent []byte) *excelize.File {
	buffer := bytes.NewBuffer(fileContent)
	data, err := excelize.OpenReader(buffer)
	if err != nil {
		LogError("%v", err)
	}
	return data
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
	sheets := file.GetSheetList()

	for _, sheet := range sheets {
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
					replaceWord := strings.Split(cellValue, "-")
					row[colIndex] = replaceWord[len(replaceWord)-1]

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

func replaceUnitLayoutFieldInXLSX(file *excelize.File, indexOfCell int) error {
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
				if cellValue != "Unit Type" {
					if cellValue == row[colIndex] {
						// Заменяем значение ячейки на замену
						words := strings.Split(row[colIndex], " ")

						for _, num := range words {
							_, errCon := strconv.Atoi(num)
							if errCon == nil {
								row[colIndex] = num + " BR"
							}
						}

						columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
						if err1 != nil {
							return err1
						}

						err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), row[colIndex])
						if err1 != nil {
							return err1
						}
						break
					}
				}
			}

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
					word = strings.ReplaceAll(word, "and", "/")
					word = strings.ReplaceAll(word, "+", "/")
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
