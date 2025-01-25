package main

func main() {
	// time.Sleep(10 * time.Second)
	go consumeTasks()
	scheduleTaskProduction()
}
