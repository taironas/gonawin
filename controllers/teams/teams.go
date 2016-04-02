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
//	GET	/j/teams/				Retrieves all teams.
//	POST	/j/teams/new/				Creates a new team.
//	GET	/j/teams/show/[0-9]+/			Retrieves the team with the given id.
//	POST	/j/teams/update/[0-9]+/			Updates the team with the given id.
//	POST	/j/teams/destroy/[0-9]+/		Destroys the team with the given id.
//	POST	/j/teams/allow/[0-9]+/			Allow a user to be a member of a team with the given id.
//	POST	/j/teams/deny/[0-9]+/			Deny entrance of user to be a member of a team with the given id.
//	POST	/j/teams/join/[0-9]+/			Make a user join a team with the given id.
//	GET	/j/teams/search/			Search for all teams respecting the query "q"
//	GET	/j/teams/[0-9]+/members/		Retrieves all members of a team with the given id.
//	GET	/j/teams/[0-9]+/ranking/		Retrieves the ranking of a team with the given id.
//	GET	/j/teams/[0-9]+/accuracies/		Retrieves all the tournament accuracies of a team with the given id.
//	GET	/j/teams/[0-9]+/accuracies/[0-9]+/	Retrieves accuracies of a team with the given id for the specified tournament.
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

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// TeamData holds basic information of a Team entity.
//
type TeamData struct {
	Name        string
	Description string
	Visibility  string
}

// PriceData holds basic information of a Price entity.
//
type PriceData struct {
	Description string
}

// Index handler, use it to get the team data.
//      GET     /j/teams/?			List users not joined by user.
// Parameters:
//   'page' a int indicating the page number.
//   'count' a int indicating the number of teams per page number. default value is 25
// Response: array of JSON formatted teams.
//
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "teams index handler:"

	// get count parameter, if not present count is set to 25
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
		} else {
			page = p
		}
	}

	// fetch teams
	teams := mdl.GetNotJoinedTeams(c, u, count, page)

	tvm := buildIndexTeamsViewModel(teams)

	return templateshlp.RenderJSON(w, c, tvm)
}

type indexTeamViewModel struct {
	ID           int64 `json:"Id"`
	Name         string
	Private      bool
	Accuracy     float64
	MembersCount int64
	ImageURL     string
}

func buildIndexTeamsViewModel(teams []*mdl.Team) []indexTeamViewModel {
	ts := make([]indexTeamViewModel, len(teams))
	for i, t := range teams {
		ts[i].ID = t.ID
		ts[i].Name = t.Name
		ts[i].Private = t.Private
		ts[i].Accuracy = t.Accuracy
		ts[i].MembersCount = t.MembersCount
		ts[i].ImageURL = helpers.TeamImageURL(t.Name, t.ID)
	}

	return ts
}

// New handler, use it to create a new team.
//	POST	/j/teams/new/		Creates a new team.
// Response: JSON formatted team.
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

	team, err := mdl.CreateTeam(c, tData.Name, tData.Description, u.ID, tData.Visibility == "Private")
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
	if updatedUser, err := mdl.UserByID(c, u.ID); err != nil {
		log.Errorf(c, "User not found %v", u.ID)
	} else {
		updatedUser.Publish(c, "team", "created a new team", team.Entity(), mdl.ActivityEntity{})
	}

	// return the newly created team
	tvm := buildNewTeamsViewModel(team)

	return templateshlp.RenderJSON(w, c, tvm)
}

type newTeamViewModel struct {
	MessageInfo string `json:",omitempty"`
	Team        mdl.TeamJSON
}

func buildNewTeamsViewModel(team *mdl.Team) newTeamViewModel {
	var tJSON mdl.TeamJSON
	fieldsToKeep := []string{"ID", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &tJSON, fieldsToKeep)

	msg := fmt.Sprintf("The team %s was correctly created!", team.Name)

	tvm := newTeamViewModel{MessageInfo: msg, Team: tJSON}

	return tvm
}

// Show handler, use it to get the team data to show.
//	GET	/j/teams/show/[0-9]+/		Retrieves the team with the given id.
// Response: JSON formatted team.
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

	// build players json
	var players []*mdl.User
	if players, err = team.Players(c); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// build tournaments json
	tournaments := team.Tournaments(c)

	svm := buildShowViewModel(c, team, u, players, tournaments)

	return templateshlp.RenderJSON(w, c, svm)

}

type showViewModel struct {
	Team        mdl.TeamJSON              `json:",omitempty"`
	Joined      bool                      `json:",omitempty"`
	RequestSent bool                      `json:",omitempty"`
	Players     []playerViewModel         `json:",omitempty"`
	Tournaments []showTournamentViewModel `json:",omitempty"`
	ImageURL    string                    `json:",omitempty"`
}

func buildShowViewModel(c appengine.Context, t *mdl.Team, u *mdl.User, players []*mdl.User, tournaments []*mdl.Tournament) showViewModel {
	// build team json
	var tJSON mdl.TeamJSON
	fieldsToKeep := []string{"ID", "Name", "Description", "AdminIds", "Private", "TournamentIds", "Accuracy"}
	helpers.InitPointerStructure(t, &tJSON, fieldsToKeep)

	pvm := buildPlayersViewModel(c, players)
	tvm := buildShowTournamentViewModel(c, tournaments)

	return showViewModel{
		tJSON,
		t.Joined(c, u),
		mdl.WasTeamRequestSent(c, t.ID, u.ID),
		pvm,
		tvm,
		helpers.TeamImageURL(t.Name, t.ID),
	}
}

