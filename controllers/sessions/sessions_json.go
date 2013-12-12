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
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	
	"appengine"
	
	userhlp "github.com/santiaago/purple-wing/helpers/user"
	authhlp "github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/log"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
)

func JsonGoogleAuth(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	access_token := r.FormValue("access_token")
	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	
	userInfo := userhlp.GPlusUserInfo{id, email, name}
	
	var err error
	var user *usermdl.User
	
	log.Infof(c, "JsonGoogleAuth: access_token=%s, name=%s, email=%s", access_token, name, email)
	
	if authhlp.IsAuthorizedWithGoogle(&userInfo) {
		if user, err = authhlp.SignupUser(w, r, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil{
			log.Errorf(c, " SignupUser: %v", err)
			http.Redirect(w, r, root, http.StatusFound)
			return
		}
		
		// store in memcache auth key in memcaches
		authhlp.StoreAuthKey(c, user.Id, user.Auth)
	}
}

func renderJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	js, _ := json.Marshal(data)
	
	fmt.Fprint(w, js)
}
