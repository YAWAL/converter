package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"os"

	"bytes"

	"path"

	parser "github.com/YAWAL/converter/csvparser"
	"github.com/YAWAL/converter/policy"
)

const defaultCSVFileName = "List of Actions.xlsx"

func main() {
	var m map[string][][]string
	var err error
	var ext string
	var fileName string
	if len(os.Args) > 1 {
		fileName = os.Args[1]
		ext = path.Ext(fileName)
	} else {
		log.Printf("No argument for filename. Used default filename '%s'", defaultCSVFileName)
		fileName = defaultCSVFileName
		ext = path.Ext(fileName)
	}

	if ext == ".csv" {
		m, err = parser.CSV(fileName, ';')
		if err != nil {
			log.Fatalf("Cannot parse csv: %s", err)
		}
	}
	if ext == ".xlsx" {
		m, err = parser.XLSX(fileName)
		if err != nil {
			log.Fatalf("Cannot parse xlsx: %s", err)
		}
	}
	for k := range m {
		//channel for getting Policies from parser.Parse
		readerChan := make(chan policy.Policy, 4)
		go parser.Parse(m[k], readerChan)
		//if directory already exists we get error, but we need just skip this action, not panic
		if err := os.Mkdir(k, os.ModePerm); err != nil && !os.IsExist(err) {
			log.Fatalf("Cannot create directory for policies: %s", err)
		}
		for c := range readerChan {
			marshaledPolicies, err := json.Marshal(&c)
			if err != nil {
				log.Fatalf("Cannot marshal policy '%s' : %s", c.Name, err)
			}
			newName := ReplaceRuneWith(c.Name, ':', '_')
			if err = ioutil.WriteFile(k+"/"+newName+".json", marshaledPolicies, 0666); err != nil {
				log.Fatalf("Cannot save json file for policy '%s': %s", c.Name, err)
			}

		}
	}

	log.Printf("Successfully parsed and saved")
}

//ReplaceRuneWith - return copy of string with changed rune1 to rune2
func ReplaceRuneWith(str string, char1, char2 rune) string {
	buffer := bytes.Buffer{}
	for _, c := range str {
		if c == char1 {
			buffer.WriteRune(char2)
		} else {
			buffer.WriteRune(c)
		}
	}
	return buffer.String()
}
