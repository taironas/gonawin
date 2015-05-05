/*
 * Copyright (c) 2014 Santiago Arias | Remy Jourde
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

// Package teams provides the JSON handlers to handle teams data in gonawin app.
//
// It provides the following methods
//
//	GET	/j/teams/				Retreives all teams.
//	POST	/j/teams/new/				Creates a new team.
//	GET	/j/teams/show/[0-9]+/			Retreives the team with the given id.
//	POST	/j/teams/update/[0-9]+/			Updates the team with the given id.
//	POST	/j/teams/destroy/[0-9]+/		Destroys the team with the given id.
//	POST	/j/teams/invite/[0-9]+/			Request an invitation to a private team with the given id.
//	POST	/j/teams/allow/[0-9]+/			Allow a user to be a member of a team with the given id.
//	POST	/j/teams/deny/[0-9]+/			Deny entrance of user to be a member of a team with the given id.
//	POST	/j/teams/join/[0-9]+/			Make a user join a team with the given id.
//	GET	/j/teams/search/			Search for all teams respecting the query "q"
//	GET	/j/teams/[0-9]+/members/		Retreives all members of a team with the given id.
//	GET	/j/teams/[0-9]+/ranking/		Retreives the ranking of a team with the given id.
//	GET	/j/teams/[0-9]+/accuracies/		Retreives all the tournament accuracies of a team with the given id.
//	GET	/j/teams/[0-9]+/accuracies/[0-9]+/	Retreives accuracies of a team with the given id for the specified tournament.
//
//
// Every method below gives more information about every API call, its parameters, and its resutls.
package teams

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"appengine"

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

type TeamData struct {
	Name        string
	Description string
	Visibility  string
}

type PriceData struct {
	Description string
}

// Index handler, use it to get the team data.
//      GET     /j/teams/?			List users not joined by user.
// Parameters:
//   'page' a int indicating the page number.
//   'count' a int indicating the number of teams per page number. default value is 25
// Reponse: array of JSON formatted teams.
//
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "teams index handler:"

	// get count parameter, if not present count is set to 20
	strcount := r.FormValue("count")
	count := int64(25)
	if len(strcount) > 0 {
		if n, err := strconv.ParseInt(strcount, 0, 64); err != nil {
			log.Errorf(c, "%s: error during conversion of count parameter: %v", desc, err)
		} else {
			count = n
		}
	}

	// get page parameter, if not present set page to the first one.
	strpage := r.FormValue("page")
	page := int64(1)
	if len(strpage) > 0 {
		if p, err := strconv.ParseInt(strpage, 0, 64); err != nil {
			log.Errorf(c, "%s error during conversion of page parameter: %v", desc, err)
			page = 1
		} else {
			page = p
		}
	}
	// fetch teams
	teams := mdl.GetNotJoinedTeams(c, u, count, page)

	if len(teams) == 0 {
		return templateshlp.RenderEmptyJsonArray(w, c)
	}

	type team struct {
		Id           int64  `json:",omitempty"`
		Name         string `json:",omitempty"`
		AdminIds     []int64
		Private      bool
		Accuracy     float64
		MembersCount int64
		ImageURL     string
	}
	ts := make([]team, len(teams))
	for i, t := range teams {
		ts[i].Id = t.Id
		ts[i].Name = t.Name
		ts[i].Private = t.Private
		ts[i].Accuracy = t.Accuracy
		ts[i].MembersCount = t.MembersCount
		ts[i].ImageURL = helpers.TeamImageURL(t.Name, t.Id)
	}

	return templateshlp.RenderJson(w, c, ts)
}

// New handler, use it to create a new team.
//	POST	/j/teams/new/				Creates a new team.
//
func New(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team New Handler:"

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
	}

	var tData TeamData
	err = json.Unmarshal(body, &tData)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
	}

	if len(tData.Name) <= 0 {
		log.Errorf(c, "%s 'Name' field cannot be empty", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeNameCannotBeEmpty)}
	}

	if t := mdl.FindTeams(c, "KeyName", helpers.TrimLower(tData.Name)); t != nil {
		log.Errorf(c, "%s That team name already exists.", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamAlreadyExists)}
	}

	team, err := mdl.CreateTeam(c, tData.Name, tData.Description, u.Id, tData.Visibility == "Private")
	if err != nil {
		log.Errorf(c, "%s error when trying to create a team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
	}
	// join the team
	if err = team.Join(c, u); err != nil {
		log.Errorf(c, "%s error when trying to create a team relationship: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotCreate)}
	}
	// publish new activity
	if updatedUser, err := mdl.UserById(c, u.Id); err != nil {
		log.Errorf(c, "User not found %v", u.Id)
	} else {
		updatedUser.Publish(c, "team", "created a new team", team.Entity(), mdl.ActivityEntity{})
	}

	// return the newly created team
	var tJson mdl.TeamJson
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

	msg := fmt.Sprintf("The team %s was correctly created!", team.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Team        mdl.TeamJson
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// Show handler, use it to get the team data to show.
//	GET	/j/teams/show/[0-9]+/			Retreives the team with the given id.
//
func Show(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Show Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	// get data for json team
	// build team json
	var tJson mdl.TeamJson
	fieldsToKeep := []string{"Id", "Name", "Description", "AdminIds", "Private", "TournamentIds", "Accuracy"}
	helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

	// build players json
	players := team.Players(c)
	type player struct {
		Id       int64  `json:",omitempty"`
		Username string `json:",omitempty"`
		Alias    string
		Score    int64
		ImageURL string
	}
	ps := make([]player, len(players))
	for i, p := range players {
		ps[i].Id = p.Id
		ps[i].Username = p.Username
		ps[i].Alias = p.Alias
		ps[i].Score = p.Score
		ps[i].ImageURL = helpers.UserImageURL(p.Name, p.Id)
	}

	// build tournaments json
	tournaments := team.Tournaments(c)
	type tournament struct {
		Id                int64  `json:",omitempty"`
		Name              string `json:",omitempty"`
		ParticipantsCount int
		TeamsCount        int
		Progress          float64
		ImageURL          string
	}
	ts := make([]tournament, len(tournaments))
	for i, t := range tournaments {
		ts[i].Id = t.Id
		ts[i].Name = t.Name
		ts[i].ParticipantsCount = len(t.UserIds)
		ts[i].TeamsCount = len(t.TeamIds)
		ts[i].Progress = t.Progress(c)
		ts[i].ImageURL = helpers.TournamentImageURL(t.Name, t.Id)
	}

	teamData := struct {
		Team        mdl.TeamJson
		Joined      bool
		RequestSent bool
		Players     []player
		Tournaments []tournament
		ImageURL    string
	}{
		tJson,
		team.Joined(c, u),
		mdl.WasTeamRequestSent(c, team.Id, u.Id),
		ps,
		ts,
		helpers.TeamImageURL(team.Name, team.Id),
	}
	return templateshlp.RenderJson(w, c, teamData)

}

// Update handler, use it to update a team from a given id.
//	POST	/j/teams/update/[0-9]+/			Updates the team with the given id.
//
func Update(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Update Handler:"
	extract := extract.NewContext(c, desc, r)
	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	if !mdl.IsTeamAdmin(c, team.Id, u.Id) {
		log.Errorf(c, "%s user is not admin", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamUpdateForbiden)}
	}

	// only work on name and private. Other values should not be editable
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s Error when reading request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	var updatedData TeamData
	err = json.Unmarshal(body, &updatedData)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	updatedPrivate := updatedData.Visibility == "private"
	log.Errorf(c, "%s %v. %v", desc, team.Name, team.Private)
	log.Errorf(c, "%s %v. %v", desc, updatedData.Name, updatedPrivate)
	log.Errorf(c, "%s visibility %v.", desc, updatedData.Visibility)
	log.Errorf(c, "%s updateddata %v.", desc, updatedData)

	if helpers.IsStringValid(updatedData.Name) &&
		(updatedData.Name != team.Name || updatedData.Description != team.Description || updatedPrivate != team.Private) {
		if updatedData.Name != team.Name {
			// be sure that a team with that name does not exist in datastore.
			if t := mdl.FindTeams(c, "KeyName", helpers.TrimLower(updatedData.Name)); t != nil {
				log.Errorf(c, "%s That team name already exists.", desc)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamAlreadyExists)}
			}
			// update data
			team.Name = updatedData.Name
		}
		team.Description = updatedData.Description
		team.Private = updatedPrivate
		team.Update(c)
	} else {
		log.Errorf(c, "%s Cannot update because updated is not valid.", desc)
		log.Errorf(c, "%s Update name = %s", desc, updatedData.Name)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	// publish new activity
	u.Publish(c, "team", "updated team", team.Entity(), mdl.ActivityEntity{})

	// keep only needed fields for json api
	var tJson mdl.TeamJson
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

	msg := fmt.Sprintf("The team %s was correctly updated!", team.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Team        mdl.TeamJson
	}{
		msg,
		tJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Destroy handler, use to to destroy a team.
//	POST	/j/teams/destroy/[0-9]+/		Destroys the team with the given id.
//
func Destroy(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Destroy Handler:"
	extract := extract.NewContext(c, desc, r)
	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	if !mdl.IsTeamAdmin(c, team.Id, u.Id) {
		log.Errorf(c, "%s user is not admin", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamDeleteForbiden)}
	}

	// delete all team-user relationships
	for _, player := range team.Players(c) {
		if err := player.RemoveTeamId(c, team.Id); err != nil {
			log.Errorf(c, "%s error when trying to destroy team relationship: %v", desc, err)
		} else if u.Id == player.Id {
			// Be sure that current user has the latest data,
			// as the u.Publish method will update again the user,
			// we don't want to override the team ID removal.
			u = player
		}
	}

	// delete all tournament-team relationships
	for _, tournament := range team.Tournaments(c) {
		if err := tournament.RemoveTeamId(c, team.Id); err != nil {
			log.Errorf(c, "%serror when trying to destroy tournament relationship: %v", desc, err)
		}
	}
	// delete the team
	team.Destroy(c)

	// publish new activity
	u.Publish(c, "team", "deleted team", team.Entity(), mdl.ActivityEntity{})

	msg := fmt.Sprintf("The team %s was correctly deleted!", team.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
	}{
		msg,
	}

	// return destroyed status
	return templateshlp.RenderJson(w, c, data)
}

// RequestInvite handler, use it to request an invitation to a team.
// A user sends a team a request to join a team if team is private.
// use this handler when you wish to request an invitation to a team.
// this is done when the team in set as 'private' and the user wishes to join it.
func RequestInvite(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Request Invite Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	if _, err := mdl.CreateTeamRequest(c, team.Id, team.Name, u.Id, u.Username); err != nil {
		log.Errorf(c, "%s teams.Invite, error when trying to create a team request: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotInvite)}
	}

	// return destroyed status
	return templateshlp.RenderJson(w, c, "team request was created")
}

// SendInvite handler, use it to send an invitation to gonawin.
// A team sends a user an invitation with tuple (user, team) to join a team.
func SendInvite(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Send User Invitation Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	var user *mdl.User
	user, err = extract.User()
	if err != nil {
		return err
	}

	if _, err := mdl.CreateUserRequest(c, team.Id, user.Id); err != nil {
		log.Errorf(c, "%s teams.SendInvite, error when trying to create a user request: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotInvite)}
	}

	// publish new activity
	user.Publish(c, "invitation", "has been invited to join team ", team.Entity(), mdl.ActivityEntity{})

	return templateshlp.RenderJson(w, c, "user request was created")
}

// Invited handler, use it to get all the users who were invited to a team.
func Invited(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	desc := "Team Invited Handler:"
	c := appengine.NewContext(r)
	extract := extract.NewContext(c, desc, r)

	var teamId int64
	var err error
	teamId, err = extract.TeamId()
	if err != nil {
		return err
	}

	urs := mdl.FindUserRequests(c, "TeamId", teamId)
	var ids []int64
	for _, ur := range urs {
		ids = append(ids, ur.UserId)
	}

	users := mdl.UsersByIds(c, ids)

	// filter team information to return in json api
	fieldsToKeep := []string{"Id", "Username", "Alias", "Score"}
	usersJson := make([]mdl.UserJson, len(users))
	helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

	data := struct {
		Users []mdl.UserJson
	}{
		usersJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// AllowRequest handler, use it to allow a user to join a team.
// use this handler to allow a request send by a user on a team.
// after this, the user that send the request will be part of the team
func AllowRequest(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Allow Request Handler:"
	extract := extract.NewContext(c, desc, r)

	var teamRequest *mdl.TeamRequest
	var err error
	teamRequest, err = extract.TeamRequest()
	if err != nil {
		return err
	}

	// join user to the team
	var team *mdl.Team
	team, err = mdl.TeamById(c, teamRequest.TeamId)
	if err != nil {
		log.Errorf(c, "%s team not found. id: %v, err: %v", desc, teamRequest.TeamId, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamRequestNotFound)}
	}
	user, err := mdl.UserById(c, teamRequest.UserId)
	if err != nil {
		log.Errorf(c, "%s user not found, err: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
	}

	team.Join(c, user)
	// request is no more needed so clear it from datastore
	teamRequest.Destroy(c)

	return templateshlp.RenderJson(w, c, "team request was handled")
}

// Deny handler.
// use this handler to deny a request send by a user on a team.
// the user will not be able to be part of the team
func DenyRequest(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Deny Request Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var teamRequest *mdl.TeamRequest

	teamRequest, err = extract.TeamRequest()
	if err != nil {
		return err
	}

	// request is no more needed so clear it from datastore
	teamRequest.Destroy(c)

	return templateshlp.RenderJson(w, c, "team request was handled")
}

// Search handler, use it to get all teams that match the search.
//	GET	/j/teams/search/			Search for all teams respecting the query "q"
//
func Search(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	keywords := r.FormValue("q")
	if r.Method != "GET" || len(keywords) == 0 {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Search Handler:"

	words := helpers.SetOfStrings(keywords)
	ids, err := mdl.GetTeamInvertedIndexes(c, words)
	if err != nil {
		log.Errorf(c, "%s teams.Index, error occurred when getting indexes of words: %v", desc, err)
		data := struct {
			MessageDanger string `json:",omitempty"`
		}{
			"Oops! something went wrong, we are unable to perform search query.",
		}
		return templateshlp.RenderJson(w, c, data)
	}
	result := mdl.TeamScore(c, keywords, ids)
	log.Infof(c, "%s result from TeamScore: %v", desc, result)
	teams := mdl.TeamsByIds(c, result)
	log.Infof(c, "%s ByIds result %v", desc, teams)
	if len(teams) == 0 {
		msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "team")
		data := struct {
			MessageInfo string `json:",omitempty"`
		}{
			msg,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	// filter team information to return in json api
	type team struct {
		Id           int64
		Name         string
		AdminIds     []int64
		Private      bool
		Accuracy     float64
		MembersCount int64
		ImageURL     string
	}

	ts := make([]team, len(teams))
	for i, t := range teams {
		ts[i].Id = t.Id
		ts[i].Name = t.Name
		ts[i].AdminIds = t.AdminIds
		ts[i].Private = t.Private
		ts[i].Accuracy = t.Accuracy
		ts[i].MembersCount = t.MembersCount
		ts[i].ImageURL = helpers.TeamImageURL(t.Name, t.Id)
	}

	// we should not directly return an array. so we add an extra layer.
	data := struct {
		Teams []team `json:",omitempty"`
	}{
		ts,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Members handler, use it to get all members of a team.
//	/j/teams/[0-9]+/members/	GET			use this handler to get members of a team.
//
func Members(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Members Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error

	team, err = extract.Team()
	if err != nil {
		return err
	}

	// build members json
	members := team.Players(c)
	fieldsToKeepForMember := []string{"Id", "Username", "Alias", "Score"}
	membersJson := make([]mdl.UserJson, len(members))
	helpers.TransformFromArrayOfPointers(&members, &membersJson, fieldsToKeepForMember)

	data := struct {
		Members []mdl.UserJson
	}{
		membersJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Prices handler, use it to get the team's prices.
func Prices(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Prices Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error

	team, err = extract.Team()
	if err != nil {
		return err
	}

	log.Infof(c, "%s ready to build a price array", desc)
	prices := team.Prices(c)

	data := struct {
		Prices []*mdl.Price
	}{
		prices,
	}

	return templateshlp.RenderJson(w, c, data)
}

// PriceByTournament, use it to list the prices of a team for a specific tournament.
//
// Use this handler to get the price of a team for a specific tournament.
//	GET	/j/teams/[0-9]+/prices/[0-9]+/	Retreives price of a team with the given id for the specified tournament.
//
func PriceByTournament(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method == "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Prices by tournament Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	t, err = extract.Team()
	if err != nil {
		return err
	}

	var tournamentId int64
	tournamentId, err = extract.TournamentId()
	if err != nil {
		return err
	}

	log.Infof(c, "%s ready to get price", desc)
	p := t.PriceByTournament(c, tournamentId)

	data := struct {
		Price *mdl.Price
	}{
		p,
	}
	return templateshlp.RenderJson(w, c, data)
}

// UpdatePrice handler, use it to update the price of a team for a specific tournament.
//
// Use this handler to get the price of a team for a specific tournament.
//	GET	/j/teams/[0-9]+/prices/update/[0-9]+/	Update Retreives price of a team with the given id for the specified tournament.
//
func UpdatePrice(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team update price Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	t, err = extract.Team()
	if err != nil {
		return err
	}

	var tournamentId int64
	tournamentId, err = extract.TournamentId()
	if err != nil {
		return err
	}

	log.Infof(c, "%s ready to get price", desc)
	p := t.PriceByTournament(c, tournamentId)

	// only work on name and private. Other values should not be editable
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s Error when reading request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	var priceData PriceData
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	if helpers.IsStringValid(priceData.Description) && (p.Description != priceData.Description) {
		// update data
		p.Description = priceData.Description
		p.Update(c)
	} else {
		log.Errorf(c, "%s Cannot update because updated is not valid.", desc)
		log.Errorf(c, "%s Update description = %s", desc, priceData.Description)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTeamCannotUpdate)}
	}

	data := struct {
		Price *mdl.Price
	}{
		p,
	}
	return templateshlp.RenderJson(w, c, data)
}
