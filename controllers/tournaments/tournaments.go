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
	"errors"
	"html/template"
	"net/http"
	"fmt"
	"time"

	"appengine"	

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
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

// index tournaments handler
func Index(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	var data indexData	
	if r.Method == "GET"{
		tournaments := tournamentmdl.FindAll(c)

		data.Tournaments = tournaments
		data.TournamentInputSearch = ""
		
	} else if r.Method == "POST" {
		if query := r.FormValue("TournamentInputSearch"); len(query) == 0 {
			http.Redirect(w, r, "tournaments", http.StatusFound)
			return
		} else {
			words := helpers.SetOfStrings(query)
			ids, err := tournamentinvmdl.GetIndexes(c, words)
			
			if err != nil {
				log.Errorf(c, " tournaments.Index, error occurred when getting indexes of words: %v", err)
			}
			
			result := searchmdl.TournamentScore(c, query, ids)
			
			tournaments := tournamentmdl.ByIds(c, result)
			data.Tournaments = tournaments
			data.TournamentInputSearch = query
		}
	} else {
		helpers.Error404(w)
		return
	}
	
	t := template.Must(template.New("tmpl_tournament_index").
		ParseFiles("templates/tournament/index.html"))

	funcs := template.FuncMap{
		"Tournaments": func() bool {return true},
	}
	
	templateshlp.RenderWithData(w, r, c, t, data, funcs, "renderTournamentIndex")
}

