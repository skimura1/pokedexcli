package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	for {
		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan();
		userInput := scanner.Text()
		cleanedInput := cleanInput(userInput)
		fmt.Println("Your command was:", cleanedInput[0])
	}
}
