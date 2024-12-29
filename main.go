package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"os/signal"
	"io/ioutil"
	"os/user"
)

var history []string

// loadHistory loads the history from the ~/.ccsh_history file.
func loadHistory() {
	usr, _ := user.Current()
	historyFile := usr.HomeDir + "/.ccsh_history"
	if data, err := ioutil.ReadFile(historyFile); err == nil {
		history = strings.Split(string(data), "\n")
	}
}

// saveHistory saves the history to the ~/.ccsh_history file.
func saveHistory() {
	usr, _ := user.Current()
	historyFile := usr.HomeDir + "/.ccsh_history"
	ioutil.WriteFile(historyFile, []byte(strings.Join(history, "\n")), 0644)
}

// clearHistory clears the command history in memory and on disk.
func clearHistory() {
	history = []string{}
	usr, _ := user.Current()
	historyFile := usr.HomeDir + "/.ccsh_history"
	ioutil.WriteFile(historyFile, []byte(""), 0644)
}

// showHistory displays the command history.
func showHistory() {
	for i, command := range history {
		fmt.Printf("%d: %s\n", i+1, command)
	}
}

// handleSignal captures an interrupt signal (Ctrl+C) and ensures a graceful shutdown.
func handleSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for range signalChan {
			fmt.Println("\nExiting shell... Goodbye!")
			saveHistory()
			os.Exit(0)
		}
	}()
}

func main() {
	loadHistory() // Load history when the shell starts
	defer saveHistory() // Save history on exit

	handleSignal() // Handle Ctrl+C gracefully

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("ccsh> ")
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		// Command parsing and exit condition
		if command == "exit" {
			break
		}

		// Handle special commands like 'cd' and 'pwd'
		args := strings.Fields(command)
		if args[0] == "cd" {
			if len(args) < 2 {
				fmt.Println("Usage: cd <directory>")
			} else if err := os.Chdir(args[1]); err != nil {
				fmt.Println("Error:", err)
			}
			continue
		} else if args[0] == "pwd" {
			if dir, err := os.Getwd(); err == nil {
				fmt.Println(dir)
			} else {
				fmt.Println("Error:", err)
			}
			continue
		} else if args[0] == "clear" {
			clearHistory()
			fmt.Println("History cleared!")
			continue
		} else if args[0] == "history" {
			showHistory()
			continue
		}

		// Regular command execution using exec.Command
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command and append to history if it succeeds
		if err := cmd.Run(); err != nil {
			fmt.Println("Error:", err)
		} else {
			history = append(history, command)
		}
	}
}
