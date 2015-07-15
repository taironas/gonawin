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

// Package tournaments provides the JSON handlers to handle tournaments data in gonawin app.
package tournaments

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"appengine"

	"github.com/taironas/route"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

type TournamentData struct {
	Name        string
	Description string
}

// Index handler, use it to get the data of current tournaments.
func Index(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "tournament index handler:"

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
			page = 1
		} else {
			page = p
		}
	}
	tournaments := mdl.FindAllTournaments(c, count, page)
	if len(tournaments) == 0 {
		return templateshlp.RenderEmptyJsonArray(w, c)
	}

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

	return templateshlp.RenderJson(w, c, ts)
}

// New handler, use it to create a new tournament.
func New(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
	c := appengine.NewContext(r)
	desc := "Tournament New Handler:"

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotCreate)}
	}

	var tData TournamentData
	err = json.Unmarshal(body, &tData)
	if err != nil {
		log.Errorf(c, "%s Error when decoding request body: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotCreate)}
	}

	if len(tData.Name) <= 0 {
		log.Errorf(c, "%s 'Name' field cannot be empty", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeNameCannotBeEmpty)}
	}

	if t := mdl.FindTournaments(c, "KeyName", helpers.TrimLower(tData.Name)); t != nil {
		log.Errorf(c, "%s That tournament name already exists.", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentAlreadyExists)}
	}

	tournament, err := mdl.CreateTournament(c, tData.Name, tData.Description, time.Now(), time.Now(), u.Id)
	if err != nil {
		log.Errorf(c, "%s error when trying to create a tournament: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotCreate)}
	}
	// return the newly created tournament
	fieldsToKeep := []string{"Id", "Name"}
	var tJson mdl.TournamentJson
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	u.Publish(c, "tournament", "created a tournament", tournament.Entity(), mdl.ActivityEntity{})

	msg := fmt.Sprintf("The tournament %s was correctly created!", tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJson
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)

}

// Show handler, use it to get the data of a specific tournament.
func Show(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Show Handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error
	tournament, err = extract.Tournament()
	if err != nil {
		return err
	}

	participants := tournament.Participants(c)
	teams := tournament.Teams(c)

	fieldsToKeep := []string{"Id", "Name", "Description", "AdminIds", "IsFirstStageComplete"}
	var tournamentJson mdl.TournamentJson
	helpers.InitPointerStructure(tournament, &tournamentJson, fieldsToKeep)

	participantFieldsToKeep := []string{"Id", "Username", "Alias"}
	participantsJson := make([]mdl.UserJson, len(participants))
	helpers.TransformFromArrayOfPointers(&participants, &participantsJson, participantFieldsToKeep)

	teamsJson := make([]mdl.TeamJson, len(teams))
	helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)

	progress := tournament.Progress(c)

	// formatted start and end
	const layout = "2 January 2006"
	start := tournament.Start.Format(layout)
	end := tournament.End.Format(layout)

	remainingDays := int64(tournament.Start.Sub(time.Now()).Hours() / 24)

	imageURL := helpers.TournamentImageURL(tournament.Name, tournament.Id)

	data := struct {
		Tournament    mdl.TournamentJson
		Joined        bool
		Participants  []mdl.UserJson
		Teams         []mdl.TeamJson
		Progress      float64
		Start         string
		End           string
		RemainingDays int64
		ImageURL      string
	}{
		tournamentJson,
		tournament.Joined(c, u),
		participantsJson,
		teamsJson,
		progress,
		start,
		end,
		remainingDays,
		imageURL,
	}

	return templateshlp.RenderJson(w, c, data)

}

// Destrou handler, use it to detroy a tournament.
func Destroy(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Destroy Handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error

	tournament, err = extract.Tournament()
	if err != nil {
		return err
	}

	if !mdl.IsTournamentAdmin(c, tournament.Id, u.Id) {
		log.Errorf(c, "%s user is not admin", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentDeleteForbiden)}
	}

	// delete all tournament-user relationships
	for _, participant := range tournament.Participants(c) {
		if err := participant.RemoveTournamentId(c, tournament.Id); err != nil {
			log.Errorf(c, " %s error when trying to remove tournament id from user: %v", desc, err)
		} else if u.Id == participant.Id {
			// Be sure that current user has the latest data,
			// as the u.Publish method will update again the user,
			// we don't want to override the tournament ID removal.
			u = participant
		}
	}

	// delete all tournament-team relationships
	for _, team := range tournament.Teams(c) {
		if err := tournament.TeamLeave(c, team); err != nil {
			log.Errorf(c, "%s error when trying to destroy team relationship: %v", desc, err)
		}
	}

	// delete matches of first stage
	if err := mdl.DestroyMatches(c, tournament.Matches1stStage); err != nil {
		log.Errorf(c, "%s error when trying to destroy tournament's matches of first stage: %v", desc, err)
	}

	// delete matches of second stage
	if err := mdl.DestroyMatches(c, tournament.Matches2ndStage); err != nil {
		log.Errorf(c, "%s error when trying to destroy tournament's matches of second stage: %v", desc, err)
	}

	// delete groups
	if err := mdl.DestroyGroups(c, tournament.GroupIds); err != nil {
		log.Errorf(c, "%s error when trying to destroy tournament's groups: %v", desc, err)
	}

	// delete the tournament
	tournament.Destroy(c)

	// publish new activity
	u.Publish(c, "tournament", "deleted tournament", tournament.Entity(), mdl.ActivityEntity{})

	msg := fmt.Sprintf("The tournament %s has been destroyed!", tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
	}{
		msg,
	}

	// return destroyed status
	return templateshlp.RenderJson(w, c, data)
}

