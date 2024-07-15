package utilities

import "strings"

// CountWords counts the number of words in a text
func CountWords(text string) int {
	return len(strings.Fields(text))
}

// EstimateReadingTime estimates the reading time in minutes based on the word count
func EstimateReadingTime(wordCount int) int {
	readingSpeed := 200 // words per minute
	return (wordCount + readingSpeed - 1) / readingSpeed
}
