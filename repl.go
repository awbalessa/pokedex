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
	callback    func() error
}

type config struct {
	next     *string
	previous *string
}

type repl struct {
	config *config
	client *pokeapi.PokeClient
}

var commands map[string]cliCommand

func Init() *repl {
	client := pokeapi.NewClient()
	r := &repl{
		config: &config{
			next:     nil,
			previous: nil,
		},
		client: &client,
	}
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    func() error { return r.commandExit() },
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    func() error { return r.commandHelp() },
		},
		"map": {
			name:        "map",
			description: "Provides the next 20 location areas",
			callback:    func() error { return r.commandMap() },
		},
		"mapb": {
			name:        "mapb",
			description: "Provides the previous 20 location areas",
			callback:    func() error { return r.commandMapb() },
		},
	}
	return r
}
func cleanInput(str string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(str)))
}

func (r *repl) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // So the compiler doesn't complain
}

func (r *repl) commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func (r *repl) commandMap() error {
	initialCall := r.config.next == nil && r.config.previous == nil

	if r.config.next == nil && !initialCall {
		fmt.Println("You've reached the end of the available locations.")
		return nil
	}
	areas, err := r.client.GetLocationAreas(r.config.next)
	if err != nil {
		return err
	}
	r.config.next = &areas.Next
	r.config.previous = &areas.Previous
	for _, result := range areas.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func (r *repl) commandMapb() error {
	if r.config.previous == nil {
		fmt.Println("You're on the first page.")
		return nil
	}
	areas, err := r.client.GetLocationAreas(r.config.previous)
	if err != nil {
		return err
	}
	r.config.next = &areas.Next
	r.config.previous = &areas.Previous
	for _, result := range areas.Results {
		fmt.Println(result.Name)
	}
	return nil
}
