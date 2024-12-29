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

func loadHistory() {
	usr, _ := user.Current()
	historyFile := usr.HomeDir + "/.ccsh_history"
	if data, err := ioutil.ReadFile(historyFile); err == nil {
		history = strings.Split(string(data), "\n")
	}
}

func saveHistory() {
	usr, _ := user.Current()
	historyFile := usr.HomeDir + "/.ccsh_history"
	ioutil.WriteFile(historyFile, []byte(strings.Join(history, "\n")), 0644)
}

func main() {

	loadHistory()
	defer saveHistory()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for range signalChan {
			fmt.Println("\nccsh> ")
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("ccsh> ")
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		if command == "exit" {
			break
		}

		args := strings.Fields(command)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

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
		}

		if err := cmd.Run(); err != nil {
			fmt.Println("No such file or directory:", err)
		}

		history = append(history, command)
	}
}

