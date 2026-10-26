package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		if scanner.Scan() {
			inputArray := cleanInput(scanner.Text())
			fmt.Printf("Your command was: %s\n", inputArray[0])
		}
	}
}
