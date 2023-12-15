package files_util

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"log"
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
	tests := []struct {
		name        string
		filePath    string
		expectedStr any
	}{
		{
			name:        "success",
			filePath:    tempFilePathCSV,
			expectedStr: csvContent{},
		},
		// todo add fail test
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
			assert.NoError(t, AppendCSV(tempFilePathCSV, []string{"quoicoubeh", "q", "22", "q"}))
			assert.NoError(t, RemoveCSVRow[csvContent](tempFilePathCSV, "foo"))
			assert.NoError(t, ValidateCSV[csvContent](tempFilePathCSV, ","))

			rows, err := ReadCSV[csvContent](tempFilePathCSV)
			log.Println("ROWS", rows)
			assert.NoError(t, err)

			assert.Contains(t, rows[0].CountryISO, "FR")
			assert.Contains(t, rows[0].FirstName, "bob")
			assert.Contains(t, rows[0].LastName, "foo")
		})
	}
}

// todo move tests there
func tmpJson(t *testing.T) {

	tests := []struct {
		name string
		pass bool
	}{
		{
			name: "",
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}

// fonctionne
func TestJSON(t *testing.T) {
	var err error

	type idiom struct {
		city string `json:"city"`
		id   int    `json:"id"`
	}

	type nested struct {
		name  string `json:"name"`
		idiom idiom  `json:"idiom"`
	}

	values := jsonContent{
		FirstName:  expectedValuesCSV[0],
		LastName:   expectedValuesCSV[1],
		Age:        43,
		CountryISO: expectedValuesCSV[3],
	}

	if err = WriteJSON(tempFilePathJSON, values); err != nil {
		assert.Error(t, err)
	}

	r, err := ReadJSON[jsonContent](tempFilePathJSON)
	if err != nil {
		t.Fatal(err)
	}

	if err = RemoveJSONLine(tempFilePathJSON, "firstName"); err != nil {
		t.Fatal(err)
	}

	r, err = ReadJSON[jsonContent](tempFilePathJSON)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("new rows", r)
	defer os.Remove(tempFilePathJSON)

	data, err := ReadJSON[jsonContent](tempFilePathJSON)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	if !assert.Contains(t, data.LastName, "foo") {
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

	if err := CreateFile(tempFilePath); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFilePath)

	_, err := os.Stat(tempFilePath)
	if err != nil {
		assert.Errorf(t, err, "failed to create test file")
	}

	assert.NoError(t, nil)
}
