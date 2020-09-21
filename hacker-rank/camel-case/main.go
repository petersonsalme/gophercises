package main

import (
	"fmt"
	"strings"
)

func main() {
	var input string
	fmt.Scanf("%s", &input)

	answer := 1
	for _, ch := range input {
		str := string(ch)

		if strings.ToUpper(str) == str {
			answer++
		}
	}

	fmt.Println(answer)
}
