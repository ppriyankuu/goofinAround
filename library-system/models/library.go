package models

import "fmt"

type Library struct {
	Books []Book
}

type LibraryActions interface {
	BorrowBook(book *Book, user *User) error
	ReturnBook(book *Book, user *User) error
}

func (l *Library) BorrowBook(book *Book, user *User) error {
	if !book.IsAvailable {
		return fmt.Errorf("Book %s is not available", book.Title)
	}

	book.IsAvailable = false
	user.BorrowedBooks = append(user.BorrowedBooks, *book)

	fmt.Printf("%s borrowed the book %s", user.Name, book.Title)
	return nil
}

func (l *Library) ReturnBook(book *Book, user *User) error {
	for i, b := range user.BorrowedBooks {
		if b.Title == book.Title {
			user.BorrowedBooks = append(user.BorrowedBooks[:i], user.BorrowedBooks[i+1:]...)
			book.IsAvailable = true
			fmt.Printf("%s returned the book %s", user.Name, book.Title)
			return nil
		}
	}
	return fmt.Errorf("book '%s' was not borrowed by %s", book.Title, user.Name)
}
