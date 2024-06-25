package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/andrwui/hyprsources/constants"
	"github.com/andrwui/hyprsources/helper"
)

type SourceFile struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Selected bool   `json:"selected"`
	Priority string `json:"priority"`
}

type SourceFilePointerArray []*SourceFile

func (s *SourceFilePointerArray) SaveToHyprland() {

	homedir := helper.GetHomedir()

	hyprlandConfigDir := homedir + constants.HYPRLAND_CONFIG

	configFileLines, err := helper.GetFileLines(hyprlandConfigDir)
	if err != nil {
		fmt.Println("Could not load hyprland config:")
		panic(err)
	}

	var lowPrioritySources []string
	var highPrioritySources []string

	for _, sourceFile := range *s {

		var lineToAppend string

		if !sourceFile.Selected {
			lineToAppend += "# "
		}

		lineToAppend += "source = " + sourceFile.Path

		switch sourceFile.Priority {
		case "low":
			lowPrioritySources = append(lowPrioritySources, lineToAppend)
		case "high":
			highPrioritySources = append(highPrioritySources, lineToAppend)
		}
	}

	var filteredConfigLines []string
	existingSourcesPattern := "source = "
	for _, line := range configFileLines {
		if !strings.Contains(line, existingSourcesPattern) {
			filteredConfigLines = append(filteredConfigLines, line)
		}
	}

	filteredConfigLines = append(filteredConfigLines, highPrioritySources...)

	filteredConfigLines = append(lowPrioritySources, filteredConfigLines...)

	finalFile := strings.Join(filteredConfigLines, "\n")

	err = os.WriteFile(hyprlandConfigDir, []byte(finalFile), fs.ModePerm)
	if err != nil {
		fmt.Println("Could not save hyprland config:")
		panic(err)
	}

}

func (s *SourceFilePointerArray) SaveToJson() {

	homedir := helper.GetHomedir()

	json, err := json.Marshal(s)
	if err != nil {
		log.Fatal("Could not save sources.json")
		panic(err)
	}

	os.WriteFile(homedir+constants.JSON_SOURCES, json, os.ModePerm)

}

func unmarshalSources(file *os.File) []SourceFile {
	configBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var sources []SourceFile

	json.Unmarshal(configBytes, &sources)

	return sources

}

func GetJsonSources() []SourceFile {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic("Could not access home directory.")
	}

	configFilePath := homedir + "/.config/hyprsources/sources.json"

	if _, err := os.Stat(configFilePath); err == nil {
		file, err := os.Open(configFilePath)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		return unmarshalSources(file)

	} else if errors.Is(err, os.ErrNotExist) {

		err := os.MkdirAll(homedir+"/.config/hyprsources/", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err2 := os.Create(configFilePath)
		if err2 != nil {
			panic(err2)
		}

		defer file.Close()
		return unmarshalSources(file)

	} else {
		panic(err)
	}
}
