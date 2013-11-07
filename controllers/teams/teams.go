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
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	"github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/handlers"
	teamrelshlp "github.com/santiaago/purple-wing/helpers/teamrels"
	tournamentrelshlp "github.com/santiaago/purple-wing/helpers/tournamentrels"
	
	teammdl "github.com/santiaago/purple-wing/models/team"
	usermdl "github.com/santiaago/purple-wing/models/user"
	searchmdl "github.com/santiaago/purple-wing/models/search"
	teaminvidmdl "github.com/santiaago/purple-wing/models/teamInvertedIndex"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
	teamrequestmdl "github.com/santiaago/purple-wing/models/teamrequest"
)

type NewForm struct {
	Name string
	Private bool
	Error string
}

type indexData struct{
	Teams []*teammdl.Team
	TeamInputSearch string	
}

func Index(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	var data indexData
	if r.Method == "GET"{
		teams := teammdl.FindAll(r)
		data.Teams = teams
		data.TeamInputSearch = ""

	} else if r.Method == "POST" {
		if query := r.FormValue("TeamInputSearch"); len(query) == 0{
			http.Redirect(w, r, "teams", http.StatusFound)
			return
		} else{
			words := helpers.SetOfStrings(query)
			ids := teaminvidmdl.GetIndexes(r,words)
			c.Infof("pw: search:%v Ids:%v",query, ids)
			result := searchmdl.TeamScore(r, query, ids)
			
			teams := teammdl.ByIds(r, result)
			data.Teams = teams
			data.TeamInputSearch = query
		}
	} else{
		helpers.Error404(w)
	}
	
	t := template.Must(template.New("tmpl_team_index").
		ParseFiles("templates/team/index.html"))

	funcs := template.FuncMap{
		"Teams": func() bool {return true},
	}

	templateshlp.Render_with_data(w, r, t, data, funcs, "renderTeamIndex")
}

func New(w http.ResponseWriter, r *http.Request){
	
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
	} else{
		helpers.Error404(w)
	}

	t := template.Must(template.New("tmpl_team_new").
		ParseFiles("templates/team/new.html"))
	
	funcs := template.FuncMap{}
	
	templateshlp.Render_with_data(w, r, t, form, funcs, "renderTeamNew")
}

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
	} 

	if r.Method == "POST" && r.FormValue("Action") == "delete" {
		// delete all team-user relationships
		for _, player := range teamrelshlp.Players(r, intID) {
			if err := teamrelmdl.Destroy(r, intID, player.Id); err !=nil {
				c.Errorf("pw: error when trying to destroy team relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, tournament := range tournamentrelshlp.Teams(r, intID) {
			if err := tournamentteamrelmdl.Destroy(r, tournament.Id, intID); err !=nil {
				c.Errorf("pw: error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the team
		teammdl.Destroy(r, intID)
		
		http.Redirect(w, r, "/m/teams", http.StatusFound)
	}
	funcs := template.FuncMap{
		"Joined": func() bool { return teammdl.Joined(r, intID, auth.CurrentUser(r).Id) },
		"IsTeamAdmin": func() bool { return teammdl.IsTeamAdmin(r, intID, auth.CurrentUser(r).Id) },
		"RequestSent": func() bool { return teamrequestmdl.Sent(r, intID, auth.CurrentUser(r).Id) },
	}
	
	t := template.Must(template.New("tmpl_team_show").
		Funcs(funcs).
		ParseFiles("templates/team/show.html",
		"templates/team/players.html"))

	var team *teammdl.Team
	if team, err = teammdl.ById(r, intID); err != nil{
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
	templateshlp.Render_with_data(w, r, t, teamData, funcs, "renderTeamShow")
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
		templateshlp.Render_with_data(w, r, t, team, funcs, "renderTeamEdit")

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
	} else {
		helpers.Error404(w)
	}
}

func Invite(w http.ResponseWriter, r *http.Request){
	
	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
	}
	
	if r.Method == "POST"{
		if teamRequest := teamrequestmdl.Create(r, intID, auth.CurrentUser(r).Id); teamRequest == nil {
			appengine.NewContext(r).Errorf("pw: no team request has been created")
		}
		
		url := fmt.Sprintf("/m/teams/%d", intID)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func Request(w http.ResponseWriter, r *http.Request){
	
	if r.Method == "POST"{
		
		requestId, _ := strconv.ParseInt(r.FormValue("RequestId"), 10,64)
		
		if r.FormValue("SubmitButton") == "Accept" {
			if teamRequest, err := teamrequestmdl.ById(r, requestId); err == nil {
				// join user to the team
				teammdl.Join(r, teamRequest.TeamId, teamRequest.UserId);
			} else {
				appengine.NewContext(r).Errorf("pw: cannot find team request with id=%d", requestId)
			}
		}
		
		teamrequestmdl.Destroy(r, requestId)
		
		url := fmt.Sprintf("/m/users/%d", auth.CurrentUser(r).Id)
		http.Redirect(w, r, url, http.StatusFound)
	}
}
