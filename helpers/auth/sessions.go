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
	"strconv"
	"time"
	
	"appengine"
	"appengine/datastore"
	"appengine/memcache"

	"github.com/santiaago/purple-wing/helpers/log"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

type AuthKey struct {
	Id int64
	Key string
	Value string
	Created time.Time
}

// Store authentication key in datastore and update memcache
func StoreAuthKey(c appengine.Context, uid int64, auth string) {
	
	uidStr := fmt.Sprintf("%d", uid)
	
	// create new team
	authKeyId, _, err := datastore.AllocateIDs(c, "AuthKey", nil, 1)
	if err != nil {
		log.Errorf(c, " StoreAuthKey: %v", err)
	}
	
	key := datastore.NewKey(c, "AuthKey", "", authKeyId, nil)

	authKey := &AuthKey{ authKeyId, auth, uidStr, time.Now() }

	_, err = datastore.Put(c, key, authKey)
	if err != nil {
		log.Errorf(c, " StoreAuthKey: %v", err)
	}
	
	item := &memcache.Item{
		Key:   "auth:" + auth,
		Value: []byte(uidStr),
	}
	// Set the item, unconditionally
	if err := memcache.Set(c, item); err == memcache.ErrNotStored {
		log.Infof(c, "item with key %q already exists", item.Key)
	} else if err != nil {
		log.Errorf(c, "error adding item: %v", err)
	}
}

// fetch Authentication key with respect to auth string from memcache or datastore
func fetchAuthKey(r *http.Request, auth string) string {
	c := appengine.NewContext(r)

	// Get the item from the memcache
	if item, err := memcache.Get(c, "auth:"+auth); err == memcache.ErrCacheMiss {
		// item doesn't exist in memcache, retrieve it in datastore
		q := datastore.NewQuery("AuthKey").Filter("Key =", auth).Limit(1)

		var authKeys []*AuthKey

		if _, err := q.GetAll(c, &authKeys); err == nil && len(authKeys) > 0 {
			return authKeys[0].Value
		}
	} else if err != nil {
		log.Errorf(c, " error getting item: %v", err)
	} else {
		return fmt.Sprintf("%s", item.Value)
	}

	return ""
}

// Set cookie with authentication string
func SetAuthCookie(w http.ResponseWriter, auth string) {
	cookie := &http.Cookie{ 
		Name: "auth", 
		Value: auth, 
		Path: "/",
	}
	http.SetCookie(w, cookie)
}

// extract authentication value from http.Request cookie
func GetAuthCookie(r *http.Request) string {
	if cookie, err := r.Cookie("auth"); err == nil {
		return cookie.Value
	}
	return ""
}

// clear authentication cookie
func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "auth",
		Path:    "/",
		Expires: time.Now(),
	})
}

// returns pointer to current user, from authentication cookie.
func CurrentUser(r *http.Request, c appengine.Context) *usermdl.User {
	if KOfflineMode {
		currentUser := usermdl.Find(c, "Username", "purple")
    
    if currentUser == nil {
      currentUser, _ = usermdl.Create(c, "purple@wing.com", "purple", "wing", usermdl.GenerateAuthKey())
    }
    return currentUser
	}
	
	if auth := GetAuthCookie(r); len(auth) > 0 {
		if uid := fetchAuthKey(r, auth); len(uid) > 0 {
			log.Infof(c, " CurrentUser, uid=%s, auth=%s", uid, auth)
			userId, err := strconv.ParseInt(uid, 10, 64)
			if err != nil {
				log.Errorf(c, " CurrentUser, string value could not be parsed: %v", err)
			}
			
			return usermdl.Find(c, "Id", userId)
		}
	}

	return nil
}
