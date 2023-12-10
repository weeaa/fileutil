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
	"os"
	"reflect"
	"strings"
	"unicode"
)

const (
	folderPerm = 0777
	filePerm   = 0666
)

var ()

func ReadCSV[T any](filePath string) ([]T, error) {
	var err error
	var rows []T
	
	body, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	if err = csvutil.Unmarshal(body, &rows); err != nil {
		return nil, err
	}
	
	return rows, nil
}

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
			fmt.Println(field, valueFields[i])
			return fmt.Errorf("type mismatch")
		}
	}
	
	return nil
}

func ReadJSON[T any](filePath string) (T, error) {
	var data T
	
	file, err := os.ReadFile(filePath)
	if err != nil {
		return data, err
	}
	
	if err = json.Unmarshal(file, &data); err != nil {
		return data, err
	}
	
	return data, nil
}

func AppendJSON[T any](filePath string, data T) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	_ = file
	
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}
	
	return nil
}

func ValidateJSON(filePath string, data any) error {
	body, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	valueOf := reflect.ValueOf(data)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	typeOf := valueOf.Type()
	
	var fileData map[string]interface{}
	if err = json.Unmarshal(body, &fileData); err != nil {
		return err
	}
	
	if reflect.TypeOf(fileData) != typeOf {
		return errors.New("error validating: data mismatch [json]")
	}
	
	return nil
}

func WriteJSON(filePath string, data any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}
	return os.WriteFile(filePath, buf.Bytes(), 0777)
}

func CreateCSV(filePath string, keys [][]string) error {
	if _, err := os.Stat(filePath); err != nil {
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

func CreateFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(folderPath, folderPerm)
	}
	return nil
}

func CreateFile(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		var file *os.File
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, filePerm)
		if err != nil {
			return err
		}
		return file.Close()
	}
	return nil
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

func ValidateYAML(filePath string, data any) error {
	body, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	valueOf := reflect.ValueOf(data)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	typeOf := valueOf.Type()
	
	var fileData map[string]interface{}
	if err = yaml.Unmarshal(body, &fileData); err != nil {
		return err
	}
	
	if reflect.TypeOf(fileData) != typeOf {
		return errors.New("error validating: data mismatch [yaml]")
	}
	
	return nil
}

func CreateLogFile(filePath string) {
	os.Create(filePath)
}
