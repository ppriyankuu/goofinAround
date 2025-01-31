package utils

import "log"

func LogInfo(message string) {
	log.Println("[Info]", message)
}

func LogError(message string) {
	log.Println("[Error]", message)
}
