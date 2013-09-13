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

package user

import (
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	
	"appengine"
	"appengine/datastore"
	
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
)

type User struct {
	Id int64
	Email string
	Username string
	Name string
	Auth string
	Created time.Time
}

type GPlusUserInfo struct {
	Id string
	Email string
	Name string
}

type TwitterUserInfo struct {
	Id int64
	Name string
	Screen_name string
}

func Create(r *http.Request, email string, username string, name string, auth string) *User {
	c := appengine.NewContext(r)
	// create new user
	userId, _, _ := datastore.AllocateIDs(c, "User", nil, 1)
	key := datastore.NewKey(c, "User", "", userId, nil)
	
	user := &User{ userId, email, username, name, auth, time.Now() }

	_, err := datastore.Put(c, key, user)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	return user;
}

func Find(r *http.Request, filter string, value interface{}) *User{
	q := datastore.NewQuery("User").Filter(filter + " =", value)
	
	var users []*User
	
	if _, err := q.GetAll(appengine.NewContext(r), &users); err == nil && len(users) > 0 {
		return users[0]
	}
	
	return nil
}

func ById(r *http.Request, id int64)(*User, error){
	c := appengine.NewContext(r)

	var u User
	key := datastore.NewKey(c, "User", "", id, nil)

	if err := datastore.Get(c, key, &u); err != nil {
		c.Errorf("pw: user not found : %v", err)
		return &u, err
	}
	return &u, nil
}

func KeyById(r *http.Request, id int64)(*datastore.Key){
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "User", "", id, nil)

	return key
}

func Update(r *http.Request, id int64, u *User) error{
	c := appengine.NewContext(r)
	k := KeyById(r, id)
	if _, err := datastore.Put(c, k, u); err != nil {
		return err
	}
	return nil
}

func FetchGPlusUserInfo(r *http.Request, c *http.Client) (*GPlusUserInfo, error) {
	// Make the request.
	request, err := c.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	
	if err != nil {
		return nil, err
	}

	defer request.Body.Close()
	
	if body, err := ioutil.ReadAll(request.Body); err == nil {
		var ui *GPlusUserInfo

		if err := json.Unmarshal(body, &ui); err == nil {
			return ui, err
		}	
	}

	return nil, err
}

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

func Joined(r *http.Request, teamId int64, userId int64) bool {
	_, err := teamrelmdl.Find(r, "TeamId", teamId)
	return err == nil
}

func Join(r *http.Request, teamId int64, userId int64) error {
	if teamRel := teamrelmdl.Create(r, teamId, userId); teamRel == nil {
		return errors.New("error during team relationship creation")
	}
	
	return nil
}

func Leave(r *http.Request, teamId int64, userId int64) error {
	return teamrelmdl.Destroy(r, teamId, userId)
}