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
	"github.com/santiaago/purple-wing/helpers/log"

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

func Index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	var data indexData
	if r.Method == "GET" {
		teams := teammdl.FindAll(c)
		data.Teams = teams
		data.TeamInputSearch = ""

	} else if r.Method == "POST" {
		if query := r.FormValue("TeamInputSearch"); len(query) == 0 {
			http.Redirect(w, r, "teams", http.StatusFound)
			return
		} else {
			words := helpers.SetOfStrings(query)
			ids, err := teaminvidmdl.GetIndexes(c, words)
			
			if err != nil {
				log.Errorf(c, " teams.Index, error occurred when getting indexes of words: %v", err)
			}
			
			result := searchmdl.TeamScore(c, query, ids)
			
			teams := teammdl.ByIds(c, result)
			data.Teams = teams
			data.TeamInputSearch = query
		}
	} else {
		helpers.Error404(w)
	}
	
	t := template.Must(template.New("tmpl_team_index").
		ParseFiles("templates/team/index.html"))

	funcs := template.FuncMap{
		"Teams": func() bool {return true},
	}

	templateshlp.RenderWithData(w, r, c, t, data, funcs, "renderTeamIndex")
}

func New(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	var form NewForm
	if r.Method == "GET" {
		form.Name = ""
		form.Error = ""
	} else if r.Method == "POST" {
		form.Name = r.FormValue("Name")
		form.Private = (r.FormValue("Visibility") == "Private")
		
		if len(form.Name) <= 0 {
			form.Error = "'Name' field cannot be empty"
		} else if t := teammdl.Find(c, "KeyName", helpers.TrimLower(form.Name)); t != nil {
			form.Error = "That team name already exists."
		} else {
			team, err := teammdl.Create(c, form.Name, auth.CurrentUser(r, c).Id, form.Private)
			if err != nil {
				log.Errorf(c, " error when trying to create a team: %v", err)
			}
			// join the team
			_, err = teamrelmdl.Create(c, team.Id, auth.CurrentUser(r, c).Id)
			if err != nil {
				log.Errorf(c, " error when trying to create a team relationship: %v", err)
			}
			// redirect to the newly created team page
			http.Redirect(w, r, "/m/teams/" + fmt.Sprintf("%d", team.Id), http.StatusFound)
		}
	} else{
		helpers.Error404(w)
	}

	t := template.Must(template.New("tmpl_team_new").
		ParseFiles("templates/team/new.html"))
	
	funcs := template.FuncMap{}
	
	templateshlp.RenderWithData(w, r, c, t, form, funcs, "renderTeamNew")
}

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
	} 

	if r.Method == "POST" && r.FormValue("Action") == "delete" {
		// delete all team-user relationships
		for _, player := range teamrelshlp.Players(c, intID) {
			if err := teamrelmdl.Destroy(c, intID, player.Id); err !=nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, tournament := range tournamentrelshlp.Teams(c, intID) {
			if err := tournamentteamrelmdl.Destroy(c, tournament.Id, intID); err !=nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the team
		teammdl.Destroy(c, intID)
		
		http.Redirect(w, r, "/m/teams", http.StatusFound)
		return
	} else if (r.Method != "GET"){
		log.Errorf(c, " request method not supported")
		helpers.Error404(w)
		return	
	}
	funcs := template.FuncMap{
		"Joined": func() bool { return teammdl.Joined(c, intID, auth.CurrentUser(r, c).Id) },
		"IsTeamAdmin": func() bool { return teammdl.IsTeamAdmin(c, intID, auth.CurrentUser(r, c).Id) },
		"RequestSent": func() bool { return teamrequestmdl.Sent(c, intID, auth.CurrentUser(r, c).Id) },
	}
	
	t := template.Must(template.New("tmpl_team_show").
		Funcs(funcs).
		ParseFiles("templates/team/show.html",
		"templates/team/players.html"))

	var team *teammdl.Team
	if team, err = teammdl.ById(c, intID); err != nil{
		helpers.Error404(w)
		return
	}
	
	players := teamrelshlp.Players(c, intID)
	
	teamData := struct { 
		Team *teammdl.Team
		Players []*usermdl.User 
	}{
		team,
		players,
	}
	templateshlp.RenderWithData(w, r, c, t, teamData, funcs, "renderTeamShow")
}

func Edit(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
	}
	
	if !teammdl.IsTeamAdmin(c, intID, auth.CurrentUser(r, c).Id) {
		http.Redirect(w, r, "/m", http.StatusFound)
	}
	
	if r.Method == "GET" {

		funcs := template.FuncMap{}
		
		t := template.Must(template.New("tmpl_team_edit").
			ParseFiles("templates/team/edit.html"))

		var team *teammdl.Team
		team, err = teammdl.ById(c, intID)

		if err != nil{
			helpers.Error404(w)
			return
		}
		templateshlp.RenderWithData(w, r, c, t, team, funcs, "renderTeamEdit")

	} else if r.Method == "POST"{
		
		var team *teammdl.Team
		team, err = teammdl.ById(c,intID)
		if err != nil{
			log.Errorf(c, " Team Edit handler: team not found. id: %v",intID)		
			helpers.Error404(w)
			return
		}
		// only work on name and private. Other values should not be editable
		editName := r.FormValue("Name")
		editPrivate := (r.FormValue("Visibility") == "Private")
		log.Infof(c, " Name=%s, Private=%s", editName, editPrivate)

		if helpers.IsStringValid(editName) && (editName != team.Name || editPrivate != team.Private) {
			team.Name = editName
			team.Private = editPrivate
			teammdl.Update(c, intID, team)
		}else{
			log.Errorf(c, " cannot update isStringValid: %v", helpers.IsStringValid(editName))
		}
		url := fmt.Sprintf("/m/teams/%d",intID)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		helpers.Error404(w)
	}
}

func Invite(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
		return
	}
	
	if r.Method == "POST"{
		if _, err := teamrequestmdl.Create(c, intID, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " teams.Invite, error when trying to create a team request: %v", err)
		}
		
		url := fmt.Sprintf("/m/teams/%d", intID)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func Request(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	if r.Method == "POST"{
		
		requestId, err := strconv.ParseInt(r.FormValue("RequestId"), 10,64)
		if err != nil {
			log.Errorf(c, " teams.Request, string value could not be parsed: %v", err)
		}
		
		if r.FormValue("SubmitButton") == "Accept" {
			if teamRequest, err := teamrequestmdl.ById(c, requestId); err == nil {
				// join user to the team
				teammdl.Join(c, teamRequest.TeamId, teamRequest.UserId);
			} else {
				appengine.NewContext(r).Errorf(" cannot find team request with id=%d", requestId)
			}
		}
		
		teamrequestmdl.Destroy(c, requestId)
		
		url := fmt.Sprintf("/m/users/%d", auth.CurrentUser(r, c).Id)
		http.Redirect(w, r, url, http.StatusFound)
	}
}
