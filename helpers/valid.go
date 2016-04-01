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
	"regexp"
	"strings"
)

// UsernameRegExp is the regexp for username.
//
var UsernameRegExp = regexp.MustCompile(`^[\sa-zA-Z0-9_-]{3,20}$`)

// IsUsernameValid returns true is username is a string between 3 and 20 characters.
// note: added \s just temporary until user name is splited with name and last name
// username should not have whitespaces
//
func IsUsernameValid(username string) bool {

	return UsernameRegExp.MatchString(username)
}

// PasswordRegexp is the regexp for password.
//
var PasswordRegexp = regexp.MustCompile(`^.{3,20}$`)

// IsPasswordValid returns true for all passwords that are between 3 and 20 characters.
//
func IsPasswordValid(password string) bool {

	return PasswordRegexp.MatchString(password)
}

// emailRegexp is the regexp for email.
//
var emailRegexp = regexp.MustCompile(`^[\S]+@[\S]+\.[\S]+$`)

// IsEmailValid returns true if string email has the form a@b.c
//
func IsEmailValid(email string) bool {

	return emailRegexp.MatchString(email)
}

// AreEmailsValid checks if an array of emails is valid.
//
func AreEmailsValid(emails []string) bool {
	for _, email := range emails {
		if !IsEmailValid(strings.TrimSpace(email)) {
			return false
		}
	}

	return true
}

// IsStringValid returns true if string is not empty.
func IsStringValid(s string) bool {
	return len(s) > 0
}
