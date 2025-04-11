package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/awbalessa/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

type config struct {
	next     *string
	previous *string
	pokedex  map[string]pokeapi.Pokemon
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
			pokedex:  make(map[string]pokeapi.Pokemon),
		},
		client: &client,
	}
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    func(args []string) error { return r.commandExit() },
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    func(args []string) error { return r.commandHelp() },
		},
		"map": {
			name:        "map",
			description: "Provides the next 20 location areas",
			callback:    func(args []string) error { return r.commandMap() },
		},
		"mapb": {
			name:        "mapb",
			description: "Provides the previous 20 location areas",
			callback:    func(args []string) error { return r.commandMapb() },
		},
		"explore": {
			name:        "explore",
			description: "Explores the Pokemon in the area",
			callback:    func(args []string) error { return r.commandExplore(args[1]) },
		},
		"catch": {
			name:        "catch",
			description: "Throws a Pokeball at a Pokemon",
			callback:    func(args []string) error { return r.commandCatch(args[1]) },
		},
		"inspect": {
			name:        "inspect",
			description: "Returns information about a Pokemon",
			callback:    func(args []string) error { return r.commandInspect(args[1]) },
		},
		"pokedex": {
			name:        "pokedex",
			description: "Returns all the Pokemon in your Pokedex",
			callback:    func(args []string) error { return r.commandPokedex() },
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

func (r *repl) commandExplore(name string) error {
	areaPokemon, err := r.client.ExploreArea(name)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...", name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range areaPokemon.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}

func (r *repl) commandCatch(name string) error {
	pokemon, err := r.client.GetPokemon(name)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	difficulty := max(0, (pokemon.BaseExperience / 34))
	if attempt := rand.Intn(difficulty); attempt == 0 {
		fmt.Printf("%s was caught!\n", name)
		r.config.pokedex[name] = *pokemon
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}

func (r *repl) commandInspect(name string) error {
	pokeinfo, exist := r.config.pokedex[name]
	if !exist {
		return fmt.Errorf("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Height: %d\n", pokeinfo.Height)
	fmt.Printf("Weight: %d\n", pokeinfo.Weight)
	fmt.Printf("Stats:")
	for _, stat := range pokeinfo.Stats {
		fmt.Printf("	-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:")
	for _, poketype := range pokeinfo.Types {
		fmt.Printf("	-%s\n", poketype.Type.Name)
	}
	return nil
}

func (r *repl) commandPokedex() error {
	if len(r.config.pokedex) == 0 {
		return fmt.Errorf("your pokedex is empty")
	}
	for pokemon := range r.config.pokedex {
		fmt.Println(pokemon)
	}
	return nil
}
