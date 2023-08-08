package azizi

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

	buf, errFirstRow := cmd.UpdateFirstRowInCSV(buffer, _struct.GetNameFirstRow())
	if errFirstRow != nil {
		LogError("Ошибка при добавлении строки", errFirstRow)
		return nil, errFirstRow
	}

	return buf, nil
}

func ReplaceColumnsOnSheet(data *excelize.File, newXlsxFile *excelize.File, sheetName string) error {
	cols, sheetErr := data.GetCols(sheetName)
	rows, _ := data.GetRows(sheetName)

	if sheetErr != nil {
		return sheetErr
	}

	for i, col := range cols {
		for _, colName := range col {
			colName = strings.ReplaceAll(colName, " ", "")
			colName = strings.ToLower(colName)
			switch colName {
			case "unit":
				setColumnValues(newXlsxFile, cols[i], "A") //Square
			case "Unit":
				setColumnValues(newXlsxFile, cols[i], "A") //Square
			}
		}
	}

	for _, row := range rows {
		for i, cellValue := range row {
			cellValue = strings.ReplaceAll(cellValue, " ", "")
			cellValue = strings.ToLower(cellValue)
			if strings.Contains(cellValue, "price") {
				setColumnValues(newXlsxFile, cols[i], "B") //price
			}
			if strings.Contains(cellValue, "totalarea(sqft)") {
				setColumnValues(newXlsxFile, cols[i], "C") //Square
			}
			setColumnValues(newXlsxFile, []string{}, "D") //height
			setColumnValues(newXlsxFile, []string{}, "E") //type
			if strings.Contains(cellValue, "unittype") {
				setColumnValues(newXlsxFile, cols[i], "F") //layout
			}
			if strings.Contains(cellValue, "view") {
				setColumnValues(newXlsxFile, cols[i], "G") //views
			}
		}
	}

	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
	replaceUnitTypeFieldInXLSX(newXlsxFile, 4)
	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
	replaceUnitHeightFieldInXLSX(newXlsxFile, 3)
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

			for i, cellValue := range row {
				if cellValue == row[colIndex] {
					cellValue = strings.ToLower(cellValue)
					row[colIndex] = strings.ReplaceAll(cellValue, "\nbedroom", " BR")

					if strings.Contains(cellValue, "shop") {
						row[colIndex] = "retail"
						if i >= 1 {
							row[i-1] = "commerce"
							cellValue = strings.ReplaceAll(cellValue, "commerce", "")
						}
						if i >= 2 {
							row[i-2] = "shop"
							cellValue = strings.ReplaceAll(cellValue, "shop", "")
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

					columnName1, err3 := excelize.ColumnNumberToName(i)
					if err1 != nil {
						return err1
					}

					err3 = file.SetCellValue(sheet, columnName1+strconv.Itoa(rowIndex), row[i-1])
					if err3 != nil {
						return err3
					}

					columnName2, err4 := excelize.ColumnNumberToName(i - 1)
					if err1 != nil {
						return err1
					}

					err4 = file.SetCellValue(sheet, columnName2+strconv.Itoa(rowIndex), row[i-2])
					if err4 != nil {
						return err4
					}
					break
				}
			}
			rowIndex++
		}
	}

	return nil
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
					replaceWordArray := strings.Split(cellValue, " ")
					for _, replaceWord := range replaceWordArray {
						convertWord, errAtoi := strconv.Atoi(replaceWord)
						if errAtoi == nil {
							row[colIndex] = strconv.Itoa(convertWord)
						}
					}

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

func replaceUnitTypeFieldInXLSX(file *excelize.File, indexOfCell int) error {
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
					cellValue = strings.ToLower(cellValue)
					switch cellValue {
					case "apartment":
						row[colIndex] = "Apartments"
					case "commerce":
						row[colIndex] = "commerce"
					}

					if cellValue == "" {
						row[colIndex] = "Apartments"
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
			rowIndex++
		}
	}

	return nil
}

func replaceUnitHeightFieldInXLSX(file *excelize.File, indexOfCell int) error {
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
					cellValue = strings.ToLower(cellValue)

					switch cellValue {
					case "":
						row[colIndex] = "Simplex"
					case "shop":
						row[colIndex] = ""
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
			rowIndex++
		}
	}

	return nil
}
