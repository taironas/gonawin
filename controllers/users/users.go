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
	"html/template"
	"net/http"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	teamrelshlp "github.com/santiaago/purple-wing/helpers/teamrels"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
	teammdl "github.com/santiaago/purple-wing/models/team"
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
	
	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/users/", http.StatusFound)
	}
	
	funcs := template.FuncMap{
		"Profile": func() bool {return true},
	}
	
	t := template.Must(template.New("tmpl_user_show").
		Funcs(funcs).
		ParseFiles("templates/user/show.html", 
		"templates/user/info.html",
		"templates/user/teams.html"))

	var user *usermdl.User
	user, err = usermdl.ById(r,intID)
	if err != nil{
		helpers.Error404(w)
		return
	}
	
	teams := teamrelshlp.Teams(r, intID)
	
	userData := struct { 
		User *usermdl.User
		Teams []*teammdl.Team 
	}{
		user,
		teams,
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "tmpl_user_show", userData)
	show := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template user_show: %v", err)
	}

	err = templateshlp.Render(w, r, show, &funcs, "renderUserShow")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}
