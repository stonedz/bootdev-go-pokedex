package main

import "fmt"
import "bufio"
import "os"

type cliCommand struct {
	command     string
	description string
	callback    func() error
}

func main() {
	commands := make(map[string]cliCommand)
	addCommand(commands, "help", "Prints the help message", commandHelp)
	addCommand(commands, "exit", "Exits the program", commandExit)

	for {
		fmt.Print("pokedex >")
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		scanner.Scan()
		text := scanner.Text()
		handleCommand(text, commands)
	}

}

func handleCommand(text string, commands map[string]cliCommand) {
	command, ok := commands[text]
	if ok {
		command.callback()
	} else {
		fmt.Println("Command not found")
	}

}

func addCommand(commands map[string]cliCommand, command string, description string, callback func() error) {

	commands[command] = cliCommand{command, description, callback}
}

func commandHelp() error {
	fmt.Println("Help message")
	return nil
}
func commandExit() error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}
