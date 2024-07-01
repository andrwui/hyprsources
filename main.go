package main

import (
	"fmt"

	"github.com/andrwui/hyprsources/helper"
	"github.com/andrwui/hyprsources/models"
)

func init() {
	fmt.Print("\033[H\033[2J")
}

func main() {

	jsonSources := models.GetJsonSources()
	hyprSources := helper.GetConfigSources()

	sources := models.MatchSources(jsonSources, hyprSources)

	sources.SaveToJson()

	list := models.NewList("hyprsources - Manage hyprland source files.\nPress [h] for help.\n").SetSources(sources)

	editedSources, ok := list.Display()

	if ok {
		editedSources.SaveToJson()
		editedSources.SaveToHyprland()
		fmt.Print("Changes saved.")
	} else {
		fmt.Print("Exiting...")
	}

}
