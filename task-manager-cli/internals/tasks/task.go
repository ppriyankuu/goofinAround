package tasks

import (
	"encoding/json"
	"errors"
	"os"
)

const filepath = "data/tasks.json"

func SaveTask(task string) error {
	taskList, err := LoadTasks()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	taskList = append(taskList, task)

	data, err := json.Marshal(taskList)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func LoadTasks() ([]string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var taskList []string
	err = json.Unmarshal(data, &taskList)
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func DeleteTask(index int) error {
	taskList, err := LoadTasks()
	if err != nil {
		return err
	}

	if index < 0 || index >= len(taskList) {
		return errors.New("task not found")
	}

	taskList = append(taskList[:index], taskList[index+1:]...)
	data, err := json.Marshal(taskList)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}
