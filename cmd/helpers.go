package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func LoadEnv() {
	realPath, errGetWd := os.Getwd()

	if errGetWd != nil {
		panic(errGetWd)
	}

	godotenv.Load(realPath + "/.env")
}

func ConvertXlsxToCsv(inputFile *excelize.File) (*bytes.Buffer, error) {
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

func RemoveFirstRowFromXLSX(file *excelize.File, sheetName string) (*excelize.File, error) {
	// Get all the rows from the sheet
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) <= 1 {
		// The file doesn't have enough rows to remove the first row
		return nil, fmt.Errorf("файл не содержит достаточно строк для удаления первой строки")
	}

	// Create a new XLSX file
	newFile := excelize.NewFile()

	// Copy rows, excluding the first row, to the new sheet
	for rowIndex, row := range rows {
		if rowIndex > 0 {
			for colIndex, cellValue := range row {
				cellName, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex)
				newFile.SetCellValue(sheetName, cellName, cellValue)
			}
		}
	}

	return newFile, nil
}

func RemoveFirstRowFromCSV(buffer *bytes.Buffer) (*bytes.Buffer, error) {
	// Читаем содержимое буфера в строку
	content := buffer.String()

	// Разделяем строки по символу новой строки
	lines := strings.Split(content, "\n")

	if len(lines) <= 1 {
		// Файл не содержит достаточно строк для удаления первой строки
		return nil, fmt.Errorf("файл не содержит достаточно строк для удаления первой строки")
	}

	// Удаляем первую строку из среза строк
	lines = lines[1:]

	// Собираем строки обратно в одну строку
	newContent := strings.Join(lines, "\n")

	// Создаем новый буфер и записываем в него обновленное содержимое
	newBuffer := bytes.NewBufferString(newContent)

	return newBuffer, nil
}

func RemoveAnyRowFromCSV(buffer *bytes.Buffer, rowForRemove int) (*bytes.Buffer, error) {
	// Читаем содержимое буфера в строку
	content := buffer.String()

	// Разделяем строки по символу новой строки
	lines := strings.Split(content, "\n")

	if len(lines) <= 1 {
		// Файл не содержит достаточно строк для удаления первой строки
		return nil, fmt.Errorf("файл не содержит достаточно строк для удаления первой строки")
	}

	// Удаляем первую строку из среза строк
	lines = lines[rowForRemove:]

	// Собираем строки обратно в одну строку
	newContent := strings.Join(lines, "\n")

	// Создаем новый буфер и записываем в него обновленное содержимое
	newBuffer := bytes.NewBufferString(newContent)

	return newBuffer, nil
}

func ReplaceFieldInXLSX(file *excelize.File) error {
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

			// Перебираем каждую ячейку в строке
			for colIndex, cellValue := range row {
				if cellValue == row[colIndex] {
					// Заменяем значение ячейки на замену
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

func ReplaceWhateverFieldInXLSX(file *excelize.File, indexOfCell int) error {
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
					word := strings.Replace(row[colIndex], "HK", "R", 1)
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

func RoundString(input string) (string, error) {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return input, err
	}

	rounded := fmt.Sprintf("%.2f", value)
	return rounded, nil
}

func UpdateFirstRowInCSV(buffer *bytes.Buffer, values []string) (*bytes.Buffer, error) {
	content, err := ioutil.ReadAll(buffer)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(content))
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		rows[0] = values
	}

	newBuffer := new(bytes.Buffer)
	writer := csv.NewWriter(newBuffer)

	for _, row := range rows {
		err1 := writer.Write(row)
		if err1 != nil {
			return nil, err1
		}
	}

	writer.Flush()

	return newBuffer, nil
}

func AddEmptyFirstLine(buffer *bytes.Buffer) (*bytes.Buffer, error) {
	result := bytes.NewBuffer(nil)

	_, err := fmt.Fprintln(result, "empty")

	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteTo(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
