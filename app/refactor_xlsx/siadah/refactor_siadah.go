package siadah

import (
	"bytes"
	"genieMap/cmd"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

	return refactorFile, nil
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
