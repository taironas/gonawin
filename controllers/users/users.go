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

// Package users provides the JSON handlers to handle users data in gonawin app.
package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"appengine"
	"appengine/taskqueue"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"
	mdl "github.com/taironas/gonawin/models"
)

// Index user handler, returns an http response with the information of the
// current user.
//
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)

	users := mdl.FindAllUsers(c)

	vm := buildIndexUsersViewModel(users)

	return templateshlp.RenderJson(w, c, vm)
}

func buildIndexUsersViewModel(users []*mdl.User) []mdl.UserJson {
	fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email", "Created"}
	usersJSON := make([]mdl.UserJson, len(users))
	helpers.TransformFromArrayOfPointers(&users, &usersJSON, fieldsToKeep)
	return usersJSON
}

// Show User handler, use it to get the user JSON data.
// including parameter: {teams, tournaments, teamrequests}
// 'count' parameter: default 25
// 'page' parameter: default 1
//
func Show(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	extract := extract.NewContext(c, "User show handler:", r)

	var user *mdl.User
	var err error
	if user, err = extract.User(); err != nil {
		return err
	}

	log.Infof(c, "User: %v", user)
	log.Infof(c, "User: %v", user.TeamIds)

	with := r.FormValue("including")
	params := helpers.SetOfStrings(with)

	teams := extractTeams(c, extract, user, params)
	teamRequests := extractTeamRequests(c, teams, params)
	tournaments := extractTournaments(c, user, params)
	invitations := extractInvitations(c, user, params)

	shvm := buildShowViewModel(c, user, teams, tournaments, teamRequests, invitations)

	return templateshlp.RenderJson(w, c, shvm)
}

type showViewModel struct {
	User            mdl.UserJson                   `json:",omitempty"`
	Teams           []showTeamViewModel            `json:",omitempty"`
	TeamRequests    []mdl.TeamRequestJSON          `json:",omitempty"`
	Tournaments     []showTournamentViewModel      `json:",omitempty"`
	TournamentStats []showTournamentStatsViewModel `json:",omitempty"`
	Invitations     []mdl.TeamJSON                 `json:",omitempty"`
	ImageURL        string                         `json:",omitempty"`
}

func buildShowViewModel(c appengine.Context, u *mdl.User, teams []*mdl.Team, tournaments []*mdl.Tournament, trs []*mdl.TeamRequest, invs []*mdl.Team) showViewModel {

	uvm := buildShowUserViewModel(u)
	tvm := buildShowTeamViewModel(teams)
	tsvm := buildShowTournamentStatsViewModel(c, tournaments)
	tourvm := buildShowTournamentViewModel(tournaments)
	trsvm := buildShowTeamRequestsViewModel(trs)
	ivm := buildShowInvitationsViewModel(invs)

	// imageURL
	imageURL := helpers.UserImageURL(u.Username, u.Id)

	return showViewModel{
		uvm,
		tvm,
		trsvm,
		tourvm,
		tsvm,
		ivm,
		imageURL,
	}
}

func buildShowUserViewModel(user *mdl.User) (u mdl.UserJson) {
	fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email", "Created", "IsAdmin", "Auth", "TeamIds", "TournamentIds", "Score"}

	helpers.InitPointerStructure(user, &u, fieldsToKeep)
	return
}

func extractTeams(c appengine.Context, extract extract.Context, user *mdl.User, params []string) (teams []*mdl.Team) {

	for _, param := range params {
		if param != "teams" {
			continue
		}
		count := extract.CountOrDefault(25)
		page := extract.Page()
		teams = user.TeamsByPage(c, count, page)
		break
	}
	return
}

func extractTeamRequests(c appengine.Context, teams []*mdl.Team, params []string) (teamRequests []*mdl.TeamRequest) {

	for _, param := range params {
		if param != "teamrequests" {
			continue
		}
		teamRequests = mdl.TeamsRequests(c, teams)
		break
	}
	return
}

func extractTournaments(c appengine.Context, user *mdl.User, params []string) (tournaments []*mdl.Tournament) {

	for _, param := range params {
		if param != "tournaments" {
			continue
		}
		tournaments = user.Tournaments(c)
		break
	}
	return
}

func extractInvitations(c appengine.Context, user *mdl.User, params []string) (invitations []*mdl.Team) {

	for _, param := range params {
		if param != "invitations" {
			continue
		}
		invitations = user.Invitations(c)
		break
	}
	return
}

type showTeamViewModel struct {
	ID           int64 `json:"Id"`
	Name         string
	Accuracy     float64
	MembersCount int64
	Private      bool
	ImageURL     string
}

