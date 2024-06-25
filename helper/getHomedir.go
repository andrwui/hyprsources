package helper

import "os"

func GetHomedir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic("Could not access home directory.")
	}
	return homedir
}
