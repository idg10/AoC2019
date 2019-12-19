package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], " <day> <inputFolder>")
		fmt.Println("E.g.: ", os.Args[0], " 1 IdgInput")
		return
	}

	day, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("First argument must be an integer")
		return
	}

	var inputFolder string
	inputFolder = os.Args[2]

	switch day {
	case 1:
		showDay1(inputFolder)
	case 2:
		showDay2(inputFolder)
	default:
		fmt.Println("No solution for day ", day)
	}
}
