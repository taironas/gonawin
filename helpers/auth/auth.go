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

	"appengine"
	"appengine/urlfetch"

	"github.com/santiaago/purple-wing/helpers/log"

	usermdl "github.com/santiaago/purple-wing/models/user"
)

const KOfflineMode bool = false

const kEmailRjourde = "remy.jourde@gmail.com"
const kEmailSarias = "santiago.ariassar@gmail.com"
const kEmailGonawinTest = "gonawin.test@gmail.com"

type GPlusUserInfo struct {
	Id    string
	Email string
	Name  string
}

// from an accessToken string, verify if google user account is valid
func CheckGoogleUserValidity(accessToken string, r *http.Request) bool {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get("https://www.google.com/accounts/AuthSubTokenInfo?bearer_token=" + accessToken)
	if err != nil {
		log.Errorf(c, " CheckUserValidity: %v", err)
	}

	return resp.StatusCode == 200
}

// Check if authorization information in HTTP.Request is valid,
// ie: if it matches a user.
func CheckAuthenticationData(r *http.Request) *usermdl.User {
	return usermdl.Find(appengine.NewContext(r), "Auth", r.Header.Get("Authorization"))
}

// Ckeck if googple plus user is admin.
// #196: Should be removed when deployed in production.
func IsAuthorizedWithGoogle(ui *GPlusUserInfo) bool {
	return ui != nil && (ui.Email == kEmailRjourde || ui.Email == kEmailSarias || ui.Email == kEmailGonawinTest)
}

// returns pointer to current user, from authentication cookie.
func CurrentOfflineUser(r *http.Request, c appengine.Context) *usermdl.User {
	if KOfflineMode {
		currentUser := usermdl.Find(c, "Username", "purple")

		if currentUser == nil {
			currentUser, _ = usermdl.Create(c, "purple@wing.com", "purple", "wing", true, usermdl.GenerateAuthKey())
		}
		return currentUser
	} else {
		return nil
	}
}
