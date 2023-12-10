package files_util

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

const (
	tempFilePathCSV  = "temp.csv"
	tempFilePathJSON = "temp.json"
)

var expectedValuesCSV = []string{"bob", "foo", "43", "FR"}

type csvContent struct {
	FirstName  string `csv:"firstName"`
	LastName   string `csv:"lastName"`
	Age        int    `csv:"age"`
	CountryISO string `csv:"countryISO"`
}

type jsonContent struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Age        int    `json:"age"`
	CountryISO string `json:"countryISO"`
}

func TestCSV(t *testing.T) {
	err := CreateCSV(tempFilePathCSV, [][]string{
		{
			"firstName",
			"lastName",
			"age",
			"countryISO",
		},
	})
	assert.NoError(t, err)
	
	assert.FileExists(t, tempFilePathCSV)
	
	defer os.Remove(tempFilePathCSV)
	assert.NoError(t, AppendCSV(tempFilePathCSV, expectedValuesCSV))
	assert.NoError(t, ValidateCSV[csvContent](tempFilePathCSV, ","))
	
	rows, err := ReadCSV[csvContent](tempFilePathCSV)
	assert.NoError(t, err)
	
	assert.Contains(t, rows[0].CountryISO, "FR")
	assert.Contains(t, rows[0].FirstName, "bob")
	assert.Contains(t, rows[0].LastName, "foo")
	
}

func TestJSON(t *testing.T) {
	var err error
	
	values := jsonContent{
		FirstName:  expectedValuesCSV[0],
		LastName:   expectedValuesCSV[1],
		Age:        43,
		CountryISO: expectedValuesCSV[3],
	}
	
	if err = WriteJSON(tempFilePathJSON, values); err != nil {
		assert.Error(t, err)
	}
	defer os.Remove(tempFilePathJSON)
	
	data, err := ReadJSON[jsonContent](tempFilePathJSON)
	if !assert.Error(t, err) {
		t.Fatal(err)
	}
	
	if !assert.Contains(t, data.FirstName, "bob") {
		t.Fatalf("data.FirstName mismatch, expected bob, got %s", data.FirstName)
	}
}

func TestYAML(t *testing.T) {
	var err error
	var buf bytes.Buffer
	
	values := map[string]any{
		"foo": true,
		"bar": 61.77,
		"baz": "hey",
	}
	
	if err = yaml.NewEncoder(&buf).Encode(values); err != nil {
		assert.Error(t, err)
	}
	
}

func TestCreateFolder(t *testing.T) {
	tempFolderPath := "tempfolder"
	
	CreateFolder(tempFolderPath)
	defer os.Remove(tempFolderPath)
	
	_, err := os.Stat(tempFolderPath)
	if err != nil {
		assert.Errorf(t, err, "failed to create test folder")
	}
	
	assert.NoError(t, nil)
}

func TestCreateFile(t *testing.T) {
	tempFilePath := "tempfile.txt"
	
	CreateFile(tempFilePath)
	defer os.Remove(tempFilePath)
	
	_, err := os.Stat(tempFilePath)
	if err != nil {
		assert.Errorf(t, err, "failed to create test file")
	}
	
	assert.NoError(t, nil)
}