type playerViewModel struct {
	ID       int64 `json:"Id"`
	Username string
	Alias    string
	Score    int64
	ImageURL string
}

func buildPlayersViewModel(c appengine.Context, players []*mdl.User) []playerViewModel {
	pvm := make([]playerViewModel, len(players))
	for i, p := range players {
		pvm[i].ID = p.ID
		pvm[i].Username = p.Username
		pvm[i].Alias = p.Alias
		pvm[i].Score = p.Score
		pvm[i].ImageURL = helpers.UserImageURL(p.Name, p.ID)
	}

	return pvm
}

type showTournamentViewModel struct {
	ID                int64 `json:"Id"`
	Name              string
	ParticipantsCount int
	TeamsCount        int
	Progress          float64
	ImageURL          string
}

func buildShowTournamentViewModel(c appengine.Context, tournaments []*mdl.Tournament) []showTournamentViewModel {
	tvm := make([]showTournamentViewModel, len(tournaments))
	for i, t := range tournaments {
		tvm[i].ID = t.ID
		tvm[i].Name = t.Name
		tvm[i].ParticipantsCount = len(t.UserIds)
		tvm[i].TeamsCount = len(t.TeamIds)
		tvm[i].Progress = t.Progress(c)
		tvm[i].ImageURL = helpers.TournamentImageURL(t.Name, t.ID)
	}

	return tvm
}

// Update handler, use it to update a team from a given id.
//	POST	/j/teams/update/[0-9]+/			Updates the team with the given id.
// Response: JSON formatted team.
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

	if !mdl.IsTeamAdmin(c, team.ID, u.ID) {
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

	tvm := buildUpdateTeamsViewModel(team)

	return templateshlp.RenderJSON(w, c, tvm)
}

type updateTeamViewModel struct {
	MessageInfo string `json:",omitempty"`
	Team        mdl.TeamJSON
}

func buildUpdateTeamsViewModel(team *mdl.Team) updateTeamViewModel {
	var tJSON mdl.TeamJSON
	fieldsToKeep := []string{"ID", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &tJSON, fieldsToKeep)

	msg := fmt.Sprintf("The team %s was correctly updated!", team.Name)

	tvm := updateTeamViewModel{MessageInfo: msg, Team: tJSON}

	return tvm
}

// Destroy handler, use to to destroy a team.
//	POST	/j/teams/destroy/[0-9]+/		Destroys the team with the given id.
// Response: JSON formatted message.
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

	if !mdl.IsTeamAdmin(c, team.ID, u.ID) {
		log.Errorf(c, "%s user is not admin", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamDeleteForbiden)}
	}

	// delete all team-user relationships
	var players []*mdl.User
	if players, err = team.Players(c); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	for _, player := range players {
		if err := player.RemoveTeamID(c, team.ID); err != nil {
			log.Errorf(c, "%s error when trying to destroy team relationship: %v", desc, err)
		} else if u.ID == player.ID {
			// Be sure that current user has the latest data,
			// as the u.Publish method will update again the user,
			// we don't want to override the team ID removal.
			u = player
		}
	}

	// delete all tournament-team relationships
	for _, tournament := range team.Tournaments(c) {
		if err := tournament.RemoveTeamID(c, team.ID); err != nil {
			log.Errorf(c, "%s error when trying to destroy tournament relationship: %v", desc, err)
		}
	}
	// delete the team
	team.Destroy(c)

	// publish new activity
	u.Publish(c, "team", "deleted team", team.Entity(), mdl.ActivityEntity{})

	tvm := buildDestroyTeamsViewModel(team)

	// return destroyed status
	return templateshlp.RenderJSON(w, c, tvm)
}

type destroyTeamViewModel struct {
	MessageInfo string `json:",omitempty"`
}

func buildDestroyTeamsViewModel(team *mdl.Team) destroyTeamViewModel {
	msg := fmt.Sprintf("The team %s was correctly deleted!", team.Name)

	tvm := destroyTeamViewModel{MessageInfo: msg}

	return tvm
}

// Members handler, use it to get all members of a team.
//	/j/teams/[0-9]+/members/	GET			use this handler to get members of a team.
// Response: array of JSON formatted users.
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
	var members []*mdl.User
	if members, err = team.Players(c); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	mvm := buildMembersViewModel(c, members)

	return templateshlp.RenderJSON(w, c, mvm)
}

type memberViewModel struct {
	ID       int64 `json:"Id"`
	Username string
	Alias    string
	Score    int64
	ImageURL string
}

type membersViewModel struct {
	Members []memberViewModel
}

func buildMembersViewModel(c appengine.Context, members []*mdl.User) membersViewModel {
	mvm := make([]memberViewModel, len(members))
	for i, m := range members {
		mvm[i].ID = m.ID
		mvm[i].Username = m.Username
		mvm[i].Alias = m.Alias
		mvm[i].Score = m.Score
		mvm[i].ImageURL = helpers.UserImageURL(m.Name, m.ID)
	}

	return membersViewModel{Members: mvm}
}
