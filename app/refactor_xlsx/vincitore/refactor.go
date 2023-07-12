package vincitore

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

	addFirstRow, errAddFirstRow := cmd.AddEmptyFirstLine(buffer)

	if errAddFirstRow != nil {
		LogError("Ошибка добавления пустой строки", errAddFirstRow)
	}

	buf, errFirstRow := cmd.UpdateFirstRowInCSV(addFirstRow, _struct.GetNameFirstRow())
	if errFirstRow != nil {
		LogError("Ошибка при добавлении строки", errFirstRow)
		return nil, errFirstRow
	}

	return buf, nil
}

func ReplaceColumnsOnSheet(data *excelize.File, newXlsxFile *excelize.File, sheetName string) error {
	cols, sheetErr := data.GetCols(sheetName)

	if sheetErr != nil {
		return sheetErr
	}

	setColumnValues(newXlsxFile, cols[1], "A")    //number
	setColumnValues(newXlsxFile, cols[5], "B")    //price
	setColumnValues(newXlsxFile, cols[3], "C")    //Square
	setColumnValues(newXlsxFile, []string{}, "D") //height
	setColumnValues(newXlsxFile, []string{}, "E") //type
	setColumnValues(newXlsxFile, cols[2], "F")    //layout
	setColumnValues(newXlsxFile, cols[4], "G")    //views

	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
	replaceUnitViewsFieldInXLSX(newXlsxFile, 6)

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
				if cellValue == row[colIndex] {

					replaceWordArray := strings.Split(cellValue, " ")
					for _, replaceWord := range replaceWordArray {
						convertWord, errAtoi := strconv.Atoi(replaceWord)
						if errAtoi == nil {
							row[colIndex] = strconv.Itoa(convertWord) + " BR"
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
					replaceWord := strings.ReplaceAll(cellValue, "View", "")
					replaceWord = strings.ReplaceAll(replaceWord, "Views", "")
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
				}
			}
			rowIndex++
		}
	}

	return nil
}
