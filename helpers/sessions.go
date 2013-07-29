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

package helpers

import (
    "net/http"
    "io"
    "fmt"
    "crypto/rand"
	"strconv"
	"time"
    
    "appengine"
    "appengine/memcache"
    
	usermdl "github.com/santiaago/purple-wing/models/user"
)

func IsAuthorized(ui *usermdl.GPlusUserInfo) bool {
	return ui != nil && (ui.Email == "remy.jourde@gmail.com" || ui.Email == "santiago.ariassar@gmail.com")
}

func StoreAuthKey(r *http.Request, uid int64, auth []byte) {
    c := appengine.NewContext(r)
	
	item := &memcache.Item{
        Key:   fmt.Sprintf("auth:%x", auth),
        Value: []byte(fmt.Sprintf("%d", uid)),
    }
    // Set the item, unconditionally
    if err := memcache.Set(c, item); err != nil {
        c.Errorf("pw: error setting item: %v", err)
    }
}

func fetchAuthKey(r *http.Request, auth string) string {
    c := appengine.NewContext(r)

    // Get the item from the memcache
    if item, err := memcache.Get(c, "auth:"+auth); err == nil {
        return fmt.Sprintf("%s", item.Value)
    } 
    
    return ""
}

func SetAuthCookie(w http.ResponseWriter, auth []byte) {
	cookie := &http.Cookie{ 
        Name: "auth", 
        Value: fmt.Sprintf("%x", auth), 
        Path: "/m",
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
		Path:    "/m",
		Expires: time.Now(),
	})
}

func GenerateAuthKey() []byte {
    b := make([]byte, 16)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return nil
    }
    return b
}

func CurrentUser(r *http.Request) *usermdl.User {
	if auth := GetAuthCookie(r); len(auth) > 0 {
		if uid := fetchAuthKey(r, auth); len(uid) > 0 {
			userId, _ := strconv.ParseInt(uid, 10, 64)
			return usermdl.Find(r, "Id", userId)
		}
	}
	
	return nil
}


func LoggedIn(r *http.Request) bool {
	if auth := GetAuthCookie(r); len(auth) > 0 {
		if u := CurrentUser(r); u != nil {
			return fmt.Sprintf("%x", u.Auth) == auth
		}
	}
	
	return false
}
