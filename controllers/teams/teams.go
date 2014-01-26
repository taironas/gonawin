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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

type TeamData struct {
	Name string
	Visibility string
}

type indexData struct{
	Teams []*teammdl.Team
	TeamInputSearch string
}

// team handler
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
		return
	}
	
	t := template.Must(template.New("tmpl_team_index").
		ParseFiles("templates/team/index.html"))

	funcs := template.FuncMap{
		"Teams": func() bool {return true},
	}

	templateshlp.RenderWithData(w, r, c, t, data, funcs, "renderTeamIndex")
}

// Team new handler
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
			return
		}
	} else{
		helpers.Error404(w)
	}

	t := template.Must(template.New("tmpl_team_new").
		ParseFiles("templates/team/new.html"))
	
	funcs := template.FuncMap{}
	
	templateshlp.RenderWithData(w, r, c, t, form, funcs, "renderTeamNew")
}

// Team show handler
func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
		return
	}

	if (r.Method == "GET"){
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
}

// Team Edit handler
func Edit(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
		return
	}

	if !teammdl.IsTeamAdmin(c, intID, auth.CurrentUser(r, c).Id) {
		http.Redirect(w, r, "/m", http.StatusFound)
		return
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
		return
	} else {
		helpers.Error404(w)
	}
}

// Team destroy handler
func Destroy(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	teamID, err := handlers.PermalinkID(r, c, 4)
	if err != nil{
		http.Redirect(w,r, "/m/teams/", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		// delete all team-user relationships
		for _, player := range teamrelshlp.Players(c, teamID) {
			if err := teamrelmdl.Destroy(c, teamID, player.Id); err !=nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, tournament := range tournamentrelshlp.Teams(c, teamID) {
			if err := tournamentteamrelmdl.Destroy(c, tournament.Id, teamID); err !=nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the team
		teammdl.Destroy(c, teamID)

		http.Redirect(w, r, "/m/teams", http.StatusFound)
		return
	}
}

// invite handler
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
		return
	}
}

// Request handler
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
		return
	}
}

// json index handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		teams := teammdl.FindAll(c)
		if len(teams) == 0{
			return templateshlp.RenderEmptyJsonArray(w, c)
		}

		return templateshlp.RenderJson(w, c, teams)
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json new handler
func NewJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	
	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return helpers.InternalServerError{ errors.New("Error when reading request body") }
		}
		
		var data TeamData
		err = json.Unmarshal(body, &data)
		if err != nil {
				return helpers.InternalServerError{ errors.New("Error when decoding request body") }
		}
		
		if len(data.Name) <= 0 {
			return helpers.InternalServerError{errors.New("'Name' field cannot be empty")}
		} else if t := teammdl.Find(c, "KeyName", helpers.TrimLower(data.Name)); t != nil {
			return helpers.InternalServerError{ errors.New("That team name already exists") }
		} else {
			team, err := teammdl.Create(c, data.Name, u.Id, data.Visibility == "Private")
			if err != nil {
				log.Errorf(c, "error when trying to create a team: %v", err)
				return helpers.InternalServerError{ errors.New("error when trying to create a team") }
			}
			// join the team
			_, err = teamrelmdl.Create(c, team.Id, u.Id)
			if err != nil {
				log.Errorf(c, " error when trying to create a team relationship: %v", err)
				return helpers.InternalServerError{ errors.New("error when trying to create a team relationship") }
			}
			// return the newly created team
			return templateshlp.RenderJson(w, c, team)
		}
	} else{
		return helpers.BadRequest{ errors.New("not supported") }
	}
}

