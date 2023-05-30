package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
)

func LoadEnv() {
	realPath, errGetWd := os.Getwd()

	if errGetWd != nil {
		panic(errGetWd)
	}

	if err := godotenv.Load(realPath + "/.env"); err != nil {
		panic(err)
	}
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
