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
 
package settings

import (
	"net/http"
	"html/template"
	
	"appengine"
	
	"github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
)

func Profile(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	if r.Method == "GET" {
		funcs := template.FuncMap{
			"Profile": func() bool {return true},
		}
		
		t := template.Must(template.New("tmpl_settings_profile").
			Funcs(funcs).
			ParseFiles("templates/settings/profile.html"))
		
		templateshlp.RenderWithData(w, r, t, auth.CurrentUser(r), funcs, "renderProfile")

	}else if r.Method == "POST"{
		currentUser := auth.CurrentUser(r)
		
		editUserName := r.FormValue("Username")
		
		if helpers.IsUsernameValid(editUserName) && editUserName != currentUser.Username{
			currentUser.Username = editUserName
			usermdl.Update(r, currentUser)
		} else {
			log.Errorf(c, " cannot update current user info")
		}
		
		http.Redirect(w, r, "/m/settings/edit-profile", http.StatusFound)
	}
}

func Networks(w http.ResponseWriter, r *http.Request){

	t := template.Must(template.New("tmpl_settings_networks").
		ParseFiles("templates/settings/networks.html"))

	// no data
	funcs := template.FuncMap{}
	templateshlp.RenderWithData(w, r, t, nil, funcs, "renderNetworks")
}

func Email(w http.ResponseWriter, r *http.Request){

	t := template.Must(template.New("tmpl_settings_email").
		ParseFiles("templates/settings/email.html"))

	// no data
	funcs := template.FuncMap{}
	templateshlp.RenderWithData(w, r, t, nil, funcs, "renderEmail")
}