//  Update hanlder, use it to update a tournament.
func Update(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Update handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error

	tournament, err = extract.Tournament()
	if err != nil {
		return err
	}

	if !mdl.IsTournamentAdmin(c, tournament.Id, u.Id) {
		log.Errorf(c, "%s user is not admin", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentUpdateForbiden)}
	}

	// only work on name other values should not be editable
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%s error when reading request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
	}

	var updatedData TournamentData
	err = json.Unmarshal(body, &updatedData)
	if err != nil {
		log.Errorf(c, "%s error when decoding request body err: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
	}

	if helpers.IsStringValid(updatedData.Name) &&
		(updatedData.Name != tournament.Name || updatedData.Description != tournament.Description) {
		if updatedData.Name != tournament.Name {
			// be sure that team with that name does not exist in datastore
			if t := mdl.FindTournaments(c, "KeyName", helpers.TrimLower(updatedData.Name)); t != nil {
				log.Errorf(c, "%s that tournament name already exists.", desc)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentAlreadyExists)}
			}
			// update data
			tournament.Name = updatedData.Name
		}
		tournament.Description = updatedData.Description
		tournament.Update(c)
	} else {
		log.Errorf(c, "%s cannot update because updated data is not valid.", desc)
		log.Errorf(c, "%s update name = %s", desc, updatedData.Name)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
	}

	// publish new activity
	u.Publish(c, "tournament", "updated tournament", tournament.Entity(), mdl.ActivityEntity{})

	// return the updated tournament
	fieldsToKeep := []string{"Id", "Name"}
	var tJson mdl.TournamentJson
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	msg := fmt.Sprintf("The tournament %s was correctly updated!", tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJson
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// Search handler, use it to get all the tournaments that match the query.
func Search(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	keywords := r.FormValue("q")

	if r.Method != "GET" || len(keywords) == 0 {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Search handler:"

	words := helpers.SetOfStrings(keywords)
	ids, err := mdl.GetTournamentInvertedIndexes(c, words)
	if err != nil {
		log.Errorf(c, "%s tournaments.Index, error occurred when getting indexes of words: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeTournamentCannotSearch)}
	}

	result := mdl.TournamentScore(c, keywords, ids)
	log.Infof(c, "%s result from TournamentScore: %v", desc, result)

	var tournaments []*mdl.Tournament
	if tournaments, err = mdl.TournamentsByIds(c, result); err != nil {
		log.Infof(c, "%v something failed when calling TournamentsByIds: %v", desc, err)
	}
	log.Infof(c, "%s ByIds result %v", desc, tournaments)

	if len(tournaments) == 0 {
		msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "tournament")
		data := struct {
			MessageInfo string `json:",omitempty"`
		}{
			msg,
		}
		return templateshlp.RenderJson(w, c, data)
	}

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

	// we should not directly return an array. so we add an extra layer.
	data := struct {
		Tournaments []tournament `json:",omitempty"`
	}{
		ts,
	}
	return templateshlp.RenderJson(w, c, data)
}

// CandidateTeams handler, use it to get the list of teams that you can add to a tournament.
func CandidateTeams(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Candidate Teams handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error
	tournament, err = extract.Tournament()
	if err != nil {
		return err
	}

	// query teams
	var teams []*mdl.Team
	for _, teamId := range u.TeamIds {
		if team, err1 := mdl.TeamById(c, teamId); err1 == nil {
			for _, aId := range team.AdminIds {
				if aId == u.Id {
					teams = append(teams, team)
				}
			}
		} else {
			log.Errorf(c, "%v", err1)
		}
	}

	type canditateType struct {
		Team   mdl.TeamJson
		Joined bool
	}

	fieldsToKeep := []string{"Id", "Name"}
	candidatesData := make([]canditateType, len(teams))

	for counterCandidate, team := range teams {
		var tJson mdl.TeamJson
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)
		var canditate canditateType
		canditate.Team = tJson
		canditate.Joined = tournament.TeamJoined(c, team)
		candidatesData[counterCandidate] = canditate
	}

	// we should not directly return an array. so we add an extra layer.
	data := struct {
		Candidates []canditateType `json:",omitempty"`
	}{
		candidatesData,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Participants handler, use it to get the participants to a tournament.
// use this handler to get participants of a tournament.
func Participants(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Participants handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error
	tournament, err = extract.Tournament()
	if err != nil {
		return err
	}

	participants := tournament.Participants(c)

	participantFieldsToKeep := []string{"Id", "Username", "Alias"}
	participantsJson := make([]mdl.UserJson, len(participants))
	helpers.TransformFromArrayOfPointers(&participants, &participantsJson, participantFieldsToKeep)

	data := struct {
		Participants []mdl.UserJson
	}{
		participantsJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// Reset handler, use it to reset points and goals of a tournament.
func Reset(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Reset handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Tournament
	var err error
	t, err = extract.Tournament()
	if err != nil {
		return err
	}

	if err = t.Reset(c); err != nil {
		log.Errorf(c, "%s unable to reset tournament: %v error:", desc, t.Id, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	groups := mdl.Groups(c, t.GroupIds)
	groupsJson := formatGroupsJson(groups)

	msg := fmt.Sprintf("Tournament is now reset.")
	data := struct {
		MessageInfo string `json:",omitempty"`
		Groups      []GroupJson
	}{
		msg,
		groupsJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Predict handler, use it to set the predictions of a match to the current user.
func Predict(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Predict Handler:"
	extract := extract.NewContext(c, desc, r)

	var tournament *mdl.Tournament
	var err error
	tournament, err = extract.Tournament()
	if err != nil {
		return err
	}

	// check if user joined the tournament
	if !tournament.Joined(c, u) {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotAllowedToSetPrediction)}
	}

	// get match id number
	strmatchIdNumber, err2 := route.Context.Get(r, "matchId")
	if err2 != nil {
		log.Errorf(c, "%s error getting match id, err:%v", desc, err2)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchNotFoundCannotSetPrediction)}
	}

	var matchIdNumber int64
	matchIdNumber, err2 = strconv.ParseInt(strmatchIdNumber, 0, 64)
	if err2 != nil {
		log.Errorf(c, "%s error converting match id from string to int64, err:%v", desc, err2)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeMatchNotFoundCannotSetPrediction)}
	}

	match := mdl.GetMatchByIdNumber(c, *tournament, matchIdNumber)
	if match == nil {
		log.Errorf(c, "%s unable to get match with id number :%v", desc, matchIdNumber)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchNotFoundCannotSetPrediction)}
	}
	result1 := r.FormValue("result1")
	result2 := r.FormValue("result2")
	var r1, r2 int
	if r1, err = strconv.Atoi(result1); err != nil {
		log.Errorf(c, "%s unable to get results, error: %v not number 1", desc, err)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeCannotSetPrediction)}
	}
	if r2, err = strconv.Atoi(result2); err != nil {
		log.Errorf(c, "%s unable to get results, error: %v not number 2", desc, err)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeCannotSetPrediction)}
	}
	msg := ""
	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(tournament); tb == nil {
		log.Errorf(c, "%s TournamentBuilder not found")
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	mapIdTeams := tb.MapOfIdTeams(c, tournament)
	var p *mdl.Predict
	if p = mdl.FindPredictByUserMatch(c, u.Id, match.Id); p == nil {
		log.Infof(c, "%s predict entity for pair (%v, %v) not found, so we create one.", desc, u.Id, match.Id)
		if predict, err1 := mdl.CreatePredict(c, u.Id, int64(r1), int64(r2), match.Id); err1 != nil {
			log.Errorf(c, "%s unable to create Predict for match with id:%v error: %v", desc, match.Id, err1)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeCannotSetPrediction)}
		} else {
			// add p.Id to User predict table.
			if err = u.AddPredictId(c, predict.Id); err != nil {
				log.Errorf(c, "%s unable to add predict id in user entity: error: %v", desc, err)
				return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeCannotSetPrediction)}
			}
			p = predict
		}
		msg = fmt.Sprintf("You set a prediction: %s %d:%d %s.", mapIdTeams[match.TeamId1], p.Result1, p.Result2, mapIdTeams[match.TeamId2])

	} else {
		// predict already exist so just update resulst.
		p.Result1 = int64(r1)
		p.Result2 = int64(r2)
		if err := p.Update(c); err != nil {
			log.Errorf(c, "%s unable to edit predict entity. %v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeCannotSetPrediction)}
		}
		msg = fmt.Sprintf("Your prediction is now updated: %s %d:%d %s.", mapIdTeams[match.TeamId1], p.Result1, p.Result2, mapIdTeams[match.TeamId2])
	}

	data := struct {
		MessageInfo string `json:",omitempty"`
		Predict     *mdl.Predict
	}{
		msg,
		p,
	}

	// publish activity
	verb := fmt.Sprintf("predicted %d-%d for", p.Result1, p.Result2)
	object := mdl.ActivityEntity{Id: match.Id, Type: "match", DisplayName: mapIdTeams[match.TeamId1] + "-" + mapIdTeams[match.TeamId2]}
	u.Publish(c, "predict", verb, object, tournament.Entity())

	return templateshlp.RenderJson(w, c, data)
}
