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

package sessions

import (
	"errors"
	"net/http"
	
	"appengine"
	
	"github.com/santiaago/purple-wing/helpers"
	userhlp "github.com/santiaago/purple-wing/helpers/user"
	authhlp "github.com/santiaago/purple-wing/helpers/auth"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
)

func JsonGoogleAuth(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	userInfo := userhlp.GPlusUserInfo{r.FormValue("id"), r.FormValue("email"), r.FormValue("name")}
	
	var err error
	var user *usermdl.User
	
	if !authhlp.CheckUserValidity(r.FormValue("access_token"), r) {
		return helpers.InternalServerError{errors.New("Access token is not valid")}
	}
	if !authhlp.IsAuthorizedWithGoogle(&userInfo) {
		return helpers.Forbidden{errors.New("You are not authorized to log in to purple-wing")}
	}	
	if user, err = authhlp.SigninUser(w, r, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil{
		return helpers.InternalServerError{errors.New("Error occurred during signin process")}
	}

	// return user
	return templateshlp.RenderJson(w, c, user)
}
