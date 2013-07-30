/*
 * Copyright (c) 2013 Santiago Arias | Remy Jourde | Carlos Bernal
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

package auth

import (
	"net/http"
	"fmt"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
)

func IsAuthorized(ui *usermdl.GPlusUserInfo) bool {
	return ui != nil && (ui.Email == "remy.jourde@gmail.com" || ui.Email == "santiago.ariassar@gmail.com")
}

// LoggedIn is true is the AuthCookie exist and match your user.Auth property
func LoggedIn(r *http.Request) bool {
	if auth := GetAuthCookie(r); len(auth) > 0 {
		if u := CurrentUser(r); u != nil {
			return fmt.Sprintf("%x", u.Auth) == auth
		}
	}
	
	return false
}

// IsAdmin is true if you are logged in and belong to the below users.
func IsAdmin(r *http.Request) bool {

	if LoggedIn(r){
		if u := CurrentUser(r); u != nil{
			return (u.Email == "remy.jourde@gmail.com" || u.Email == "santiago.ariassar@gmail.com")
		}
	}
	return false
}

// IsUser is true if you are logged in, can either be an admin or not.
func IsUser(r *http.Request) bool {
	return LoggedIn(r)
}
