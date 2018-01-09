package csvparser

import (
	"encoding/csv"
	"errors"
	"os"
	"path"
	"strings"
)

func CSV(fileName string, delimiter rune) (map[string][][]string, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return make(map[string][][]string), err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return make(map[string][][]string), err
	}
	if len(records) < 1 {
		return make(map[string][][]string), errors.New("0 records founded")
	}
	baseName := path.Base(fileName)
	name := strings.Split(baseName, ".")
	res := make(map[string][][]string)
	res[name[0]] = records
	return res, nil
}
