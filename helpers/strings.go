/*
* Copyright (c) 2014 Santiago Arias | Remy Jourde
*
* Permission to use, copy, modify, and distribute this software for any
* purpose with or without fee is hereby granted, provided that the above
* copyright notice and this permission notice appear in all copies.
*
* THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
* WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
* MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
* ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
* WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
* ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
* OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package helpers

import (
	"strings"
)

// TrimLower returns a lower case slice of the string s, with all leading and trailing white space removed, as defined by Unicode.
//
func TrimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// SetOfStrings returns a set of string for a given strings.
//
func SetOfStrings(s string) []string {
	slice := strings.Split(TrimLower(s), " ")
	set := ""
	for _, w := range slice {
		if !StringContains(set, w) {
			if len(set) == 0 {
				set = w
			} else {
				set = set + " " + w
			}
		}
	}
	return strings.Split(set, " ")
}

// SliceContains checks if a given string exists in a slice of strings.
//
func SliceContains(slice []string, s string) bool {
	for _, w := range slice {
		if w == s {
			return true
		}
	}
	return false
}

// StringContains checks if a given string exists in a string.
//
func StringContains(strToSplit string, s string) bool {
	slice := strings.Split(strToSplit, " ")
	for _, w := range slice {
		if w == s {
			return true
		}
	}
	return false
}

// Intersect computes the instersection of two strings.
// From two strings with format "str1 str2" and "str2 str3" in this example the result is "str2".
//
func Intersect(a string, b string) string {
	sa := SetOfStrings(a)
	sb := SetOfStrings(b)
	intersect := ""
	for _, val := range sa {
		if SliceContains(sb, val) {
			if len(intersect) == 0 {
				intersect = val
			} else {
				intersect = intersect + " " + val
			}
		}
	}
	return intersect
}

// CountTerm counts the occurence of a word into a slice of words.
//
func CountTerm(words []string, w string) int64 {
	var c int64
	for _, wi := range words {
		if wi == w {
			c = c + 1
		}
	}
	return c
}
