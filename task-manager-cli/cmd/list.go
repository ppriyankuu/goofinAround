package cmd

import (
	"fmt"
	"task-manager/internals/tasks"
)

func ListTasks() {
	taskList, err := tasks.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	if len(taskList) == 0 {
		fmt.Println("No tasks found!")
		return
	}

	fmt.Println("Tasks:")
	for i, task := range taskList {
		fmt.Printf("%d. %s\n", i+1, task)
	}
}
