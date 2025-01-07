package cmd

import (
	"fmt"
	"strings"
	"task-manager/internals/tasks"
)

func AddTask(args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a task description")
		return
	}

	task := strings.Join(args, " ")
	err := tasks.SaveTask(task)
	if err != nil {
		fmt.Printf("Error adding task: %v\n", err)
		return
	}

	fmt.Printf("Task added: %s\n", task)
}
