package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetUserString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s", prompt)

	res, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	res = strings.TrimSpace(res)

	return res
}
