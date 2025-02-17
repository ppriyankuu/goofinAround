package models

type Room struct {
	ID       string
	Users    map[string]bool
	MaxUsers int
}
