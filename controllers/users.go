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

func UserShow(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	t, err := template.ParseFiles("templates/user/show.html")
	
	user := models.User{ 1, "test@example.com", "John Doe", time.Now() }

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_user_show", user)
	show := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template user_show: %q", err)
	}

	renderUser(c, w, helpers.Content{template.HTML(show)})
}

func UserEdit(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	t, err := template.ParseFiles("templates/user/edit.html")

	user := models.User{ 1, "test@example.com", "John Doe", time.Now() }

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_user_edit", user)
	show := buf.Bytes()

	if err != nil{
		c.Errorf("pw: error in parse template user_edit: %q", err)
	}

	renderUser(c, w, helpers.Content{template.HTML(show)})
}

// renderShowUser executes the user show/edit template.
func renderUser(c appengine.Context, w http.ResponseWriter, content helpers.Content) {
	tmpl, err := template.ParseFiles("templates/layout/application.html", 
									 "templates/layout/container.html",
									 "templates/layout/header.html",
									 "templates/layout/footer.html",
									 "templates/layout/scripts.html" )
	if err != nil{
		c.Errorf("error in parse files: %q", err)
	}

	err = tmpl.ExecuteTemplate(w,"tmpl_application",content)
	if err != nil{
		c.Errorf("error in execute template: %q", err)
	}
}