func buildShowTeamViewModel(teams []*mdl.Team) []showTeamViewModel {
	ts := make([]showTeamViewModel, len(teams))
	for i, t := range teams {
		ts[i].ID = t.ID
		ts[i].Name = t.Name
		ts[i].MembersCount = t.MembersCount
		ts[i].Private = t.Private
		ts[i].ImageURL = helpers.TeamImageURL(t.Name, t.ID)
	}
	return ts
}

type showTournamentStatsViewModel struct {
	ID                int64 `json:"Id"`
	Name              string
	ParticipantsCount int
	TeamsCount        int
	Progress          float64
	ImageURL          string
}

func buildShowTournamentStatsViewModel(c appengine.Context, tournaments []*mdl.Tournament) []showTournamentStatsViewModel {

	stats := make([]showTournamentStatsViewModel, len(tournaments))
	for i, t := range tournaments {
		stats[i].ID = t.Id
		stats[i].Name = t.Name
		stats[i].ParticipantsCount = len(t.UserIds)
		stats[i].TeamsCount = len(t.TeamIds)
		stats[i].Progress = t.Progress(c)
		stats[i].ImageURL = helpers.TournamentImageURL(t.Name, t.Id)
	}
	return stats
}

type showTournamentViewModel struct {
	ID       int64 `json:"Id"`
	Name     string
	UserIds  []int64
	TeamIds  []int64
	ImageURL string
}

func buildShowTournamentViewModel(tournaments []*mdl.Tournament) []showTournamentViewModel {

	tournaments2 := make([]showTournamentViewModel, len(tournaments))
	for i, t := range tournaments {
		tournaments2[i].ID = t.Id
		tournaments2[i].UserIds = t.UserIds
		tournaments2[i].TeamIds = t.TeamIds
		tournaments2[i].ImageURL = helpers.TournamentImageURL(t.Name, t.Id)
	}
	return tournaments2
}

func buildShowTeamRequestsViewModel(teamRequests []*mdl.TeamRequest) []mdl.TeamRequestJSON {

	fieldsToKeep := []string{"Id", "TeamId", "TeamName", "UserId", "UserName"}
	trs := make([]mdl.TeamRequestJSON, len(teamRequests))
	helpers.TransformFromArrayOfPointers(&teamRequests, &trs, fieldsToKeep)
	return trs
}

func buildShowInvitationsViewModel(invitations []*mdl.Team) []mdl.TeamJSON {

	fieldsToKeep := []string{"Id", "Name"}
	inv := make([]mdl.TeamJSON, len(invitations))
	helpers.TransformFromArrayOfPointers(&invitations, &inv, fieldsToKeep)
	return inv
}

// use this structure to get information of user in order to update it.
type userData struct {
	User struct {
		Username string
		Name     string
		Alias    string
		Email    string
	}
}

// Update user handler, use this handler to update a user entity.
//
func Update(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User update handler:"
	extract := extract.NewContext(c, desc, r)

	var userID int64
	var err error

	if userID, err = extract.UserID(); err != nil {
		return err
	} else if userID != u.Id {
		log.Errorf(c, "%s error user ids do not match. url id:%s user id: %s", desc, userID, u.Id)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
	}

	var updatedData *userData
	if updatedData, err = userDataFromHTTPRequest(c, desc, r); err != nil {
		return err
	}

	var update bool
	if shouldUpdateUserEmail(updatedData.User.Email, u.Email) {
		u.Email = updatedData.User.Email
		update = true
	}

	if shouldUpdateUserAlias(updatedData.User.Alias, u.Alias) {
		u.Alias = updatedData.User.Alias
		update = true
	}

	if !update {
		return nothingToUpdate(c, w)
	}

	if err = u.Update(c); err != nil {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
	}

	uvm := buildUpdateViewModel(u)

	return templateshlp.RenderJson(w, c, uvm)
}

type updateViewModel struct {
	MessageInfo string       `json:",omitempty"`
	User        mdl.UserJson `json:",omitempty"`
}

func buildUpdateViewModel(u *mdl.User) updateViewModel {

	fieldsToKeep := []string{"Id", "Username", "Name", "Alias", "Email"}
	var uJSON mdl.UserJson
	helpers.InitPointerStructure(u, &uJSON, fieldsToKeep)

	return updateViewModel{
		"User was correctly updated.",
		uJSON,
	}
}

func shouldUpdateUserEmail(new, old string) bool {
	return helpers.IsEmailValid(new) && new != old
}

