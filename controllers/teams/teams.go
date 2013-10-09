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
	"fmt"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	"github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/handlers"
	teamrelshlp "github.com/santiaago/purple-wing/helpers/teamrels"
	
	teammdl "github.com/santiaago/purple-wing/models/team"
	usermdl "github.com/santiaago/purple-wing/models/user"
	searchmdl "github.com/santiaago/purple-wing/models/search"
	teaminvidmdl "github.com/santiaago/purple-wing/models/teamInvertedIndex"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
)

type NewForm struct {
	Name string
	Private bool
	Error string
}

func Index(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	funcs := template.FuncMap{
		"Teams": func() bool {return true},
	}
	
	t := template.Must(template.New("tmpl_team_index").
		ParseFiles("templates/team/index.html"))
	if r.Method == "GET"{
		teams := teammdl.FindAll(r)
		indexData := struct { 
			Teams []*teammdl.Team
			TeamInputSearch string
		}{
			teams,
			"",
		}
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf,"tmpl_team_index", indexData)
		index := buf.Bytes()
		
		if err != nil{
			c.Errorf("pw: error in parse template team_index: %v", err)
		}

		err = templateshlp.Render(w, r, index, &funcs, "renderTeamIndex")
		if err != nil{
			c.Errorf("pw: error when calling Render from helpers: %v", err)
		}
	}else if r.Method == "POST"{
		query := r.FormValue("TeamInputSearch")
		words := helpers.SetOfStrings(query)
		ids := teaminvidmdl.GetIndexes(r,words)
		c.Infof("pw: search:%v Ids:%v",query, ids)
		//result := searchmdl.Score(r, query, ids)
		searchmdl.Score(r, query, ids)

		teams := teammdl.FindAll(r)
		indexData := struct { 
			Teams []*teammdl.Team
			TeamInputSearch string
		}{
			teams,
			"q:"+r.FormValue("TeamInputSearch"),
		}
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf,"tmpl_team_index", indexData)
		index := buf.Bytes()
		
		if err != nil{
			c.Errorf("pw: error in parse template team_index: %v", err)
		}

		err = templateshlp.Render(w, r, index, &funcs, "renderTeamIndex")
		if err != nil{
			c.Errorf("pw: error when calling Render from helpers: %v", err)
		}
	}
}

func New(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	funcs := template.FuncMap{}
	
	t := template.Must(template.New("tmpl_team_new").
		ParseFiles("templates/team/new.html"))
	
	var form NewForm
	if r.Method == "GET" {
		form.Name = ""
		form.Error = ""
	} else if r.Method == "POST" {
		form.Name = r.FormValue("Name")
		form.Private = (r.FormValue("Visibility") == "Private")
		
		if len(form.Name) <= 0 {
			form.Error = "'Name' field cannot be empty"
		} else if t := teammdl.Find(r, "KeyName", helpers.TrimLower(form.Name)); t != nil {
			form.Error = "That team name already exists."
		} else {
			team := teammdl.Create(r, form.Name, auth.CurrentUser(r).Id, form.Private)
			// join the team
			teamrelmdl.Create(r, team.Id, auth.CurrentUser(r).Id)
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

	err = templateshlp.Render(w, r, edit, &funcs, "renderTeamNew")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
	}
	
	funcs := template.FuncMap{
		"Joined": func() bool { return teammdl.Joined(r, intID, auth.CurrentUser(r).Id) },
		"IsTeamAdmin": func() bool { return teammdl.IsTeamAdmin(r, intID, auth.CurrentUser(r).Id) },
	}
	
	t := template.Must(template.New("tmpl_team_show").
		Funcs(funcs).
		ParseFiles("templates/team/show.html",
		"templates/team/players.html"))

	var team *teammdl.Team
	team, err = teammdl.ById(r, intID)
	
	if err != nil{
		helpers.Error404(w)
		return
	}
	
	players := teamrelshlp.Players(r, intID)
	
	teamData := struct { 
		Team *teammdl.Team
		Players []*usermdl.User 
	}{
		team,
		players,
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "tmpl_team_show", teamData)
	show := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template team_show: %v", err)
	}

	err = templateshlp.Render(w, r, show, &funcs, "renderTeamShow")
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %v", err)
	}
}

func Edit(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
	}
	
	if !teammdl.IsTeamAdmin(r, intID, auth.CurrentUser(r).Id) {
		http.Redirect(w, r, "/m", http.StatusFound)
	}
	
	if r.Method == "GET" {

		funcs := template.FuncMap{}
		
		t := template.Must(template.New("tmpl_team_edit").
			ParseFiles("templates/team/edit.html"))

		var team *teammdl.Team
		team, err = teammdl.ById(r, intID)

		if err != nil{
			helpers.Error404(w)
			return
		}

		var buf bytes.Buffer
		err = t.ExecuteTemplate(&buf,"tmpl_team_edit", team)
		edit := buf.Bytes()

		if err != nil{
			c.Errorf("pw: error in parse template team_edit: %v", err)
		}

		err = templateshlp.Render(w, r, edit, &funcs, "renderTeamEdit")
		if err != nil{
			c.Errorf("pw: error when calling Render from helpers: %v", err)
		}
	}else if r.Method == "POST"{
		
		var team *teammdl.Team
		team, err = teammdl.ById(r,intID)
		if err != nil{
			c.Errorf("pw: Team Edit handler: team not found. id: %v",intID)		
			helpers.Error404(w)
			return
		}
		// only work on name and private. Other values should not be editable
		editName := r.FormValue("Name")
		editPrivate := (r.FormValue("Visibility") == "Private")
		c.Infof("pw: Name=%s, Private=%s", editName, editPrivate)

		if helpers.IsStringValid(editName) && (editName != team.Name || editPrivate != team.Private) {
			team.Name = editName
			team.Private = editPrivate
			teammdl.Update(r, intID, team)
		}else{
			c.Errorf("pw: cannot update isStringValid: %v", helpers.IsStringValid(editName))
		}
		url := fmt.Sprintf("/m/teams/%d",intID)
		http.Redirect(w, r, url, http.StatusFound)
	}

}
