package main

import (
	"fmt"
	"os"
	"task-manager/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', or 'delete' sub-commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cmd.AddTask(os.Args[2:])
	case "list":
		cmd.ListTasks()
	case "delete":
		cmd.DeleteTask(os.Args[2:])
	default:
		fmt.Printf("Unknown sub-command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
