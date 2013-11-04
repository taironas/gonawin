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

package tournaments

import (
	"html/template"
	"net/http"
	"fmt"
	"time"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	"github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/handlers"
	tournamentrelshlp "github.com/santiaago/purple-wing/helpers/tournamentrels"

	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	usermdl "github.com/santiaago/purple-wing/models/user"
	teammdl "github.com/santiaago/purple-wing/models/team"
	searchmdl "github.com/santiaago/purple-wing/models/search"
	tournamentrelmdl "github.com/santiaago/purple-wing/models/tournamentrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
	tournamentinvmdl "github.com/santiaago/purple-wing/models/tournamentInvertedIndex"
)

type Form struct {
	Name string
	Error string
}

type indexData struct{
	Tournaments []*tournamentmdl.Tournament
	TournamentInputSearch string
}

func Index(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)


	var data indexData	
	if r.Method == "GET"{
		tournaments := tournamentmdl.FindAll(r)

		data.Tournaments =tournaments
		data.TournamentInputSearch = ""
		
	}else if r.Method == "POST"{
		query := r.FormValue("TournamentInputSearch")
		if len(query) == 0{
			http.Redirect(w, r, "teams", http.StatusFound)
			return
		}

		words := helpers.SetOfStrings(query)
		ids := tournamentinvmdl.GetIndexes(r, words)
		c.Infof("pw: search: %v Ids:%v", query, ids)
		result := searchmdl.TournamentScore(r, query, ids)
		
		tournaments := tournamentmdl.ByIds(r, result)
		data.Tournaments = tournaments
		data.TournamentInputSearch = query

	}else{
		helpers.Error404(w)
	}
	
	t := template.Must(template.New("tmpl_tournament_index").
		ParseFiles("templates/tournament/index.html"))

	funcs := template.FuncMap{
		"Tournaments": func() bool {return true},
	}
	
	templateshlp.Render_with_data(w, r, t, data, funcs, "renderTournamentIndex")
}

func New(w http.ResponseWriter, r *http.Request){
	
	var form Form
	if r.Method == "GET" {
		form.Name = ""
		form.Error = ""
	} else if r.Method == "POST" {
		form.Name = r.FormValue("Name")
		
		if len(form.Name) <= 0 {
			form.Error = "'Name' field cannot be empty"
		} else if t := tournamentmdl.Find(r, "KeyName", helpers.TrimLower(form.Name)); t != nil {
			form.Error = "That tournament name already exists."
		} else {
			tournament := tournamentmdl.Create(r, form.Name, "description foo",time.Now(),time.Now(), auth.CurrentUser(r).Id)
			// redirect to the newly created tournament page
			http.Redirect(w, r, "/m/tournaments/" + fmt.Sprintf("%d", tournament.Id), http.StatusFound)
		}
	} else {
		helpers.Error404(w)
	}
	
	t := template.Must(template.New("tmpl_tournament_new").
		ParseFiles("templates/tournament/new.html"))
	
	funcs := template.FuncMap{}
	
	templateshlp.Render_with_data(w, r, t, form, funcs, "renderTournamentNew")
}

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/tournaments/", http.StatusFound)
		return
	}
	
	if r.Method == "POST" && r.FormValue("Action") == "delete" {
		// delete all tournament-user relationships
		for _, participant := range tournamentrelshlp.Participants(r, intID) {
			if err := tournamentrelmdl.Destroy(r, intID, participant.Id); err !=nil {
			c.Errorf("pw: error when trying to destroy tournament relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, team := range tournamentrelshlp.Teams(r, intID) {
			if err := tournamentteamrelmdl.Destroy(r, intID, team.Id); err !=nil {
			c.Errorf("pw: error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the tournament
		tournamentmdl.Destroy(r, intID)
		
		http.Redirect(w, r, "/m/tournaments", http.StatusFound)
	} else {
		helpers.Error404(w)
		return
	}
	
	funcs := template.FuncMap{
		"Joined": func() bool { return tournamentmdl.Joined(r, intID, auth.CurrentUser(r).Id) },
		"TeamJoined": func(teamId int64) bool { return tournamentmdl.TeamJoined(r, intID, teamId) },
		"IsTournamentAdmin": func() bool { return tournamentmdl.IsTournamentAdmin(r, intID, auth.CurrentUser(r).Id) },
	}
	
	t := template.Must(template.New("tmpl_tournament_show").
		Funcs(funcs).
		ParseFiles("templates/tournament/show.html",
		"templates/tournament/participants.html",
		"templates/tournament/teams.html",
		"templates/tournament/candidateTeams.html"))

	var tournament *tournamentmdl.Tournament
	tournament, err = tournamentmdl.ById(r, intID)
	
	if err != nil{
		helpers.Error404(w)
		return
	}
	
	participants := tournamentrelshlp.Participants(r, intID)
	teams := tournamentrelshlp.Teams(r, intID)
	candidateTeams := teammdl.Find(r, "AdminId", auth.CurrentUser(r).Id)
	
	tournamentData := struct { 
		Tournament *tournamentmdl.Tournament
		Participants []*usermdl.User
		Teams []*teammdl.Team
		CandidateTeams []*teammdl.Team
	}{
		tournament,
		participants,
		teams,
		candidateTeams,
	}
	templateshlp.Render_with_data(w, r, t, tournamentData, funcs, "renderTournamentShow")
}

func Edit(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r,3)
	if err != nil{
		http.Redirect(w,r, "/m/tournaments/", http.StatusFound)
		return
	}
	
	if !tournamentmdl.IsTournamentAdmin(r, intID, auth.CurrentUser(r).Id) {
		http.Redirect(w, r, "/m", http.StatusFound)
		return
	}

	var tournament *tournamentmdl.Tournament
	tournament, err = tournamentmdl.ById(r, intID)
	if err != nil{
		c.Errorf("pw: Tournament Edit handler: tournament not found. id: %v",intID)
		helpers.Error404(w)
		return
	}
		
	if r.Method == "GET" {

		funcs := template.FuncMap{}
		
		t := template.Must(template.New("tmpl_tournament_edit").
			ParseFiles("templates/tournament/edit.html"))

		templateshlp.Render_with_data(w, r, t, tournament, funcs, "renderTournamentEdit")

	}else if r.Method == "POST"{
		
		// only work on name other values should not be editable
		editName := r.FormValue("Name")

		if helpers.IsStringValid(editName) && editName != tournament.Name{
			tournament.Name = editName
			tournamentmdl.Update(r, intID, tournament)
		}else{
			c.Errorf("pw: cannot update %v", helpers.IsStringValid(editName))
		}
		url := fmt.Sprintf("/m/tournaments/%d",intID)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		helpers.Error404(w)
	}
}
