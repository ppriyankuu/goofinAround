package cmd

import (
	"fmt"
	"strconv"
	"task-manager/internals/tasks"
)

func DeleteTask(args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide the task number to delete")
		return
	}

	taskIndex, err := strconv.Atoi(args[0])
	if err != nil || taskIndex <= 0 {
		fmt.Println("Invalid task number")
		return
	}

	err = tasks.DeleteTask(taskIndex - 1)
	if err != nil {
		fmt.Printf("Error deleting task: %v\n", err)
		return
	}

	fmt.Printf("Task %d deleted successfully\n", taskIndex)
}
