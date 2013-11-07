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
	"io"
	"errors"
	"fmt"
	"crypto/rand"
	"strconv"
	"time"
	
	"appengine"
	"appengine/datastore"
	"appengine/memcache"

	usermdl "github.com/santiaago/purple-wing/models/user"
)

type AuthKey struct {
	Id int64
	Key string
	Value string
	Created time.Time
}

func StoreAuthKey(r *http.Request, uid int64, auth string) {
	c := appengine.NewContext(r)
	
	uidStr := fmt.Sprintf("%d", uid)
	
	// create new team
	authKeyId, _, err := datastore.AllocateIDs(c, "AuthKey", nil, 1)
	if err != nil {
		c.Errorf("pw: StoreAuthKey: %v", err)
	}
	
	key := datastore.NewKey(c, "AuthKey", "", authKeyId, nil)

	authKey := &AuthKey{ authKeyId, auth, uidStr, time.Now() }

	_, err = datastore.Put(c, key, authKey)
	if err != nil {
		c.Errorf("pw: StoreAuthKey: %v", err)
	}
	
	item := &memcache.Item{
		Key:   "auth:" + auth,
		Value: []byte(uidStr),
	}
    // Set the item, unconditionally
	if err := memcache.Set(c, item); err == memcache.ErrNotStored {
		c.Infof("item with key %q already exists", item.Key)
	} else if err != nil {
		c.Errorf("error adding item: %v", err)
	}
}

func fetchAuthKey(r *http.Request, auth string) string {
	c := appengine.NewContext(r)

	// Get the item from the memcache
	if item, err := memcache.Get(c, "auth:"+auth); err == memcache.ErrCacheMiss {
		// item doesn't exist in memcache, retrieve it in datastore
		q := datastore.NewQuery("AuthKey").Filter("Key =", auth).Limit(1)

		var authKeys []*AuthKey

		if _, err := q.GetAll(appengine.NewContext(r), &authKeys); err == nil && len(authKeys) > 0 {
			return authKeys[0].Value
		}
	} else if err != nil {
		c.Errorf("pw: error getting item: %v", err)
	} else {
		return fmt.Sprintf("%s", item.Value)
	}

	return ""
}

func SetAuthCookie(w http.ResponseWriter, auth string) {
	cookie := &http.Cookie{ 
		Name: "auth", 
		Value: auth, 
		Path: "/",
	}
	http.SetCookie(w, cookie)
}

func GetAuthCookie(r *http.Request) string {
	if cookie, err := r.Cookie("auth"); err == nil {
		return cookie.Value
	}
	return ""
}

func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "auth",
		Path:    "/",
		Expires: time.Now(),
	})
}

func GenerateAuthKey() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return ""
	}
	return fmt.Sprintf("%x", b)
}

func CurrentUser(r *http.Request) *usermdl.User {
	c := appengine.NewContext(r)
	if auth := GetAuthCookie(r); len(auth) > 0 {
		if uid := fetchAuthKey(r, auth); len(uid) > 0 {
			c.Infof("pw: CurrentUser, uid=%s, auth=%s", uid, auth)
			userId, err := strconv.ParseInt(uid, 10, 64)
			if err != nil {
				c.Errorf("pw: CurrentUser, string value could not be parsed: %v", err)
			}
			
			return usermdl.Find(r, "Id", userId)
		}
	}

	return nil
}

func SignupUser(w http.ResponseWriter, r *http.Request, queryName string, email string, screenName string, name string) error{

	c := appengine.NewContext(r)
	var user *usermdl.User
	// find user
	if user = usermdl.Find(r, "Username", queryName); user == nil {
		// create user if it does not exist
		if userCreate, err := usermdl.Create(r, email, screenName, name, GenerateAuthKey()); err != nil{
			c.Errorf("Signup: %v", err)
			return errors.New("helpers/auth: Unable to create user.")
		}else{
			user = userCreate
		}
	}
	// set 'auth' cookie
	SetAuthCookie(w, user.Auth)
	// store in memcache auth key in memcaches
	StoreAuthKey(r, user.Id, user.Auth)
	return nil
}
