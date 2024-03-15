package files_util

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jszwec/csvutil"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"unicode"
)

const (
	folderPerm    = 0777
	filePerm      = 0777
	DefaultCSVSep = ","
)

func CreateFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(folderPath, folderPerm)
	}
	return nil
}

func CreateFile(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		var file *os.File
		file, err = os.Create(filePath)
		if err != nil {
			return err
		}
		return file.Close()
	}
	return nil
}

// CreateCSV creates a new CSV file at the given filePath and writes the keys to it.
// It returns an error if there was a problem creating or writing to the file.
func CreateCSV(filePath string, keys [][]string) error {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		var file *os.File

		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}

		if err = csv.NewWriter(file).WriteAll(keys); err != nil {
			return err
		}

		return file.Close()
	}

	return nil
}

// ReadCSV reads a CSV file located at the given filePath and unmarshal its contents into a slice of type T.
// The function returns the unmarshalled rows and any error encountered during the operation.
func ReadCSV[T any](filePath string) ([]T, error) {
	var err error
	var rows []T

	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = csvutil.Unmarshal(body, &rows)
	return rows, err
}

func WriteCSV(filePath string, rows [][]string) error {
	for _, row := range rows {
		if err := AppendCSV(filePath, row); err != nil {
			return err
		}
	}
	return nil
}

// AppendCSV appends records to a CSV file located at the given filePath.
func AppendCSV(filePath string, records []string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	if err = w.Write(records); err != nil {
		return err
	}

	w.Flush()

	return f.Close()
}

func tmp(filePath, fieldValue string) error {
	// Open CSV file
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read File into a variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	// Loop through lines & turn them into object
	for i, line := range lines {
		if line[0] == fieldValue {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}

	// Open CSV file in WRITE mode
	f, err = os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write new datafile
	err = csv.NewWriter(f).WriteAll(lines)
	if err != nil {
		return err
	}

	fmt.Println("Removed row from CSV for ", fieldValue)

	return nil
}

func RemoveCSVRow[T any](filePath string, key any) error {
	rows, err := ReadCSV[T](filePath)
	if err != nil {
		return err
	}

	log.Println("les rows", rows)
	for i, row := range rows {
		if strings.Contains(fmt.Sprint(row), fmt.Sprint(key)) {
			rows = append(rows[:i], rows[i+1:]...)
		}
	}

	log.Println("new rows", rows)

	var rowsStr [][]string
	for _, row := range rows {
		rowsStr = append(rowsStr, strings.Split(fmt.Sprint(row), ","))
	}

	log.Println("rowsStr", rowsStr)
	if err = WriteCSV(filePath, rowsStr); err != nil {
		log.Fatal(err)
	}
	fmt.Println("before nrows")
	nrows, err := ReadCSV[T](filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("nrows", nrows)

	return nil
}

/*
func RemoveCSVRowN[T any](filePath string, index int) error {
	rows, err := ReadCSV[T](filePath)
	if err != nil {
		return err
	}

	for i, row := range rows {
		_ = row
		if i == index {

		}
	}

	return nil
}
*/

// func RemoveCSVRowContains(key string) {}

/*
func RemoveCSVRowsContains[T any](filePath, key string) error {
	rows, err := ReadCSV[T](filePath)
	if err != nil {
		return err
	}

	for _, row := range rows {
		_ = row
	}

	return nil
}
*/

// ValidateCSV validates the fields of a CSV file against the fields of a struct.
func ValidateCSV[T any](filePath, sep string) error {
	var data T

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	if !fileScanner.Scan() {
		return fmt.Errorf("error empty file: %s", filePath)
	}

	fields := strings.Split(fileScanner.Text(), sep)
	for i, field := range fields {
		r := []rune(field)
		r[0] = unicode.ToUpper(r[0])
		fields[i] = string(r)
	}

	valueOf := reflect.ValueOf(data)

	var valueFields []string
	for i := 0; i < valueOf.NumField(); i++ {
		valueFields = append(valueFields, valueOf.Type().Field(i).Name)
	}

	for i, field := range fields {
		if field != valueFields[i] {
			return fmt.Errorf("field name mismatch")
		}
	}

	return nil
}

// ReadJSON reads a JSON file located at the given filePath and unmarshals its contents into a value of type T.
// The function returns the unmarshaled data and any error encountered during the operation.
func ReadJSON[T any](filePath string) (T, error) {
	var data T

	file, err := os.ReadFile(filePath)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(file, &data)
	return data, err
}

func AppendJSON[T any](filePath string, data T) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return WriteJSON(filePath, buf)
}

func ValidateJSON[T any](filePath string) error {
	var data T

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&data); err != nil {
		return fmt.Errorf("error while decoding JSON: %w", err)
	}

	jsonFields := reflect.ValueOf(data).Elem().Type()

	structFields := reflect.TypeOf(new(T)).Elem()

	for i := 0; i < structFields.NumField(); i++ {
		structField := structFields.Field(i)
		structFieldName := structField.Name

		jsonFieldName := structFieldName
		if tag, ok := structField.Tag.Lookup("json"); ok {
			jsonFieldName = tag
		}

		jsonFieldNameRunes := []rune(jsonFieldName)
		if len(jsonFieldNameRunes) > 0 {
			jsonFieldNameRunes[0] = unicode.ToLower(jsonFieldNameRunes[0])
		}

		jsonField, ok := jsonFields.FieldByName(string(jsonFieldNameRunes))

		if !ok {
			return fmt.Errorf("json field '%s' not found", structFieldName)
		}

		if jsonField.Type != structField.Type {
			return fmt.Errorf("field type mismatch '%s': JSON has '%s', struct has '%s'",
				structFieldName, jsonField.Type, structField.Type)
		}
	}

	return nil
}

