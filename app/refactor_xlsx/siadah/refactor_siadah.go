package siadah

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

func DoBookCSV(path string) (*bytes.Buffer, error) {
	fileContent, downloadFileErr := downloadFile(path)
	if downloadFileErr != nil {
		LogError("Ошибка при загрузке файла: %v", downloadFileErr)
		return nil, downloadFileErr
	}

	sheetName := "Luma22"

	data := GetDataFromBytes(fileContent)

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

	setColumnValues(newXlsxFile, cols[2], "A")    //number
	setColumnValues(newXlsxFile, cols[7], "B")    //price
	setColumnValues(newXlsxFile, cols[4], "C")    //Square
	setColumnValues(newXlsxFile, []string{}, "D") //height
	setColumnValues(newXlsxFile, []string{}, "E") //type
	setColumnValues(newXlsxFile, cols[3], "F")    //layout
	setColumnValues(newXlsxFile, cols[8], "G")    //views

	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
	replaceUnitHeightFieldInXLSX(newXlsxFile, 3)
	replaceUnitTypeFieldInXLSX(newXlsxFile, 4)
	cmd.AddLastRowWithEmptyWord(newXlsxFile)

	buffer, err3 := cmd.ConvertXlsxToCsv(newXlsxFile)

	if err3 != nil {
		LogError("Ошибка при конвертации XLSX в CSV: %v", err3)
		return nil, err3
	}

	refactorFile, errRefactor := cmd.RemoveAnyRowFromCSV(buffer, 4)

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

					if strings.Contains(cellValue, "simplex") {
						if i >= 2 {
							row[i-2] = "Simplex"
							cellValue = strings.ReplaceAll(cellValue, "simplex", "")
						}
					}
					if strings.Contains(cellValue, "duplex") /*|| strings.Contains(cellValue, "loft")*/ {
						if i >= 2 {
							row[i-2] = "Duplex"
							cellValue = strings.ReplaceAll(cellValue, "duplex", "")
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

					arrayLayout := strings.Split(cellValue, " ")
					for _, arrayValue := range arrayLayout {
						valueInt, errAtoi := strconv.Atoi(arrayValue)
						if errAtoi == nil {
							row[colIndex] = strconv.Itoa(valueInt) + " BR"
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
					wordArraay := strings.Split(cellValue, "-")
					row[colIndex] = wordArraay[len(wordArraay)-1]

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
