package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func createShipMatrix(fileReader *bufio.Reader, size int) [][]byte {

	shipInfo := make([][]byte, size)
	for i := range shipInfo {
		shipInfo[i] = make([]byte, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			shipInfo[i][j] = '_'
		}
	}

	line, error := fileReader.ReadString('\n')
	if error != nil {
		log.Println(error)
		os.Exit(-1)
	}

	for _, xyCoordinate := range strings.Split(strings.TrimSpace(line), ",") {
		xyCoordinateData := strings.Split(xyCoordinate, ":")
		x, _ := strconv.ParseInt(xyCoordinateData[0], 10, 0)
		y, _ := strconv.ParseInt(xyCoordinateData[1], 10, 0)
		shipInfo[x][y] = 'B'
	}

	return shipInfo
}

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

type missileInfoDimensions struct {
	x, y int
}

func putMiscile(fileReader *bufio.Reader, data chan<- missileInfoDimensions) {
	line, error := fileReader.ReadString('\n')
	if error != nil {
		log.Println(error)
		os.Exit(-1)
	}

	for _, xyCoordinate := range strings.Split(strings.TrimSpace(line), ":") {
		xyCoordinateData := strings.Split(xyCoordinate, ",")
		x, _ := strconv.ParseInt(xyCoordinateData[0], 10, 0)
		y, _ := strconv.ParseInt(xyCoordinateData[1], 10, 0)
		data <- missileInfoDimensions{x: int(x), y: int(y)}
	}
	close(data)
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
	if getSize < 0 || getSize >= 10 {
		fmt.Printf("Grid Size should between 0 and 10")
		os.Exit(-1)
	}
	numberOfShips := getGridSizeFunction(fileReader)
	println(numberOfShips, getSize)
	player1ShipMatrix := createShipMatrix(fileReader, getSize)
	player2ShipMatrix := createShipMatrix(fileReader, getSize)
	missilesTotal := getGridSizeFunction(fileReader)
	player1moves := make(chan missileInfoDimensions)
	player2moves := make(chan missileInfoDimensions)
	go putMiscile(fileReader, player1moves)
	go putMiscile(fileReader, player2moves)
}
