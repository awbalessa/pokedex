package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		args := strings.Fields(input)
		if cmd, exists := commands[args[0]]; exists {
			err := cmd.callback(args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
