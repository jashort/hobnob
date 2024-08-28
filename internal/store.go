package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Action struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Data      string    `json:"data"`
}

type Alias struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type Datafile struct {
	Actions []Action `json:"actions"`
	Version int      `json:"version"`
}

type Data struct {
	Aliases []Alias
	Notes   []Note
	Actions []Action
	Version int
}

func AddAlias(data *Data, alias Alias) error {
	for _, a := range data.Aliases {
		if strings.ToLower(a.Alias) == strings.ToLower(alias.Alias) {
			return errors.New(fmt.Sprintf("The alias %s already exists for %s", a.Alias, a.Name))
		}
	}
	data.Aliases = append(data.Aliases, alias)
	jsonRecord, _ := json.Marshal(alias)
	data.Actions = append(
		data.Actions, Action{
			Timestamp: time.Now(),
			Action:    "addAlias",
			Data:      string(jsonRecord),
		})
	return nil
}

func AddNote(data *Data, note Note) {
	data.Notes = append(data.Notes, note)
	jsonRecord, _ := json.Marshal(note)
	data.Actions = append(data.Actions, Action{
		Timestamp: time.Now(),
		Action:    "addNote",
		Data:      string(jsonRecord),
	})
}

func (data *Data) LookupName(name string) string {
	for _, a := range data.Aliases {
		if strings.ToLower(a.Alias) == strings.ToLower(name) {
			return a.Name
		}
	}
	return name
}

func Save(filename string, actions []Action) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	data := Datafile{
		Actions: actions,
		Version: 1,
	}
	asJson, err := json.MarshalIndent(data, "", " ")
	_, err = f.Write(asJson)
	if err != nil {
		return err
	}
	return nil
}

func LoadAll(filename string) (Data, error) {
	rawData, err := Load(filename)
	if err != nil {
		panic(err)
	}
	data := Data{
		Aliases: []Alias{},
		Notes:   []Note{},
		Actions: rawData.Actions,
		Version: rawData.Version,
	}

	for _, action := range rawData.Actions {
		if action.Action == "addAlias" {
			var dat Alias
			if err := json.Unmarshal([]byte(action.Data), &dat); err != nil {
				panic(err)
			}
			data.Aliases = append(data.Aliases, dat)
		} else if action.Action == "addNote" {
			var dat Note
			if err := json.Unmarshal([]byte(action.Data), &dat); err != nil {
				panic(err)
			}
			data.Notes = append(data.Notes, dat)
		}
	}

	return data, nil
}

func Load(filename string) (Datafile, error) {
	dataFile, err := os.Open(filename)
	if err != nil {
		emptyDatafile := Datafile{
			Actions: []Action{},
			Version: 1,
		}
		return emptyDatafile, nil
	}
	defer func(dataFile *os.File) {
		err := dataFile.Close()
		if err != nil {
			panic(err)
		}
	}(dataFile)

	jsonParser := json.NewDecoder(dataFile)
	var data Datafile
	if err = jsonParser.Decode(&data); err != nil {
		panic(err)
	}
	return data, nil
}
