package main

import "math"

type Recommendation struct {
	Book  Book
	Score float64
}

type Reader struct {
	Name  string
	Books []Book
}

type set map[Book]struct{}

func (s set) Contains(b Book) bool {
	_, ok := s[b]
	return ok
}

func newSet(books ...Book) set {
	setOfBooks := make(map[Book]struct{}, len(books))
	for _, book := range books {
		setOfBooks[book] = struct{}{}
	}
	return setOfBooks
}

func recommend(allReaders []Reader, target Reader, n int) []Recommendation {
	read := newSet(target.Books...)
	recommendations := map[Book]float64{}
	for _, reader := range allReaders {
		if reader.Name == target.Name {
			continue
		}
		var similarity float64
		for _, book := range reader.Books {
			if read.Contains(book) {
				similarity++
			}
		}
		if similarity == 0 {
			continue
		}
		score := math.Log(similarity) + 1
		for _, book := range reader.Books {
			if !read.Contains(book) {
				recommendations[book] += score
			}
		}
	}
	return NewListOfRecommendationsOf(recommendations)
}

func NewListOfRecommendationsOf(recommendations map[Book]float64) []Recommendation {
	var recommends []Recommendation
	for book, score := range recommendations {
		recommends = append(recommends, Recommendation{
			Book:  book,
			Score: score,
		})
	}
	return recommends
}
