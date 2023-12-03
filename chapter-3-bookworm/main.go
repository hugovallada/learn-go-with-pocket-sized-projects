package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var bookwormFile string
	flag.StringVar(&bookwormFile, "file", "testdata/bookworms.json", "Path for the bookwormFile")
	flag.Parse()
	bookworms, err := loadBookworms(bookwormFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load bookworms: %s\n", err)
		os.Exit(1)
	}
	commonBooks := findCommonBooks(bookworms)
	fmt.Println("Here are the common books:")
	displayBooks(commonBooks)
}

func displayBooks(books []Book) {
	for _, book := range books {
		fmt.Printf("- %s by %s\n", book.Title, book.Author)
	}
}
