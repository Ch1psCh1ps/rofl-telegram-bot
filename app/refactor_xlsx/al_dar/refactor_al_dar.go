package al_dar

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"genieMap/cmd"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func DoBookCSV(path string) (*bytes.Buffer, error) {
	filecsv, convertErr := ConvertCSVtoXLSX(path)
	if convertErr != nil {
		LogError("Ошибка при конвертации CSV в XLSX: %v", convertErr)
		return nil, convertErr
	}

	//fileContent, downloadFileErr := downloadFile(path)
	//fmt.Println(path)
	//if downloadFileErr != nil {
	//	LogError("Ошибка при загрузке файла: %v", downloadFileErr)
	//	return nil, downloadFileErr
	//}

	//filecsv, convertErr := ConvertCSVtoXLSX(fileContent)
	//if convertErr != nil {
	//	LogError("Ошибка при конвертации CSV в XLSX: %v", convertErr)
	//	return nil, convertErr
	//}
	//
	//// Сохранение файла XLSX
	//saveErr := SaveXLSXFile(filecsv, "output.xlsx")
	//if saveErr != nil {
	//	LogError("Ошибка при сохранении файла XLSX: %v", saveErr)
	//	return nil, saveErr
	//}

	//data := GetDataFromBytes(fileContent)
	//
	//defer func() {
	//	if err := data.Close(); err != nil {
	//		LogError("%v", err)
	//	}
	//}()

	cols, sheetErr := GetSheet(filecsv, "Sheet1")
	if sheetErr != nil {
		LogError("Ошибка при получении листа из файла XLSX: %v", sheetErr)
		return nil, sheetErr
	}

	newXlsxFile := excelize.NewFile()

	defer func() {
		if err := newXlsxFile.Close(); err != nil {
			LogError("%v", err)
		}
	}()

	setColumnValues(newXlsxFile, cols[2], "A")
	setColumnValues(newXlsxFile, cols[11], "B")
	setColumnValues(newXlsxFile, cols[10], "C")
	setColumnValues(newXlsxFile, []string{}, "D")
	setColumnValues(newXlsxFile, cols[4], "E")
	setColumnValues(newXlsxFile, cols[5], "F")
	setColumnValues(newXlsxFile, cols[6], "G")

	cmd.ReplaceFieldInXLSX(newXlsxFile)
	cmd.ReplaceWhateverFieldInXLSX(newXlsxFile, 5)

	buf, err3 := convertXlsxToCsv(newXlsxFile)

	if err3 != nil {
		LogError("Ошибка при конвертации XLSX в CSV: %v", err3)
		return nil, err3
	}

	//buffer := new(bytes.Buffer)
	//if err := newXlsxFile.Write(buffer); err != nil {
	//	LogError("%v", err)
	//	return nil, err
	//}

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

func ConvertCSVtoXLSX(csvURL string) (*excelize.File, error) {
	// Получение данных из ссылки
	response, err := http.Get(csvURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить данные из ссылки: %v", err)
	}

	// Чтение CSV-файла
	reader := csv.NewReader(response.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать CSV-файл: %v", err)
	}

	defer response.Body.Close()

	//// Создание файла на диске
	//file1, err := os.Create("/Users/anastasyaplotnikova/GolandProjects/genieMap/telegram.csv")
	//if err != nil {
	//	return nil, fmt.Errorf("не удалось создать файл: %v", err)
	//}
	//defer file1.Close()
	//
	//// Копирование содержимого из response.Body в файл
	//_, err = io.Copy(file1, response.Body)
	//if err != nil {
	//	return nil, fmt.Errorf("не удалось записать файл: %v", err)
	//}
	//
	//// Открытие существующего файла CSV
	//file, err := os.Open("/Users/anastasyaplotnikova/GolandProjects/genieMap/telegram.csv")
	//if err != nil {
	//	return nil, fmt.Errorf("не удалось открыть файл CSV: %v", err)
	//}
	//defer file.Close()

	// Создание нового файла XLSX
	xlsxFile := excelize.NewFile()

	// Создание нового листа
	sheetName := "Sheet1"

	// Запись данных из CSV-файла в XLSX
	for rowIndex, record := range records {
		for colIndex, cellValue := range record {
			cell, _ := excelize.ColumnNumberToName(colIndex + 1)
			cell += fmt.Sprintf("%d", rowIndex+1)
			xlsxFile.SetCellValue(sheetName, cell, cellValue)
		}
	}

	//xlsxFilePathForSave := "/Users/anastasyaplotnikova/GolandProjects/genieMap/refactor.xlsx"
	//
	////// Сброс указателя позиции чтения/записи в начало файла XLSX
	////xlsxFile.SetActiveSheet(0)
	//
	//// Сохранение файла XLSX на диске
	//if err1 := xlsxFile.SaveAs(xlsxFilePathForSave); err1 != nil {
	//	return nil, fmt.Errorf("не удалось сохранить файл XLSX: %v", err1)
	//}

	return xlsxFile, nil
}

func convertXlsxToCsv(inputFile *excelize.File) (*bytes.Buffer, error) {
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

func DownloadFileFromTelegramBot(token, fileID, filePath string) ([]byte, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", token, fileID)
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о файле: %v", err)
	}
	defer response.Body.Close()

	// Чтение ответа в формате JSON
	// Здесь вам потребуется ваша собственная логика для извлечения информации о файле
	// из ответа Telegram API, чтобы получить ссылку на файл

	fileURL := "https://api.telegram.org/file/bot" + token + "/" + filePath

	// Загрузка файла
	fileResponse, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить файл: %v", err)
	}
	defer fileResponse.Body.Close()

	// Чтение содержимого файла в []byte
	fileContent, err := ioutil.ReadAll(fileResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать содержимое файла: %v", err)
	}

	return fileContent, nil
}

func SaveXLSXFile(file *excelize.File, outputPath string) error {
	if err := file.SaveAs(outputPath); err != nil {
		return fmt.Errorf("не удалось сохранить файл XLSX: %v", err)
	}
	return nil
}

func GetDataFromBytes(fileContent []byte) *excelize.File {
	buffer := bytes.NewBuffer(fileContent)
	data, err := excelize.OpenReader(buffer)
	if err != nil {
		LogError("%v", err)
	}
	return data
}
