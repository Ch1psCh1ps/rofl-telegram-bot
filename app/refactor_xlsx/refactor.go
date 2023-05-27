package refactor_xlsx

import (
	"bytes"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//func DoBookXlsx(path string) {
//	fileContent, downloadFileErr := downloadFile(path)
//	if downloadFileErr != nil {
//		LogError("Ошибка при загрузке файла: %v", downloadFileErr)
//		return
//	}
//
//	data := GetDataFromBytes(fileContent)
//	cols := GetSheet(data)
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
//	setColumnValues(newXlsxFile, cols[1], "A")
//	setColumnValues(newXlsxFile, cols[9], "B")
//	setColumnValues(newXlsxFile, cols[8], "C")
//	setColumnValues(newXlsxFile, cols[2], "D")
//	setColumnValues(newXlsxFile, cols[4], "E")
//	setColumnValues(newXlsxFile, cols[10], "F")
//
//	if err := newXlsxFile.SaveAs("refactor_available_units.xlsx"); err != nil {
//		LogError("%v", err)
//	}
//
//	log.Println("Файл успешно обработан")
//}

func DoBookXlsx(path string) (*bytes.Buffer, error) {
	fileContent, downloadFileErr := downloadFile(path)
	if downloadFileErr != nil {
		LogError("Ошибка при загрузке файла: %v", downloadFileErr)
		return nil, downloadFileErr
	}

	data := GetDataFromBytes(fileContent)
	cols := GetSheet(data)

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

	setColumnValues(newXlsxFile, cols[1], "A")
	setColumnValues(newXlsxFile, cols[9], "B")
	setColumnValues(newXlsxFile, cols[8], "C")
	setColumnValues(newXlsxFile, cols[2], "D")
	setColumnValues(newXlsxFile, cols[4], "E")
	setColumnValues(newXlsxFile, cols[10], "F")

	buffer := new(bytes.Buffer)
	if err := newXlsxFile.Write(buffer); err != nil {
		LogError("%v", err)
		return nil, err
	}

	return buffer, nil
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

func GetSheet(data *excelize.File) [][]string {
	cols, err := data.GetCols("Лист1")
	if err != nil {
		LogError("%v", err)
		return [][]string{}
	}
	return cols
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

func OpenNonexistentFileForTest() {
	file, err := os.Open("nonexistent.txt")
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()
}

func OpenNonexistentFileForTestWithErr() error {
	file, err := os.Open("nonexistent.txt")
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
		return err
	}
	defer file.Close()
	return nil
}
