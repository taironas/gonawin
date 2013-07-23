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
	"bytes"
	"html/template"
	"net/http"
	"time"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/models"
)

func AdminShow(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	t, err := template.ParseFiles("templates/admin/show.html")
	
	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_admin_show", nil)
	show := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template admin_show: %v", err)
	}

	err = helpers.Render(c, w, show, nil, "renderAdminShow")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}

}

func AdminUsers(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	t, err := template.ParseFiles("templates/admin/users.html",
		"templates/user/info.html")

	// sample of users
	user1 := models.User{ 1, "test1@example.com", "John Doe 1", nil, time.Now() }
	user2 := models.User{ 1, "test2@example.com", "John Doe 2", nil, time.Now() }
	user3 := models.User{ 1, "test3@example.com", "John Doe 3", nil, time.Now() }
	users := [] models.User{user1, user2, user3}
	// end samlpe of users

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_admin_users_show", users)
	show := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template admin_users_show: %v", err)
	}

	err = helpers.Render(c, w, show, nil, "renderAdminUsersShow")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}

}
