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

package controllers

import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"fmt"
	"code.google.com/p/goauth2/oauth"
	
	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/models"
)

var CurrentUser *models.User = nil

// Set up a configuration.
func config(host string) *oauth.Config{
	return &oauth.Config{
		ClientId:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Scope:        "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		RedirectURL:  fmt.Sprintf("http://%s/oauth2callback", host),
	}
}

func SessionAuth(w http.ResponseWriter, r *http.Request){
	url := config(r.Host).AuthCodeURL(r.URL.RawQuery)
	http.Redirect(w, r, url, http.StatusFound)
}

func SessionAuthCallback(w http.ResponseWriter, r *http.Request){
	// Exchange code for an access token at OAuth provider.
	code := r.FormValue("code")
	t := &oauth.Transport{
		Config: config(r.Host),
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(r),
		},
	}
	
	var userInfo *models.GPlusUserInfo
	
	if _, err := t.Exchange(code); err == nil {
		userInfo, _ = models.FetchUserInfo(r, t.Client())
	}
	if helpers.IsAuthorized(userInfo) {
		var user *models.User
		// find user
		if user = models.Find(r, userInfo.Email); user == nil {
			// create user if it does not exist
			user = models.Create(r, userInfo.Email, userInfo.Name)
			
			CurrentUser = user
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func SessionLogout(w http.ResponseWriter, r *http.Request){
	CurrentUser = nil
	
	http.Redirect(w, r, "/", http.StatusFound)
}

func LoggedIn() bool {
	return CurrentUser != nil
}