// json index tournaments handler
func IndexJson(w http.ResponseWriter, r *http.Request) error{
	c := appengine.NewContext(r)

	var data indexData	
	if r.Method == "GET"{
		tournaments := tournamentmdl.FindAll(c)

		data.Tournaments = tournaments
		data.TournamentInputSearch = ""
		
	} else if r.Method == "POST" {
		if query := r.FormValue("TournamentInputSearch"); len(query) == 0 {
			http.Redirect(w, r, "tournaments", http.StatusFound)
			return nil
		} else {
			words := helpers.SetOfStrings(query)
			ids, err := tournamentinvmdl.GetIndexes(c, words)
			
			if err != nil {
				log.Errorf(c, " tournaments.Index, error occurred when getting indexes of words: %v", err)
			}
			
			result := searchmdl.TournamentScore(c, query, ids)
			
			tournaments := tournamentmdl.ByIds(c, result)
			data.Tournaments = tournaments
			data.TournamentInputSearch = query
		}
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
		
	return templateshlp.RenderJson(w, c, data)
}

// new tournament handler
func New(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	var form Form
	if r.Method == "GET" {
		form.Name = ""
		form.Error = ""
	} else if r.Method == "POST" {
		form.Name = r.FormValue("Name")
		
		if len(form.Name) <= 0 {
			form.Error = "'Name' field cannot be empty"
		} else if t := tournamentmdl.Find(c, "KeyName", helpers.TrimLower(form.Name)); t != nil {
			form.Error = "That tournament name already exists."
		} else {
			tournament, err := tournamentmdl.Create(c, form.Name, "description foo",time.Now(),time.Now(), auth.CurrentUser(r, c).Id)
			if err != nil {
				log.Errorf(c, " error when trying to create a tournament: %v", err)
			}
			// redirect to the newly created tournament page
			http.Redirect(w, r, "/m/tournaments/" + fmt.Sprintf("%d", tournament.Id), http.StatusFound)
			return
		}
	} else {
		helpers.Error404(w)
		return
	}
	
	t := template.Must(template.New("tmpl_tournament_new").
		ParseFiles("templates/tournament/new.html"))
	
	funcs := template.FuncMap{}
	
	templateshlp.RenderWithData(w, r, c, t, form, funcs, "renderTournamentNew")
}

// json new tournament handler
func NewJson(w http.ResponseWriter, r *http.Request) error{
	c := appengine.NewContext(r)
	
	var form Form
	if r.Method == "GET" {
		form.Name = ""
		form.Error = ""
	} else if r.Method == "POST" {
		form.Name = r.FormValue("Name")
		
		if len(form.Name) <= 0 {
			form.Error = "'Name' field cannot be empty"
		} else if t := tournamentmdl.Find(c, "KeyName", helpers.TrimLower(form.Name)); t != nil {
			form.Error = "That tournament name already exists."
		} else {
			tournament, err := tournamentmdl.Create(c, form.Name, "description foo",time.Now(),time.Now(), auth.CurrentUser(r, c).Id)
			if err != nil {
				log.Errorf(c, " error when trying to create a tournament: %v", err)
			}
			// redirect to the newly created tournament page
			http.Redirect(w, r, "/j/tournaments/" + fmt.Sprintf("%d", tournament.Id), http.StatusFound)
			return nil
		}
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
		
	return templateshlp.RenderJson(w, c, form)
}

// show tournament handler
func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		log.Errorf(c, " Unable to find ID in request %v",r)
		http.Redirect(w, r, "/m/tournaments/", http.StatusFound)
		return
	}
	
	if r.Method == "POST" && r.FormValue("Action") == "delete" {
		// delete all tournament-user relationships
		for _, participant := range tournamentrelshlp.Participants(c, intID) {
			if err := tournamentrelmdl.Destroy(c, intID, participant.Id); err !=nil {
			log.Errorf(c, " error when trying to destroy tournament relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, team := range tournamentrelshlp.Teams(c, intID) {
			if err := tournamentteamrelmdl.Destroy(c, intID, team.Id); err !=nil {
			log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the tournament
		tournamentmdl.Destroy(c, intID)
		
		http.Redirect(w, r, "/m/tournaments", http.StatusFound)
		return
	} else if r.Method != "GET"{
		log.Errorf(c, " request method not supported")
		helpers.Error404(w)
		return
	}
	
	funcs := template.FuncMap{
		"Joined": func() bool { return tournamentmdl.Joined(c, intID, auth.CurrentUser(r, c).Id) },
		"TeamJoined": func(teamId int64) bool { return tournamentmdl.TeamJoined(c, intID, teamId) },
		"IsTournamentAdmin": func() bool { return tournamentmdl.IsTournamentAdmin(c, intID, auth.CurrentUser(r, c).Id) },
	}
	
	t := template.Must(template.New("tmpl_tournament_show").
		Funcs(funcs).
		ParseFiles("templates/tournament/show.html",
		"templates/tournament/participants.html",
		"templates/tournament/teams.html",
		"templates/tournament/candidateTeams.html"))

	var tournament *tournamentmdl.Tournament
	tournament, err = tournamentmdl.ById(c, intID)
	
	if err != nil{
		helpers.Error404(w)
		return
	}
	
	participants := tournamentrelshlp.Participants(c, intID)
	teams := tournamentrelshlp.Teams(c, intID)
	candidateTeams := teammdl.Find(c, "AdminId", auth.CurrentUser(r, c).Id)
	
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
	templateshlp.RenderWithData(w, r, c, t, tournamentData, funcs, "renderTournamentShow")
}

// Json show tournament handler
func ShowJson(w http.ResponseWriter, r *http.Request) error{
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		log.Errorf(c, " Unable to find ID in request %v",r)
		return helpers.NotFound{err}
	}
	
	if r.Method == "POST" && r.FormValue("Action") == "delete" {
		// delete all tournament-user relationships
		for _, participant := range tournamentrelshlp.Participants(c, intID) {
			if err := tournamentrelmdl.Destroy(c, intID, participant.Id); err !=nil {
			log.Errorf(c, " error when trying to destroy tournament relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, team := range tournamentrelshlp.Teams(c, intID) {
			if err := tournamentteamrelmdl.Destroy(c, intID, team.Id); err !=nil {
			log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the tournament
		tournamentmdl.Destroy(c, intID)
		
		http.Redirect(w, r, "/m/tournaments", http.StatusFound)
		return nil
	} else if r.Method != "GET"{
		log.Errorf(c, " request method not supported")
		return helpers.BadRequest{errors.New("Not supported.")}
	}
	
	var tournament *tournamentmdl.Tournament
	tournament, err = tournamentmdl.ById(c, intID)	
	if err != nil{
		return helpers.NotFound{err}
	}
	
	participants := tournamentrelshlp.Participants(c, intID)
	teams := tournamentrelshlp.Teams(c, intID)
	candidateTeams := teammdl.Find(c, "AdminId", auth.CurrentUser(r, c).Id)
	
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
	return templateshlp.RenderJson(w, c, tournamentData)
}

//  Edit tournament handler
func Edit(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		http.Redirect(w,r, "/m/tournaments/", http.StatusFound)
		return
	}
	
	if !tournamentmdl.IsTournamentAdmin(c, intID, auth.CurrentUser(r, c).Id) {
		http.Redirect(w, r, "/m", http.StatusFound)
		return
	}

	var tournament *tournamentmdl.Tournament
	tournament, err = tournamentmdl.ById(c, intID)
	if err != nil{
		log.Errorf(c, " Tournament Edit handler: tournament not found. id: %v",intID)
		helpers.Error404(w)
		return
	}
		
	if r.Method == "GET" {

		funcs := template.FuncMap{}
		
		t := template.Must(template.New("tmpl_tournament_edit").
			ParseFiles("templates/tournament/edit.html"))

		templateshlp.RenderWithData(w, r, c, t, tournament, funcs, "renderTournamentEdit")
		return
	} else if r.Method == "POST" {
		
		// only work on name other values should not be editable
		editName := r.FormValue("Name")

		if helpers.IsStringValid(editName) && editName != tournament.Name{
			tournament.Name = editName
			tournamentmdl.Update(c, intID, tournament)
		} else {
			log.Errorf(c, " cannot update %v", helpers.IsStringValid(editName))
		}
		url := fmt.Sprintf("/m/tournaments/%d",intID)
		http.Redirect(w, r, url, http.StatusFound)
		return
	} else {
		helpers.Error404(w)
		return
	}
}

//  Json Edit tournament handler
func EditJson(w http.ResponseWriter, r *http.Request) error{
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		return helpers.NotFound{err}
	}
	
	if !tournamentmdl.IsTournamentAdmin(c, intID, auth.CurrentUser(r, c).Id) {
		http.Redirect(w, r, "/j", http.StatusFound)
		return nil
	}

	var tournament *tournamentmdl.Tournament
	tournament, err = tournamentmdl.ById(c, intID)
	if err != nil{
		log.Errorf(c, " Tournament Edit handler: tournament not found. id: %v",intID)
		return helpers.NotFound{err}
	}
		
	if r.Method == "GET" {
		return templateshlp.RenderJson(w, c, tournament)
	} else if r.Method == "POST" {
		
		// only work on name other values should not be editable
		editName := r.FormValue("Name")

		if helpers.IsStringValid(editName) && editName != tournament.Name{
			tournament.Name = editName
			tournamentmdl.Update(c, intID, tournament)
		} else {
			log.Errorf(c, " cannot update %v", helpers.IsStringValid(editName))
		}
		url := fmt.Sprintf("/j/tournaments/%d",intID)
		http.Redirect(w, r, url, http.StatusFound)
		return nil
	} else {
		return helpers.BadRequest{errors.New("Not supported.")}
	}
}