// json show handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	
	intID, err := handlers.PermalinkID(r, c, 4)
	if err != nil{
		return helpers.NotFound{err}
	}

	if (r.Method == "GET"){
		var team *teammdl.Team
		if team, err = teammdl.ById(c, intID); err != nil{
			return helpers.NotFound{err}
		}
		// get data for json team
		var teamData helpers.TeamJson
		teamData.Id = team.Id
		teamData.Name = team.Name
		teamData.Private = team.Private
		teamData.Joined = teammdl.Joined(c, intID, u.Id)
		teamData.RequestSent = 	teamrequestmdl.Sent(c, intID, u.Id)
		teamData.AdminId = team.AdminId
		// get compress players data
		playersFull := teamrelshlp.Players(c, intID)
		playersCompress := make([]helpers.UserJsonZip, len(playersFull))
		playerCounter := 0
		for _, player := range playersFull{
			playersCompress[playerCounter].Id = player.Id
			playersCompress[playerCounter].Username = player.Username
			playerCounter++
		}
		teamData.Players = playersCompress

		return templateshlp.RenderJson(w, c, teamData)
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json update handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	
	teamID, err := handlers.PermalinkID(r, c, 4)
	if err != nil{
		return helpers.NotFound{err}
	}
	
	if r.Method == "POST" {
		if !teammdl.IsTeamAdmin(c, teamID, u.Id) {
			return helpers.BadRequest{errors.New("Team can only be updated by the team administrator")}
		}
	
		var team *teammdl.Team
		team, err = teammdl.ById(c, teamID)
		if err != nil{
			log.Errorf(c, " Team Edit handler: team not found. id: %v",teamID)
			return helpers.NotFound{err}
		}
		// only work on name and private. Other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return helpers.InternalServerError{ errors.New("Error when reading request body") }
		}
		
		var updatedData TeamData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
				return helpers.InternalServerError{ errors.New("Error when decoding request body") }
		}
		
		updatedPrivate := updatedData.Visibility == "Private"

		if helpers.IsStringValid(updatedData.Name) && (updatedData.Name != team.Name || updatedPrivate != team.Private) {
			team.Name = updatedData.Name
			team.Private = updatedPrivate
			teammdl.Update(c, teamID, team)
		} else {
			log.Errorf(c, "Cannot update because updated data are not valid")
			log.Errorf(c, "Update name = %s", updatedData.Name)
		}
		
		// return the updated team
		return templateshlp.RenderJson(w, c, team)
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json destroy handler
func DestroyJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)

	teamID, err := handlers.PermalinkID(r, c, 4)
	if err != nil{
		return helpers.NotFound{err}
	}

	if r.Method == "POST" {
		if !teammdl.IsTeamAdmin(c, teamID, u.Id) {
			return helpers.BadRequest{errors.New("Team can only be deleted by the team administrator")}
		}
		
		// delete all team-user relationships
		for _, player := range teamrelshlp.Players(c, teamID) {
			if err := teamrelmdl.Destroy(c, teamID, player.Id); err !=nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, tournament := range tournamentrelshlp.Teams(c, teamID) {
			if err := tournamentteamrelmdl.Destroy(c, tournament.Id, teamID); err !=nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the team
		teammdl.Destroy(c, teamID)

		// return destroyed status
		return templateshlp.RenderJson(w, c, "team has been destroyed")
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}

// json invite handler
// use this handler when you wish to request an invitation to a team.
// this is done when the team in set as 'private' and the user wishes to join it.
func InviteJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
		
	if r.Method == "POST"{
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil{
			return helpers.NotFound{err}
		}
		
		if _, err := teamrequestmdl.Create(c, intID, u.Id); err != nil {
			log.Errorf(c, " teams.Invite, error when trying to create a team request: %v", err)
			return helpers.InternalServerError{errors.New("Error when sending invite.")}
		}
		// return destroyed status
		return templateshlp.RenderJson(w, c, "team request was created")
	}
	return helpers.NotFound{errors.New("Not supported.")}
}

// Json Allow handler
// use this handler to allow a request send by a user on a team.
// after this, the user that that send the request will be part of the team
func AllowRequestJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	
	if r.Method == "POST"{
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, " teams.AllowRequest, id could not be extracter from url: %v", err)
			return helpers.NotFound{err}
		}
		
		if teamRequest, err := teamrequestmdl.ById(c, requestId); err == nil {
			// join user to the team
			teammdl.Join(c, teamRequest.TeamId, teamRequest.UserId);
		} else {
			appengine.NewContext(r).Errorf(" cannot find team request with id=%d", requestId)
		}
		// request is no more needed so clear it from datastore
		teamrequestmdl.Destroy(c, requestId)
		
		return templateshlp.RenderJson(w, c, "team request was handled")

	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}


// Json Deny handler
// use this handler to deny a request send by a user on a team.
// the user will not be able to be part of the team
func DenyRequestJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	
	if r.Method == "POST"{
		requestId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, " teams.AllowRequest, id could not be extracter from url: %v", err)
			return helpers.NotFound{err}
		}
		
		// request is no more needed so clear it from datastore
		teamrequestmdl.Destroy(c, requestId)
		
		return templateshlp.RenderJson(w, c, "team request was handled")

	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json search handler
// use this handler to search for a team.
func SearchJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	log.Infof(c, "json search handler.")
	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0){
		words := helpers.SetOfStrings(keywords)
		ids, err := teaminvidmdl.GetIndexes(c, words)
		if err != nil {
			log.Errorf(c, " teams.Index, error occurred when getting indexes of words: %v", err)
		}
		result := searchmdl.TeamScore(c, keywords, ids)
		log.Infof(c, "result from TeamScore: %v", result)
		teams := teammdl.ByIds(c, result)
		log.Infof(c, "ByIds result %v", teams)
		if len(teams) == 0{
			// we build an array instead to returning "null" which is what the json encoder does when data is empty.
			// as angularjs expects either an array or an object, in the search case we expect and array. 
			// when there are not results found we build and empty array with a "not found" Message that should be handled by the client .
			msg := fmt.Sprintf("Your search - %s - did not match any team.", keywords)
			type msgStruct struct{
				Message string
			}
			var data [1]msgStruct
			data[0].Message = msg
			return templateshlp.RenderJson(w, c, data)
		}
		return templateshlp.RenderJson(w, c, teams)
	} else {
		return helpers.BadRequest{errors.New("not supported")}
	}
}
