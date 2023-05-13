package wildfrostbot

import (
	"encoding/json"
	"os"
	"sort"
)

type Library struct {
	Articles map[string]Embeder
	Commands []string
}

// MakeLibrary initializes the library and reads the data
func MakeLibrary() (*Library, error) {
	var lib = &Library{
		Articles: map[string]Embeder{},
	}

	if err := loadDictionary(); err != nil {
		return nil, err
	}

	if err := lib.MakeArticles(); err != nil {
		return nil, err
	}

	if err := lib.MakeCategories(); err != nil {
		return nil, err
	}

	lib.CreateAliases()

	// Get list of commands for autocomplete
	lib.Commands = make([]string, 0, len(lib.Articles))
	for k := range lib.Articles {
		lib.Commands = append(lib.Commands, k)
	}
	sort.Strings(lib.Commands)

	return lib, nil
}

// loadDictionary loads the dictionary into memory
func loadDictionary() error {
	replacementDict = make(map[string]string)
	fobj, err := os.Open(dictionaryPath)
	if err != nil {
		return err
	}
	return json.NewDecoder(fobj).Decode(&replacementDict)
}
