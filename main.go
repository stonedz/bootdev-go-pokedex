package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/stonedz/bootdev-go-pokedex/internal/cache"
)

type cliCommand struct {
	command     string
	description string
	callback    func(conf *config, cache *cache.Cache) error
}

type config struct {
	Next string
	Prev string
}

type mapLocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func main() {
	conf := config{Next: "", Prev: ""}
	cache := cache.NewCache(10 * time.Second)
	commands := make(map[string]cliCommand)
	addCommand(commands, "help", "Prints the help message", commandHelp)
	addCommand(commands, "exit", "Exits the program", commandExit)
	addCommand(commands, "map", "Prints the name of 20 pokemon locations", commandMap)
	addCommand(commands, "mapb", "Prints previous 20 pokemon locations", commandMapb)

	for {
		fmt.Print("pokedex >")
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		scanner.Scan()
		text := scanner.Text()
		handleCommand(text, commands, &conf, cache)
	}

}

func handleCommand(text string, commands map[string]cliCommand, conf *config, cache *cache.Cache) {
	command, ok := commands[text]
	if ok {
		err := command.callback(conf, cache)
		if err != nil {
			fmt.Println(err)
			commandExit(conf, cache)
		}
	} else {
		fmt.Println("Command not found")
	}

}

func addCommand(commands map[string]cliCommand, command string, description string, callback func(conf *config, cache *cache.Cache) error) {

	commands[command] = cliCommand{command, description, callback}
}

func commandHelp(conf *config, cache *cache.Cache) error {
	fmt.Println("Help message")
	return nil
}
func commandExit(conf *config, cache *cache.Cache) error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

func commandMap(conf *config, cache *cache.Cache) error {
	req := ""
	if conf.Next != "" {
		req = conf.Next
	} else {
		req = "https://pokeapi.co/api/v2/location/"
	}
	fmt.Println("Map locations...")
	res, err := http.Get(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return errors.New("Response failed!")
	}
	if err != nil {
		return err
	}
	myMap := mapLocations{}
	err = json.Unmarshal(body, &myMap)
	if err != nil {
		return err
	}
	//fmt.Println(myMap)
	conf.Next = myMap.Next
	conf.Prev = myMap.Previous
	for _, v := range myMap.Results {
		fmt.Println(v.Name)
	}

	return nil
}

func commandMapb(conf *config, cache *cache.Cache) error {
	fmt.Println("Mapb...")
	if conf.Prev == "" {
		return errors.New("No previous page!")
	}

	conf.Next = conf.Prev
	commandMap(conf, cache)
	return nil
}
