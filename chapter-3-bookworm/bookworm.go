package main

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
)

type Bookworm struct {
	Name  string `json:"name"`
	Books []Book `json:"books"`
}

type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}

// loadBookworms reads the file and returns the list of bookworms, and their beloved books, found therein.
func loadBookworms(filePath string) ([]Bookworm, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buffedReader := bufio.NewReaderSize(f, 1024*1024)
	var booksworms []Bookworm
	err = json.NewDecoder(buffedReader).Decode(&booksworms)
	if err != nil {
		return nil, err
	}
	return booksworms, nil
}

// findCommonBooks returns books thar are on more than one bookworm
func findCommonBooks(bookwoms []Bookworm) []Book {
	var commonBooks []Book
	booksOnShelves := booksCount(bookwoms)
	for book, count := range booksOnShelves {
		if count > 1 {
			commonBooks = append(commonBooks, book)
		}
	}
	return sortBooks(commonBooks)
}

// booksCount registers all the books and their occurrences from the bookworms shelves.
func booksCount(bookworms []Bookworm) map[Book]uint {
	counter := make(map[Book]uint)
	for _, bookworm := range bookworms {
		for _, book := range bookworm.Books {
			counter[book]++
		}
	}
	return counter
}

type byAuthor []Book

func (b byAuthor) Len() int { return len(b) }

func (b byAuthor) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byAuthor) Less(i, j int) bool {
	if b[i].Author != b[j].Author {
		return b[i].Author < b[j].Author
	}
	return b[i].Title < b[j].Title
}

func sortBooks(books []Book) []Book {
	sort.Sort(byAuthor(books))
	return books
}
