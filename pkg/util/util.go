// Package util contains some utilities we will use trough the program.
package util

import (
	"regexp"
	"strconv"
	"unicode"
)

// StringIsArgument will report if a string is a number.
// This is useful because usually assignments that are numbers are actually
// function arguments.
func StringIsArgument(input string) bool {
	_, err := strconv.ParseFloat(input, 64)

	return err == nil
}

// StringIsAllUpper will report if a string is all UPPERCASE.
// This is useful because usually all uppercase variables are env variables.
func StringIsAllUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

// IsNumber will report if a string is a numeric.
func IsNumber(s string) bool {
	return regexp.MustCompile(`\d`).MatchString(s)
}

// UniqueSliceElements will deduplicate a slice.
func UniqueSliceElements[T comparable](inputSlice []T) []T {
	uniqueSlice := make([]T, 0, len(inputSlice))
	seen := make(map[T]bool, len(inputSlice))

	for _, element := range inputSlice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}

	return uniqueSlice
}
