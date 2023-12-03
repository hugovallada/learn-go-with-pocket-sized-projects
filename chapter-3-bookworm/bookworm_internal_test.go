package main

import "testing"

var (
	handmaidsTale = Book{Author: "Margaret Atwood", Title: "The Handmaid's Tale"}
	oryxAndCrake  = Book{Author: "Margaret Atwood", Title: "Oryx and Crake"}
	theBellJar    = Book{Author: "Sylvia Plath", Title: "The Bell Jar"}
	janeEyre      = Book{Author: "Charlotte Brontë", Title: "Jane Eyre"}
)

func TestLoadBookworms(t *testing.T) {
	type testCase struct {
		bookwormsFile string
		want          []Bookworm
		wantErr       bool
	}

	testCases := map[string]testCase{
		"file exists": {
			bookwormsFile: "testdata/bookworms.json",
			want: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			wantErr: false,
		},
		"file doesn't exist": {
			bookwormsFile: "testdata/no_file_here.json",
			want:          nil,
			wantErr:       true,
		},
		"invalid JSON": {
			bookwormsFile: "testdata/invalid.json",
			want:          nil,
			wantErr:       true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			got, err := loadBookworms(testCase.bookwormsFile)
			if !testCase.wantErr && err != nil {
				t.Fatalf("expected no error, got one %s", err.Error())
			}
			if testCase.wantErr && err == nil {
				t.Fatalf("expected an error %s, got none", err.Error())
			}
			if !equalBookworms(got, testCase.want, t) {
				t.Fatalf("different result: got %s, expected %s", got, testCase.want)
			}
		})
	}
}

func equalBookworms(bookworms, target []Bookworm, t *testing.T) bool {
	t.Helper()
	if len(bookworms) != len(target) {
		return false
	}
	for i := range bookworms {
		if bookworms[i].Name != target[i].Name {
			return false
		}

		if !equalBooks(bookworms[i].Books, target[i].Books, t) {
			return false
		}
	}
	return true
}

func equalBooks(books, target []Book, t *testing.T) bool {
	t.Helper()
	if len(books) != len(target) {
		return false
	}
	for i := range books {
		if books[i] != target[i] {
			return false
		}
	}
	return true
}

func equalBooksCount(got, want map[Book]uint, t *testing.T) bool {
	t.Helper()
	if len(got) != len(want) {
		return false
	}
	for book, targetCount := range want {
		count, ok := got[book]
		if !ok || targetCount != count {
			return false
		}
	}
	return true
}

func TestBookCount(t *testing.T) {
	tt := map[string]struct {
		input []Bookworm
		want  map[Book]uint
	}{
		"nominal use case": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			want: map[Book]uint{handmaidsTale: 2, theBellJar: 1, oryxAndCrake: 1, janeEyre: 1},
		},
		"no bookworms": {
			input: []Bookworm{},
			want:  map[Book]uint{},
		},
		"bookworm without books": {
			input: []Bookworm{
				{Name: "Fadi"},
				{Name: "Peggy"},
			},
		},
		"bookworm with twice the same book": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, handmaidsTale}},
			},
			want: map[Book]uint{handmaidsTale: 2},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := booksCount(tc.input)
			if !equalBooksCount(got, tc.want, t) {
				t.Fatalf("got a different list of books: %v, expected %v", got, tc.want)
			}
		})
	}
}

func TestFindCommonBooks(t *testing.T) {
	tt := map[string]struct {
		input []Bookworm
		want  []Book
	}{

		"no common book": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, janeEyre}},
			},
			want: nil,
		},
		"one common book": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale}},
			},
			want: []Book{handmaidsTale},
		},
		"three bookworms have the same books on their shelves": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale}},
				{Name: "Faruha", Books: []Book{handmaidsTale}},
			},
			want: []Book{handmaidsTale},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := findCommonBooks(tc.input)
			if !equalBooks(tc.want, got, t) {
				t.Fatalf("got a different list of books %v, expected %v", got, tc.want)
			}
		})
	}
}
