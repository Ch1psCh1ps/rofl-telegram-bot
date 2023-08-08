package MAg

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

const (
	PriceIndex  = 2
	LayoutIndex = 5
)

func DoBookCSV(path string) (*bytes.Buffer, error) {
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

	err := ReplaceColumnsOnSheet(data, newXlsxFile)
	if err != nil {
		LogError("Ошибка при замене колонок%: %v", err)
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

func ReplaceColumnsOnSheet(data *excelize.File, newXlsxFile *excelize.File) error {
	sheetList := data.GetSheetList()
	currentRow := 1
	for _, sheetName := range sheetList {
		cols, sheetErr := data.GetCols(sheetName)

		if sheetErr != nil {
			return sheetErr
		}

		setColumnValues(newXlsxFile, cols[2], "A", &currentRow)    //number 0
		setColumnValues(newXlsxFile, cols[10], "B", &currentRow)   //price 1
		setColumnValues(newXlsxFile, cols[8], "C", &currentRow)    //Square 2
		setColumnValues(newXlsxFile, []string{}, "D", &currentRow) //height 3
		setColumnValues(newXlsxFile, cols[5], "E", &currentRow)    //type 4
		setColumnValues(newXlsxFile, cols[3], "F", &currentRow)    //layout 5
		setColumnValues(newXlsxFile, cols[6], "G", &currentRow)    //views 6

		currentRow += len(cols[1]) // или любой другой столбец
	}
	replaceUnitNumberFieldInXLSX(newXlsxFile, 0)
	replaceUnitTypeFieldInXLSX(newXlsxFile, 4)
	replaceUnitLayoutFieldInXLSX(newXlsxFile, 5)
	replaceUnitViewsFieldInXLSX(newXlsxFile, 6)

	err1 := RemoveEmptyRows(newXlsxFile)

	if err1 != nil {
		fmt.Println(err1)
	}

	cmd.AddLastRowWithEmptyWord(newXlsxFile)

	return nil
}

func RemoveEmptyRows(f *excelize.File) error {
	sheet := f.GetSheetList()[0]

	rows, err := f.GetRows(sheet)
	if err != nil {
		return fmt.Errorf("error reading Excel data: %w", err)
	}

	for i := len(rows) - 1; i >= 0; i-- {
		price := rows[i][PriceIndex-1]
		roomType := rows[i][LayoutIndex-1]

		if price == "" || roomType == "" {
			f.RemoveRow(sheet, i+1)
		}
	}

	if err := f.Save(); err != nil {
		return fmt.Errorf("error saving Excel file: %w", err)
	}

	return nil
}

func setColumnValues(file *excelize.File, values []string, colPrefix string, startRow *int) {
	for i, value := range values {
		cell := colPrefix + strconv.Itoa(*startRow+i)
		file.SetCellValue("Sheet1", cell, value)
	}
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
					arrayLayout := strings.Split(cellValue, " ")
					for _, arrayValue := range arrayLayout {
						valueInt, errAtoi := strconv.Atoi(arrayValue)
						if errAtoi == nil {
							row[colIndex] = strconv.Itoa(valueInt) + " BR"
						}
					}

					replacements := map[string]string{
						"penthouse": "Penthouse",
						"vila":      "Vila",
						"townhouse": "Townhouse",
					}

					for term, replacement := range replacements {
						if strings.Contains(cellValue, term) {
							if i >= 2 {
								row[i-2] = replacement

								prevColName, err1 := excelize.ColumnNumberToName(i - 1)
								if err1 != nil {
									return err1
								}
								err2 = file.SetCellValue(sheet, prevColName+strconv.Itoa(rowIndex), row[i-2])
								if err2 != nil {
									return err2
								}
							}
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
					replaceWordArray := strings.Split(cellValue, " ")

					if len(replaceWordArray) > 1 {
						row[colIndex] = strings.ReplaceAll(replaceWordArray[len(replaceWordArray)-1], "-", "")
					} else {
						arrNum := strings.Split(cellValue, "-")
						row[colIndex] = arrNum[len(arrNum)-1]
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
					cellValue = strings.ToLower(cellValue)
					switch cellValue {
					case "apartment":
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
					row[colIndex] = strings.ReplaceAll(cellValue, "+", "/")

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
