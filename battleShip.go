package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <input-file> <output-file>\n", os.Args[0])
		os.Exit(-1)
	}

}
