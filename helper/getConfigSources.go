package helper

import (
	"os"
	"strings"
)

func GetConfigSources() []string {

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic("Could not access home directory.")
	}

	lines, err := GetFileLines(homedir + "/.config/hypr/hyprland.conf")
	if err != nil {
		panic("Could not find hyprland config file. \n Verify your hyprland install.")
	}

	var res []string
	for i := range len(lines) {
		if strings.Contains(strings.Join(strings.Fields(lines[i]), ""), "source=") {
			res = append(res, lines[i])
		}
	}
	return res

}
