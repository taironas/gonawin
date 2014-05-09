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
	"encoding/json"
	"fmt"
	"io/ioutil"
	golog "log"
	"math/rand"
	"net/http"

	"appengine"
	"appengine/urlfetch"
	"appengine/user"

	"github.com/santiaago/gonawin/helpers/log"

	gwconfig "github.com/santiaago/gonawin/config"
	mdl "github.com/santiaago/gonawin/models"
)

var (
	config       *gwconfig.GwConfig
	KOfflineMode bool
)

func init() {
	// read config file.
	var err error
	if config, err = gwconfig.ReadConfig(""); err != nil {
		golog.Printf("Error: unable to read config file; %v", err)
	}
	KOfflineMode = config.OfflineMode
}

type UserInfo struct {
	Id    string
	Email string
	Name  string
}

type TwitterUserInfo struct {
	Id          int64
	Name        string
	Screen_name string
}

// Check user validity from an accessToken string. verify if google user account is valid.
func CheckUserValidity(r *http.Request, url string, accessToken string) bool {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get(url + "=" + accessToken)
	if err != nil {
		log.Errorf(c, " CheckUserValidity: %v", err)
	}

	return resp.StatusCode == 200
}

// Check if authorization information in HTTP.Request is valid,
// ie: if it matches a user.
func CheckAuthenticationData(r *http.Request) *mdl.User {
	return mdl.FindUser(appengine.NewContext(r), "Auth", r.Header.Get("Authorization"))
}

// Is app in offline mode and email an offline user.
func isEmailOfflineUser(email string) bool {
	if !config.OfflineMode {
		return false
	}
	return config.OfflineUser.Email == email
}

// Is gonnawin admin checks if user is gonawin admin.
func IsGonawinAdmin(c appengine.Context) bool {
	return user.IsAdmin(c)
}

// unmarshal twitter response
func FetchTwitterUserInfo(r *http.Response) (*TwitterUserInfo, error) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		var ui *TwitterUserInfo

		if err = json.Unmarshal(body, &ui); err == nil {
			return ui, err
		}
	}
	return nil, err
}

// returns pointer to current user, from authentication cookie.
func CurrentOfflineUser(r *http.Request, c appengine.Context) *mdl.User {
	var u *mdl.User
	if config.OfflineMode {
		if currentUser := mdl.FindUser(c, "Username", config.OfflineUser.Name); currentUser == nil {
			u, _ = mdl.CreateUser(c, config.OfflineUser.Email, config.OfflineUser.Name, config.OfflineUser.Username, "", true, mdl.GenerateAuthKey())
		} else {
			u = currentUser
		}
	}
	return u
}

// returns user information from Google Accounts user
// if on development server only email (example@example.com) will be present.
// So Id and Name will be added.
func GetUserGoogleInfo(u *user.User) UserInfo {
	if appengine.IsDevAppServer() {
		return UserInfo{Id: fmt.Sprintf("%d", rand.Int63()), Email: u.Email, Name: "John Smith"}
	}
	return UserInfo{Id: u.ID, Email: u.Email, Name: u.String()}
}
