package main

import (
	"fmt"
	"library-system/helpers"
	"library-system/models"
)

func main() {

	// Initializing Books
	books := []models.Book{
		{Title: "The Go Programming Language", Author: "Alan A. A. Donovan", IsAvailable: true},
		{Title: "Clean Code", Author: "Robert C. Martin", IsAvailable: true},
		{Title: "Design Patterns", Author: "Erich Gamma", IsAvailable: true},
	}

	// Initializing Library
	library := &models.Library{Books: books}

	// Initializing Users
	user01 := &models.User{Name: "Cranky"}
	user02 := &models.User{Name: "Gogoi"}

	if err := library.BorrowBook(&library.Books[0], user01); err != nil {
		fmt.Println("\n", err)
	}

	if err := library.BorrowBook(&library.Books[0], user02); err != nil {
		fmt.Println("\n", err)
	}

	if err := library.BorrowBook(&library.Books[2], user02); err != nil {
		fmt.Println("\n", err)
	}

	if err := library.ReturnBook(&library.Books[0], user01); err != nil {
		fmt.Println("\n", err)
	}

	if err := library.ReturnBook(&library.Books[0], user02); err != nil {
		fmt.Println("\n", err)
	}

	if err := library.ReturnBook(&library.Books[1], user01); err != nil {
		fmt.Println("\n", err)
	}

	fmt.Println("\nFinal Status:")
	fmt.Println("Library Books")
	for _, book := range library.Books {
		fmt.Printf("Title: %s, Available: %t\n", book.Title, book.IsAvailable)
	}

	fmt.Println("\nUsers:")
	fmt.Printf("%s borrowed books: %v\n", user01.Name, helpers.GetBorrowedBookTitles(user01))
	fmt.Printf("%s borrowed books: %v\n", user02.Name, helpers.GetBorrowedBookTitles(user02))
}
