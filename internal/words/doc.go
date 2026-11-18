// Package words holds the built-in word lists, one list per file.
//
// Every list is sorted, has no repeats, and uses only lower-case letters
// a-z and digits; the slug package's tests enforce that. To change a list,
// edit its file. To add one, add a file here and a function in
// ../../words.go that returns a copy of it.
package words
