package main

import (
	"bufio"
	"fmt"
	"io"
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
func missileLaunchingFunction(ship [][]byte, target missileInfoDimensions) bool {
	switch ship[target.x][target.y] {
	case '_':
		ship[target.x][target.y] = 'O'
	case 'B':
		ship[target.x][target.y] = 'X'
		return true
	}
	return false
}

func matrixPrint(matrix [][]byte, size int, writer io.Writer) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Fprintf(writer, "%c ", matrix[i][j])
		}
		fmt.Fprintf(writer, "\n")
	}
	fmt.Fprintf(writer, "\n")
}

func finalResult(p1Count, p2Count int, writer io.Writer) {
	fmt.Fprintf(writer, "P1:%d\n", p1Count)
	fmt.Fprintf(writer, "P2:%d\n", p2Count)

	switch {
	case p1Count > p2Count:
		fmt.Fprintln(writer, "Player 1 wins")
	case p1Count < p2Count:
		fmt.Fprintln(writer, "Player 2 wins")
	case p1Count == p2Count:
		fmt.Fprintln(writer, "It is a draw")
	}
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
	var player1Hits, player2Hits int
	for i := 0; i < missilesTotal; i++ {
		if missileLaunchingFunction(player1ShipMatrix, <-player1moves) {
			player2Hits++
		}
		if missileLaunchingFunction(player2ShipMatrix, <-player2moves) {
			player1Hits++
		}
	}
	writer := io.MultiWriter(outputFile, os.Stdout)

	fmt.Fprintf(writer, "Player1\n")
	matrixPrint(player1ShipMatrix, getSize, writer)
	fmt.Fprintf(writer, "Player2\n")
	matrixPrint(player2ShipMatrix, getSize, writer)

	finalResult(player1Hits, player2Hits, writer)

}
