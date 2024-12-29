package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("ccsh> ")
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)

		if command == "exit" {
			break
		}

		cmd := exec.Command(command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println("No such file or directory:", err)
		}
	}
}

