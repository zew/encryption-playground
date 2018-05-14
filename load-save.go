package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/zew/util"
)

// Load a Set from JSON file
func Load(fn ...string) *Set {

	setFile := "inp.json"
	if len(fn) > 0 {
		setFile = fn[0]
	}

	file, err := os.Open(setFile)
	if err != nil {
		log.Fatalf("Could not load set file %v: %v", setFile, err)
	}
	defer file.Close()
	log.Printf("Found set file: %v", setFile)

	decoder := json.NewDecoder(file)
	tmpSet := Set{}
	err = decoder.Decode(&tmpSet)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Set loaded 1\n%#s", util.IndentedDump(tmpSet))
	return &tmpSet
}

// Save a Set to JSON file
func (s *Set) Save(fn ...string) error {

	firstColLeftMostPrefix := " "
	byts, err := json.MarshalIndent(s, firstColLeftMostPrefix, "\t")
	if err != nil {
		return err
	}

	saveDir := "."
	setFile := "out.json"
	if len(fn) > 0 {
		setFile = fn[0]
	}
	savePath := path.Join(saveDir, setFile)
	err = ioutil.WriteFile(savePath, byts, 0644)
	if err != nil {
		return err
	}

	log.Printf("Saved Set to %v", savePath)
	return nil
}