func shouldUpdateUserAlias(new, old string) bool {
	return helpers.IsStringValid(new) && new != old
}

func userDataFromHTTPRequest(c appengine.Context, desc string, r *http.Request) (*userData, error) {

	// only work on name other values should not be editable
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s Error when reading request body err: %v.", desc, err)
		return nil, &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
	}

	var updatedData userData
	err = json.Unmarshal(body, &updatedData)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body err: %v.", desc, err)
		return nil, &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserCannotUpdate)}
	}
	return &updatedData, nil
}

func nothingToUpdate(c appengine.Context, w http.ResponseWriter) error {
	data := struct {
		MessageInfo string `json:",omitempty"`
	}{
		"Nothing to update.",
	}
	return templateshlp.RenderJson(w, c, data)

}

// Destroy hander, use this to remove a user from gonawin.
//	POST	/j/user/destroy/[0-9]+/		Destroys the user with the given id.
//
func Destroy(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User Destroy Handler:"
	extract := extract.NewContext(c, desc, r)

	var user *mdl.User
	var err error
	if user, err = extract.User(); err != nil {
		return err
	}

	if err = removeTeamUserRels(c, desc, user, u); err != nil {
		return err
	}
	if err = removeTournameUserRels(c, desc, user, u); err != nil {
		return err
	}

	sendTaskDeleteUserPredictions(c, desc, u)

	user.Destroy(c)

	dvm := buildDestroyUserViewModel(user)
	return templateshlp.RenderJson(w, c, dvm)
}

type destroyUserViewModel struct {
	MessageInfo string `json:",omitempty"`
}

func buildDestroyUserViewModel(user *mdl.User) destroyUserViewModel {
	msg := fmt.Sprintf("The user %s was correctly deleted!", user.Username)
	return destroyUserViewModel{msg}
}

// removeTeamUserRels remove all team - user relationships.
//
func removeTeamUserRels(c appengine.Context, desc string, requestUser, currentUser *mdl.User) error {
	var err error
	for _, teamID := range requestUser.TeamIds {
		if mdl.IsTeamAdmin(c, teamID, currentUser.Id) {
			var team *mdl.Team
			if team, err = mdl.TeamByID(c, teamID); err != nil {
				log.Errorf(c, "%s team %d not found", desc, teamID)
				return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
			}
			if err = team.RemoveAdmin(c, requestUser.Id); err != nil {
				log.Infof(c, "%s error occurred during admin deletion: %v", desc, err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserIsTeamAdminCannotDelete)}
			}
		} else {
			if err = requestUser.RemoveTeamId(c, teamID); err != nil {
				log.Errorf(c, "%s error when trying to destroy team relationship: %v", desc, err)
			}
		}
	}
	return nil
}

// removeTournameUserRels deletes all tournament-user relationships.
//
func removeTournameUserRels(c appengine.Context, desc string, requestUser, currentUser *mdl.User) error {

	var err error
	for _, tournamentID := range requestUser.TournamentIds {
		if mdl.IsTournamentAdmin(c, tournamentID, currentUser.Id) {
			var tournament *mdl.Tournament
			if tournament, err = mdl.TournamentById(c, tournamentID); err != nil {
				log.Errorf(c, "%s tournament %d not found", desc, tournamentID)
				return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
			}
			if err = tournament.RemoveAdmin(c, requestUser.Id); err != nil {
				log.Infof(c, "%s error occurred during admin deletion: %v", desc, err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeUserIsTournamentAdminCannotDelete)}
			}
		} else {
			if err := requestUser.RemoveTournamentId(c, tournamentID); err != nil {
				log.Errorf(c, "%s error when trying to destroy tournament relationship: %v", desc, err)
			}
		}
	}
	return nil
}

// sendTaskDeleteUserPredictions sends a task to delete predicts of the user.
//
func sendTaskDeleteUserPredictions(c appengine.Context, desc string, u *mdl.User) {

	log.Infof(c, "%s Sending to taskqueue: delete predicts", desc)

	var predictsIds []byte
	var err error
	if predictsIds, err = json.Marshal(u.PredictIds); err != nil {
		log.Errorf(c, "%s Error marshaling %v", desc, err)
	}

	task := taskqueue.NewPOSTTask("/a/publish/users/deletepredicts/", url.Values{
		"predict_ids": []string{string(predictsIds)},
	})

	if _, err = taskqueue.Add(c, task, ""); err != nil {
		log.Errorf(c, "%s unable to add task to taskqueue. %v", desc, err)
	} else {
		log.Infof(c, "%s add task to taskqueue successfully", desc)
	}
}

