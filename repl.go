package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/awbalessa/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	next     *string
	previous *string
}

type Repl struct {
	config *config
	client *pokeapi.PokeClient
}

var commands map[string]cliCommand

func Init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Provides the next 20 location areas",
			callback:    commandMap,
		},
	}

}
func cleanInput(str string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(str)))
}

func (r *Repl) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // So the compiler doesn't complain
}

func (r *Repl) commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func (r *Repl) commandMap() error {
	areas, err := r.client.GetLocationAreas(r.config.next)
	if err != nil {
		return err
	}
	r.config.next = areas.Next
	for _, result := range areas.Results {
		fmt.Println(result.Name)
	}
	return nil
}
