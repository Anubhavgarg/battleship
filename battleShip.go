package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func getGridSizeFunction(fileReader *bufio.Reader) int {
	size, error := fileReader.ReadString('\n')
	if error != nil {
		log.Println(error)
		os.Exit(-1)
	}

	number, error := strconv.ParseInt(strings.TrimSpace(size), 10, 0)
	if error != nil {
		log.Println(error)
		os.Exit(-1)
	}

	return int(number)
}

func main() {
	if len(os.Args) != 3 {
		os.Exit(-1)
	}
	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer outputFile.Close()
	fileReader := bufio.NewReader(inputFile)
	getSize := getGridSizeFunction(fileReader)
	println(getSize)
}