// Teams handler, use this to retrieve the teams of the current user.
// count parameter: default 12
// page parameter: default 1
//
func Teams(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User joined teams handler:"
	extract := extract.NewContext(c, desc, r)

	var user *mdl.User
	var err error
	if user, err = extract.User(); err != nil {
		return err
	}

	count := extract.CountOrDefault(25)
	page := extract.Page()

	var teams []*mdl.Team
	teams = user.TeamsByPage(c, count, page)

	tvm := buildTeamsUserViewModel(teams)
	return templateshlp.RenderJson(w, c, tvm)
}

type teamsUserViewModel struct {
	Teams []mdl.TeamJSON `json:",omitempty"`
}

func buildTeamsUserViewModel(teams []*mdl.Team) teamsUserViewModel {
	teamsFieldsToKeep := []string{"Id", "Name"}
	teamsJSON := make([]mdl.TeamJSON, len(teams))
	helpers.TransformFromArrayOfPointers(&teams, &teamsJSON, teamsFieldsToKeep)

	return teamsUserViewModel{teamsJSON}
}

// Tournaments user handler, use this to retrieve the JSON data of the tournaments of the user.
// count parameter: default 25
// page parameter: default 1
func Tournaments(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	extract := extract.NewContext(c, "User joined tournaments handler", r)

	var user *mdl.User
	var err error
	if user, err = extract.User(); err != nil {
		return err
	}

	count := extract.CountOrDefault(25)
	page := extract.Page()

	tournaments := user.TournamentsByPage(c, count, page)

	tvm := buildTournamentsUserViewModel(tournaments)

	return templateshlp.RenderJson(w, c, tvm)
}

type tournamentsUserViewModel struct {
	Tournaments []mdl.TournamentJson `json:",omitempty"`
}

func buildTournamentsUserViewModel(tournaments []*mdl.Tournament) tournamentsUserViewModel {
	fieldsToKeep := []string{"Id", "Name"}
	json := make([]mdl.TournamentJson, len(tournaments))
	helpers.TransformFromArrayOfPointers(&tournaments, &json, fieldsToKeep)

	return tournamentsUserViewModel{json}
}

// AllowInvitation handler, use it to allow an invitation to a team.
//
func AllowInvitation(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User allow invitation handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	if team, err = extract.Team(); err != nil {
		return err
	}

	// find user request
	var ur *mdl.UserRequest
	if ur = mdl.FindUserRequestByTeamAndUser(c, team.ID, u.Id); ur == nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// add user as member of team
	if err = team.Join(c, u); err != nil {
		log.Errorf(c, "Team Join Handler: error on Join team: %v", err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// destroy user request
	if err = ur.Destroy(c); err != nil {
		log.Errorf(c, "%s error when destroying user request. Error: %v", err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// publish activity
	if updatedUser, err := mdl.UserById(c, u.Id); err == nil && updatedUser != nil {
		updatedUser.Publish(c, "team", "joined team", team.Entity(), mdl.ActivityEntity{})
	}

	vm := buildAllowInvitationUserViewModel(team)
	return templateshlp.RenderJson(w, c, vm)
}

type allowInvitationUserViewModel struct {
	MessageInfo string `json:",omitempty"`
	Team        mdl.TeamJSON
}

func buildAllowInvitationUserViewModel(team *mdl.Team) allowInvitationUserViewModel {

	var json mdl.TeamJSON
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &json, fieldsToKeep)

	msg := fmt.Sprintf("You accepted invitation to team %s.", team.Name)
	return allowInvitationUserViewModel{
		msg,
		json,
	}
}

// DenyInvitation handler, use it to deny an invitation to a team.
//
func DenyInvitation(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "User deny invitation handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	if team, err = extract.Team(); err != nil {
		return err
	}

	var ur *mdl.UserRequest
	if ur = mdl.FindUserRequestByTeamAndUser(c, team.ID, u.Id); ur == nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// destroy user request.
	if err = ur.Destroy(c); err != nil {
		log.Errorf(c, "%s error when destroying user request. Error: %v", err)
	}

	vm := buildDenyInvitationUserViewModel(team.Name)
	return templateshlp.RenderJson(w, c, vm)
}

type denyInvitationUserViewModel struct {
	MessageInfo string `json:",omitempty"`
}

func buildDenyInvitationUserViewModel(name string) denyInvitationUserViewModel {
	return denyInvitationUserViewModel{
		fmt.Sprintf("You denied an invitation to team %s.", name),
	}
}
