package tiger

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
//	cols, sheetErr := data.GetCols(sheetName)
//	rows, _ := data.GetRows(sheetName)
//
//	if sheetErr != nil {
//		LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
//		return nil, sheetErr
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
//	for _, row := range rows {
//		for i, cellValue := range row {
//			if strings.Contains(cellValue, "Unit") {
//				setColumnValues(newXlsxFile, cols[i], "A") //number
//			}
//			if strings.Contains(cellValue, "Price") {
//				setColumnValues(newXlsxFile, cols[i], "B") //price
//			}
//			if strings.Contains(cellValue, "Net(sqft)") {
//				setColumnValues(newXlsxFile, cols[i], "C") //Square
//			}
//			setColumnValues(newXlsxFile, []string{}, "D") //height
//			setColumnValues(newXlsxFile, []string{}, "E") //type
//			if strings.Contains(cellValue, "Type") {
//				setColumnValues(newXlsxFile, cols[i], "F") //layout
//			}
//			setColumnValues(newXlsxFile, []string{}, "G") //views
//		}
//	}
//
//	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
//	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
//	replaceUnitHeightFieldInXLSX(newXlsxFile, 3)
//
//	buffer, err3 := cmd.ConvertXlsxToCsv(newXlsxFile)
//
//	if err3 != nil {
//		LogError("Ошибка при конвертации XLSX в CSV: %v", err3)
//		return nil, err3
//	}
//
//	buf, errFirstRow := cmd.UpdateFirstRowInCSV(buffer, _struct.GetNameFirstRow())
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

	for _, row := range rows {
		for i, cellValue := range row {
			if strings.Contains(cellValue, "Unit") {
				setColumnValues(newXlsxFile, cols[i], "A") //number
			}
			if strings.Contains(cellValue, "Price") {
				setColumnValues(newXlsxFile, cols[i], "B") //price
			}
			if strings.Contains(cellValue, "Net(sqft)") {
				setColumnValues(newXlsxFile, cols[i], "C") //Square
			}
			setColumnValues(newXlsxFile, []string{}, "D") //height
			setColumnValues(newXlsxFile, []string{}, "E") //type
			if strings.Contains(cellValue, "Type") {
				setColumnValues(newXlsxFile, cols[i], "F") //layout
			}
			setColumnValues(newXlsxFile, []string{}, "G") //views
		}
	}

	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
	replaceUnitTypeFieldInXLSX(newXlsxFile, 4)
	replaceUnitHeightFieldInXLSX(newXlsxFile, 3)

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

		rowIndex := 1

		for rows.Next() {
			row, err2 := rows.Columns()
			if err2 != nil {
				return err2
			}

			colIndex := indexOfCell

			for _, cellValue := range row {
				if cellValue == row[colIndex] {
					arrayWords := strings.Split(cellValue, "-")
					for _, word := range arrayWords {
						replaceWord, errAtoi := strconv.Atoi(word)
						if errAtoi == nil {
							row[colIndex] = strconv.Itoa(replaceWord)
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
					if strings.Contains(cellValue, "duplex") {
						if i >= 2 {
							row[i-2] = "Duplex"
							cellValue = strings.ReplaceAll(cellValue, "duplex", "")
						}
					}
					if strings.Contains(cellValue, "simplex") {
						if i >= 2 {
							row[i-2] = "Simplex"
							cellValue = strings.ReplaceAll(cellValue, "simplex", "")
						}
					}
					if strings.Contains(cellValue, "triplex") {
						if i >= 2 {
							row[i-2] = "Triplex"
							cellValue = strings.ReplaceAll(cellValue, "triplex", "")
						}
					}
					if strings.Contains(cellValue, "quadruplex") {
						if i >= 2 {
							row[i-2] = "Quadruplex"
							cellValue = strings.ReplaceAll(cellValue, "quadruplex", "")
						}
					}

					switch cellValue {
					case "fivebedroom":
						row[colIndex] = "5 BR"
					case "fourbedroom":
						row[colIndex] = "4 BR"
					case "threebedroom":
						row[colIndex] = "3 BR"
					case "twobedroom":
						row[colIndex] = "2 BR"
					case "onebedroom":
						row[colIndex] = "1 BR"
					}

					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), row[colIndex])
					if err1 != nil {
						return err1
					}

					columnName1, err3 := excelize.ColumnNumberToName(i - 1)
					if err1 != nil {
						return err1
					}

					err3 = file.SetCellValue(sheet, columnName1+strconv.Itoa(rowIndex), row[i-2])
					if err3 != nil {
						return err3
					}
					break
				}
			}
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
					if cellValue == "" {
						row[colIndex] = "Simplex"
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
