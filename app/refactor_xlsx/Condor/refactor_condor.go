package Condor

import (
	"bytes"
	"fmt"
	"genieMap/cmd"
	_struct "genieMap/structures"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"math"
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
	cols, sheetErr := data.GetCols(sheetName)

	if sheetErr != nil {
		LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
		return nil, sheetErr
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

	setColumnValues(newXlsxFile, cols[0], "A")    //number
	setColumnValues(newXlsxFile, cols[8], "B")    //price
	setColumnValues(newXlsxFile, cols[6], "C")    //Square
	setColumnValues(newXlsxFile, []string{}, "D") //height
	setColumnValues(newXlsxFile, cols[3], "E")    //type
	setColumnValues(newXlsxFile, cols[2], "F")    //layout
	setColumnValues(newXlsxFile, cols[4], "G")    //views

	replaceUnitPriceFieldInXLSX(newXlsxFile, 1)
	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
	replaceUnitTypeFieldInXLSX(newXlsxFile, 4)

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

func GetDataFromFile(filePath string) (*excelize.File, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %v", err)
	}

	return file, nil
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

func replaceUnitSquareFieldInXLSX(file *excelize.File, indexOfCell int) error {
	//не забудь поменять значение word!!!!
	// Получаем список имен листов из файла
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
					fmt.Println(cellValue, "cell")
					fmt.Println(row[0], "unit number")
					word := cellValue
					row[colIndex] = word

					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), word)
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

//func replaceUnitPriceFieldInXLSX(file *excelize.File, indexOfCell int) error {
//	//не забудь поменять значение word!!!!
//	// Получаем список имен листов из файла
//	sheets := file.GetSheetList()
//
//	for _, sheet := range sheets {
//		rows, err := file.Rows(sheet)
//		if err != nil {
//			return err
//		}
//
//		rowIndex := 1
//
//		for rows.Next() {
//			row, err2 := rows.Columns()
//			if err2 != nil {
//				return err2
//			}
//
//			colIndex := indexOfCell
//
//			for _, cellValue := range row {
//				if cellValue == row[colIndex] {
//					word1 := strings.Split(cellValue, ".")
//					fmt.Println(word1)
//					word := word1[0]
//					row[colIndex] = word
//
//					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
//					if err1 != nil {
//						return err1
//					}
//
//					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), word)
//					if err1 != nil {
//						return err1
//					}
//					break
//				}
//			}
//			rowIndex++
//		}
//	}
//
//	return nil
//}

func replaceUnitPriceFieldInXLSX(file *excelize.File, indexOfCell int) error {
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

					word := strconv.Itoa(int(parseString(cellValue)))
					row[colIndex] = word

					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), word)
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

			for _, cellValue := range row {
				if cellValue == row[colIndex] {

					wordSplit := strings.Split(cellValue, "")
					replaceWord := wordSplit[0] + " BR"
					row[colIndex] = replaceWord

					columnName, err1 := excelize.ColumnNumberToName(colIndex + 1)
					if err1 != nil {
						return err1
					}

					err1 = file.SetCellValue(sheet, columnName+strconv.Itoa(rowIndex), replaceWord)
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

					row[colIndex] = "apartments"

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

func Round(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return t + math.Copysign(1, x)
	}
	return t
}

func parseString(f string) float64 {
	if s, err := strconv.ParseFloat(f, 64); err == nil {
		s = Round(s)
		return s
	}
	return 0
}
