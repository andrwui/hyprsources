package helper

import (
	"bufio"
	"os"
)

func GetFileLines(path string) ([]string, error) {

	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil

}
