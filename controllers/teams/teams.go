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

package teams

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
	"fmt"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/auth"

	teammdl "github.com/santiaago/purple-wing/models/team"
)

type Form struct {
	Name string
	Error string
}

func Index(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	funcs := template.FuncMap{
		"Teams": func() bool {return true},
	}
	
	t := template.Must(template.New("tmpl_team_index").
		ParseFiles("templates/team/index.html"))
	
	teams := teammdl.FindAll(r)
	
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_team_index", teams)
	index := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template team_index: %v", err)
	}

	err = helpers.Render(w, r, index, &funcs, "renderTeamIndex")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}

func New(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	funcs := template.FuncMap{}
	
	t := template.Must(template.New("tmpl_team_new").
		ParseFiles("templates/team/new.html"))
	
	var form Form
	if r.Method == "GET" {
		form.Name = ""
		form.Error = ""
	} else if r.Method == "POST" {
		form.Name = r.FormValue("Name")
		
		if len(form.Name) <= 0 {
			form.Error = "'Name' field cannot be empty"
		} else if t := teammdl.Find(r, "KeyName", helpers.TrimLower(form.Name)); t != nil {
			form.Error = "That team name already exists."
		} else {
			team := teammdl.Create(r, form.Name, auth.CurrentUser(r).Id)
			// redirect to the newly created team page
			http.Redirect(w, r, "/m/teams/" + fmt.Sprintf("%d", team.Id), http.StatusFound)
		}
	}
	
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_team_new", form)
	edit := buf.Bytes()

	if err != nil{
		c.Errorf("pw: error in parse template team_new: %v", err)
	}

	err = helpers.Render(w, r, edit, &funcs, "renderTeamNew")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	funcs := template.FuncMap{}
	
	t := template.Must(template.New("tmpl_team_show").
		ParseFiles("templates/team/show.html"))
	
	team := teammdl.Team{ 1, "team foo", "Team Foo", 1, time.Now() }
	
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_team_show", team)
	show := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template team_show: %v", err)
	}

	err = helpers.Render(w, r, show, &funcs, "renderTeamShow")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	funcs := template.FuncMap{}
	
	t := template.Must(template.New("tmpl_team_show").
		ParseFiles("templates/team/show.html", 
		"templates/team/edit.html"))

	team := teammdl.Team{ 1, "team foo", "Team Foo", 1, time.Now() }

	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_team_edit", team)
	edit := buf.Bytes()

	if err != nil{
		c.Errorf("pw: error in parse template team_edit: %v", err)
	}

	err = helpers.Render(w, r, edit, &funcs, "renderTeamEdit")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}
