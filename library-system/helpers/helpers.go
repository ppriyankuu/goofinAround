package helpers

import "library-system/models"

func GetBorrowedBookTitles(user *models.User) []string {
	titles := []string{}
	for _, book := range user.BorrowedBooks {
		titles = append(titles, book.Title)
	}
	return titles
}