func WriteJSON(filePath string, data any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), filePerm)
}

// RemoveJSONLine removes a specific line in JSON
func RemoveJSONLine(filePath, key string) error {
	content, err := ReadJSON[map[string]any](filePath)
	if err != nil {
		return err
	}

	for k := range content {
		if k == key {
			delete(content, k)
		}
	}

	return WriteJSON(filePath, content)
}

func CreateYAML[T any](filePath string, dataEncoded T) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	err = yaml.NewEncoder(file).Encode(&dataEncoded)
	return err
}

func ReadYAML[T any](filePath string) (T, error) {
	var data T

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return data, err
	}

	err = yaml.NewDecoder(file).Decode(&data)
	return data, err
}

func AppendYAML[T any](filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	_ = file
	return nil
}

func ValidateYAML[T any](filePath string) error {
	var data T

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&data); err != nil {
		return fmt.Errorf("error while decoding YAML: %w", err)
	}

	yamlFields := reflect.ValueOf(data).Elem().Type()

	structFields := reflect.TypeOf(new(T)).Elem()

	for i := 0; i < structFields.NumField(); i++ {
		structField := structFields.Field(i)
		structFieldName := structField.Name

		yamlFieldName := structFieldName
		if tag, ok := structField.Tag.Lookup("yaml"); ok {
			yamlFieldName = tag
		}

		yamlFieldNameRunes := []rune(yamlFieldName)
		if len(yamlFieldNameRunes) > 0 {
			yamlFieldNameRunes[0] = unicode.ToLower(yamlFieldNameRunes[0])
		}

		yamlField, ok := yamlFields.FieldByName(string(yamlFieldNameRunes))

		if !ok {
			return fmt.Errorf("yaml field '%s' not found", structFieldName)
		}

		if yamlField.Type != structField.Type {
			return fmt.Errorf("field type mismatch '%s': YAML has '%s', struct has '%s'",
				structFieldName, yamlField.Type, structField.Type)
		}
	}

	return nil
}

func UnmarshalJSONToStruct[T any](data io.ReadCloser) (T, error) {
	var t T
	b, err := io.ReadAll(data)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(b, &t)
	return t, err
}

func ReadProxyFile(path string) (proxies []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		proxies = append(proxies, fileScanner.Text())
	}

	for i := range proxies {
		r := rand.Intn(i + 1)
		proxies[i], proxies[r] = proxies[r], proxies[i]
	}

	if len(proxies) == 0 {
		return nil, errors.New("empty proxy list")
	}

	return proxies, nil
}
