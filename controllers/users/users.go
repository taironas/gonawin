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

package users

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

type Form struct {
	Username string
	Name string
	Email string
	ErrorUsername string
	ErrorName string
	ErrorEmail string
}

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	funcs := template.FuncMap{
		"Profile": func() bool {return true},
	}
	
	t := template.Must(template.New("tmpl_user_show").
		Funcs(funcs).
		ParseFiles("templates/user/show.html", 
		"templates/user/info.html"))

	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/users/", http.StatusFound)
	}
	var user *usermdl.User
	user, err = usermdl.ById(r,intID)
	if err != nil{
		helpers.Error404(w)
		return
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_user_show", *user)
	show := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template user_show: %v", err)
	}

	err = templateshlp.Render(w, r, show, &funcs, "renderUserShow")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request){
	
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		funcs := template.FuncMap{
			"Profile": func() bool {return true},
		}
		
		t := template.Must(template.New("tmpl_user_show").
			Funcs(funcs).
			ParseFiles("templates/user/show.html", 
			"templates/user/edit.html"))
		
		intID, err := handlers.PermalinkID(r,3)
		if err != nil{
			http.Redirect(w,r, "/m/users/", http.StatusFound)
		}
		var user *usermdl.User
		user, err = usermdl.ById(r,intID)
		if err != nil{
			helpers.Error404(w)
			return
		}
		
		var buf bytes.Buffer
		err = t.ExecuteTemplate(&buf,"tmpl_user_edit", *user)
		edit := buf.Bytes()
		
		if err != nil{
			c.Errorf("pw: error in parse template user_edit: %v", err)
		}
		
		err = templateshlp.Render(w, r, edit, &funcs, "renderUserEdit")
		if err != nil{
			c.Errorf("pw: error when calling Render from helpers: %v", err)
		}
	}else if r.Method == "POST"{
		intID, err := handlers.PermalinkID(r,3)
		if err != nil{
			http.Redirect(w,r, "/m/users/", http.StatusFound)
		}
		var user *usermdl.User
		user, err = usermdl.ById(r,intID)
		if err != nil{
			c.Errorf("pw: User Edit handler: user not found. id: %v",intID)		
			helpers.Error404(w)
			return
		}
		// only work on username other values should not be editable
		editUserName := r.FormValue("Username")
		
		if helpers.IsUsernameValid(editUserName) && editUserName != user.Username{
			user.Username = editUserName
			usermdl.Update(r, intID, user)
		}else{
			c.Errorf("pw: cannot update %v", helpers.IsUsernameValid(editUserName))
		}
		url := fmt.Sprintf("/m/users/%d",intID)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func renderJsonUser(w http.ResponseWriter, user *usermdl.User) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := helpers.JsonResponse {
		"id": fmt.Sprintf("%d", user.Id),
		"email": user.Email,
		"username": user.Username,
		"name": user.Name,
		"auth": user.Auth,
		"created": user.Created,
	} 

	fmt.Fprint(w, data.String())
}
