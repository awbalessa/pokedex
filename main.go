package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Init()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if success := scanner.Scan(); !success {
			fmt.Println("error reading input")
			break
		}
		input := scanner.Text()
		if cmd, exists := commands[input]; exists {
			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
