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

package models

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"appengine"
	"appengine/datastore"
)

type User struct {
	Id int64
	Email string
	Username string
	Auth string
	Created time.Time
}

type GPlusUserInfo struct {
	Id string
	Email string
	Name string
	GivenName string
	FamilyName string
}

func Create(r *http.Request, email string, username string, auth string) *User {
	c := appengine.NewContext(r)
	// create new user
	userId, _, _ := datastore.AllocateIDs(c, "User", nil, 1)
	key := datastore.NewKey(c, "User", "", userId, nil)

	user := &User{ userId, email, username, auth, time.Now() }

	_, err := datastore.Put(c, key, user)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	return user;
}

func Find(r *http.Request, email string) *User {
	c := appengine.NewContext(r)
	
	q := datastore.NewQuery("User").Filter("Email =", email)
	
	var users []*User
	
	if _, err := q.GetAll(c, &users); err == nil && len(users) > 0 {
		return users[0]
	}
	
	return nil
}

func FetchUserInfo(r *http.Request, c *http.Client) (*GPlusUserInfo, error) {
	// Make the request.
	request, err := c.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	
	if err != nil {
		return nil, err
	}

	if body, err := ioutil.ReadAll(request.Body); err == nil {
		var ui *GPlusUserInfo

		if err := json.Unmarshal(body, &ui); err == nil {
			return ui, err
		}	
	}

	return nil, err
}