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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	mdl "github.com/santiaago/purple-wing/models"
	activitymdl "github.com/santiaago/purple-wing/models/activity"
	predictmdl "github.com/santiaago/purple-wing/models/predict"
	tournamentinvmdl "github.com/santiaago/purple-wing/models/tournamentInvertedIndex"
)

type TournamentData struct {
	Name string
}

// json index tournaments handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournaments := mdl.FindAllTournaments(c)
		if len(tournaments) == 0 {
			return templateshlp.RenderEmptyJsonArray(w, c)
		}

		fieldsToKeep := []string{"Id", "Name"}
		tournamentsJson := make([]mdl.TournamentJson, len(tournaments))
		helpers.TransformFromArrayOfPointers(&tournaments, &tournamentsJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tournamentsJson)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json new tournament handler
func NewJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "Tournament New Handler: Error when decoding request body: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotCreate)}
		}

		var data TournamentData
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Errorf(c, "Tournament New Handler: Error when decoding request body: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotCreate)}
		}

		if len(data.Name) <= 0 {
			log.Errorf(c, "Tournamnet New Handler: 'Name' field cannot be empty")
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeNameCannotBeEmpty)}
		} else if t := mdl.FindTournaments(c, "KeyName", helpers.TrimLower(data.Name)); t != nil {
			log.Errorf(c, "Tournament New Handler: That tournament name already exists.")
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentAlreadyExists)}
		} else {
			tournament, err := mdl.CreateTournament(c, data.Name, "description foo", time.Now(), time.Now(), u.Id)
			if err != nil {
				log.Errorf(c, "Tournament New Handler: error when trying to create a tournament: %v", err)
				return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotCreate)}
			}
			// return the newly created tournament
			fieldsToKeep := []string{"Id", "Name"}
			var tJson mdl.TournamentJson
			helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

			// publish new activity
			actor := activitymdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
			object := activitymdl.ActivityEntity{ID: tournament.Id, Type: "tournament", DisplayName: tournament.Name}
			target := activitymdl.ActivityEntity{}
			activitymdl.Publish(c, "tournament", "created a tournament", actor, object, target, u.Id)

			return templateshlp.RenderJson(w, c, tJson)
		}
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// Json show tournament handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {

		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournament Show Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var tournament *mdl.Tournament
		tournament, err = mdl.TournamentById(c, intID)
		if err != nil {
			log.Errorf(c, "Tournament Show Handler: tournament with id:%v was not found %v", intID, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		participants := tournament.Participants(c)
		teams := tournament.Teams(c)

		// tournament
		fieldsToKeep := []string{"Id", "Name"}
		var tournamentJson mdl.TournamentJson
		helpers.InitPointerStructure(tournament, &tournamentJson, fieldsToKeep)
		// participant
		participantFieldsToKeep := []string{"Id", "Username"}
		participantsJson := make([]mdl.UserJson, len(participants))
		helpers.TransformFromArrayOfPointers(&participants, &participantsJson, participantFieldsToKeep)
		// teams
		teamsJson := make([]mdl.TeamJson, len(teams))
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)
		// data
		data := struct {
			Tournament   mdl.TournamentJson
			Joined       bool
			Participants []mdl.UserJson
			Teams        []mdl.TeamJson
		}{
			tournamentJson,
			tournament.Joined(c, u),
			participantsJson,
			teamsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// Json tournament destroy handler
func DestroyJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {

		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournament Destroy Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFoundCannotDelete)}
		}

		if !mdl.IsTournamentAdmin(c, intID, u.Id) {
			log.Errorf(c, "Tournament Destroy Handler: user is not admin")
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentDeleteForbiden)}
		}
		var tournament *mdl.Tournament
		tournament, err = mdl.TournamentById(c, intID)
		if err != nil {
			log.Errorf(c, "Tournament Show Handler: tournament with id:%v was not found %v", intID, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		// delete all tournament-user relationships
		for _, participant := range tournament.Participants(c) {
			participant.RemoveTournamentId(c, tournament.Id)
			if err := participant.RemoveTournamentId(c, tournament.Id); err != nil {
				log.Errorf(c, " error when trying to remove tournament id from user: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, team := range tournament.Teams(c) {
			if err := tournament.TeamLeave(c, team); err != nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the tournament
		tournament.Destroy(c)

		// return destroyed status
		return templateshlp.RenderJson(w, c, "tournament has been destroyed")
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

//  Json Update tournament handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournament Update Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFoundCannotUpdate)}
		}

		if !mdl.IsTournamentAdmin(c, intID, u.Id) {
			log.Errorf(c, "Tournament Update Handler: user is not admin")
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentUpdateForbiden)}
		}

		var tournament *mdl.Tournament
		tournament, err = mdl.TournamentById(c, intID)
		if err != nil {
			log.Errorf(c, "Tournament Update handler: tournament not found. id: %v, err: %v", intID, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFoundCannotUpdate)}
		}

		// only work on name other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "Tournament Update handler: Error when reading request body err: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
		}

		var updatedData TournamentData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			log.Errorf(c, "Tournament Update handler: Error when decoding request body err: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
		}

		if helpers.IsStringValid(updatedData.Name) && updatedData.Name != tournament.Name {
			// be sure that team with that name does not exist in datastore
			if t := mdl.FindTournaments(c, "KeyName", helpers.TrimLower(updatedData.Name)); t != nil {
				log.Errorf(c, "Tournament New Handler: That tournament name already exists.")
				return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentAlreadyExists)}
			}

			tournament.Name = updatedData.Name
			tournament.Update(c)
		} else {
			log.Errorf(c, "Cannot update because updated data are not valid")
			log.Errorf(c, "Update name = %s", updatedData.Name)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
		}

		// return the updated tournament
		fieldsToKeep := []string{"Id", "Name"}
		var tJson mdl.TournamentJson
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}

}

// json search tournaments handler
func SearchJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0) {

		words := helpers.SetOfStrings(keywords)
		ids, err := tournamentinvmdl.GetIndexes(c, words)
		if err != nil {
			log.Errorf(c, "Tournament Search Handler: tournaments.Index, error occurred when getting indexes of words: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotSearch)}
		}
		result := mdl.TournamentScore(c, keywords, ids)
		log.Infof(c, "result from TournamentScore: %v", result)
		tournaments := mdl.TournamentsByIds(c, result)
		log.Infof(c, "ByIds result %v", tournaments)
		if len(tournaments) == 0 {
			msg := fmt.Sprintf("Oops! Your search - %s - did not match any %s.", keywords, "tournament")
			data := struct {
				MessageInfo string `json:",omitempty"`
			}{
				msg,
			}
			return templateshlp.RenderJson(w, c, data)
		}

		fieldsToKeep := []string{"Id", "Name"}
		tournamentsJson := make([]mdl.TournamentJson, len(tournaments))
		helpers.TransformFromArrayOfPointers(&tournaments, &tournamentsJson, fieldsToKeep)
		// we should not directly return an array. so we add an extra layer.
		data := struct {
			Tournaments []mdl.TournamentJson `json:",omitempty"`
		}{
			tournamentsJson,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json team candidates for a specific tournament:
func CandidateTeamsJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Candidate Teams Handler: error extracting permalink err:%v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *mdl.Tournament
		tournament, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "Candidate Teams Handler: tournament not found err:%v", err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		// query teams
		teams := mdl.FindTeams(c, "AdminId", u.Id)
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
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json tournament participants handler
// use this handler to get participants of a tournament.
func ParticipantsJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "Tournament Participants Handler: error extracting permalink err:%v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *mdl.Tournament
		tournament, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "Tournament Show Handler: tournament with id:%v was not found %v", tournamentId, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		participants := tournament.Participants(c)
		// participant
		participantFieldsToKeep := []string{"Id", "Username"}
		participantsJson := make([]mdl.UserJson, len(participants))
		helpers.TransformFromArrayOfPointers(&participants, &participantsJson, participantFieldsToKeep)
		// data
		data := struct {
			Participants []mdl.UserJson
		}{
			participantsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// Reset a tournament information. Reset points and goals.
func ResetJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)

		if err != nil {
			log.Errorf(c, "Tournament Reset Handler: error extracting permalink err:%v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var t *mdl.Tournament
		t, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "Tournament Update Match Result Handler: tournament with id:%v was not found %v", tournamentId, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		if err = mdl.Reset(c, t); err != nil {
			log.Errorf(c, "Tournament Reset Handler: Unable to reset tournament: %v error:", tournamentId, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeInternal)}
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
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// Set a Predict entity of a specific match for the current User.
func PredictJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Predict Handler"

	if r.Method == "POST" {
		// extract tournament
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s: error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *mdl.Tournament
		tournament, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "%s: tournament with id:%v was not found %v", desc, tournamentId, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		// extract match
		matchIdNumber, err2 := handlers.PermalinkID(r, c, 5)
		if err2 != nil {
			log.Errorf(c, "%s: error extracting permalink err:%v", desc, err2)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeMatchNotFoundCannotSetPrediction)}
		}

		match := mdl.GetMatchByIdNumber(c, *tournament, matchIdNumber)
		if match == nil {
			log.Errorf(c, "%s: unable to get match with id number :%v", desc, matchIdNumber)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeMatchNotFoundCannotSetPrediction)}
		}
		result1 := r.FormValue("result1")
		result2 := r.FormValue("result2")
		var r1, r2 int
		if r1, err = strconv.Atoi(result1); err != nil {
			log.Errorf(c, "%s: unable to get results, error: %v not number 1", desc, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeCannotSetPrediction)}
		}
		if r2, err = strconv.Atoi(result2); err != nil {
			log.Errorf(c, "%s: unable to get results, error: %v not number 2", desc, err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeCannotSetPrediction)}
		}

		if p, err1 := predictmdl.Create(c, int64(r1), int64(r2), match.Id); err1 != nil {
			log.Errorf(c, "%s: unable to create Predict for match with id:%v error: %v", desc, match.Id, err1)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeCannotSetPrediction)}
		} else {
			// add p.Id to User predict table.
			if err = u.AddPredictId(c, p.Id); err != nil {
				log.Errorf(c, "%s: unable to add predict id in user entity: error: %v", desc, err)
				return &helpers.NotFound{errors.New(helpers.ErrorCodeCannotSetPrediction)}
			}
			msg := fmt.Sprintf("Prediction is now set.")
			data := struct {
				MessageInfo string `json:",omitempty"`
				Predict     *predictmdl.Predict
			}{
				msg,
				p,
			}
			return templateshlp.RenderJson(w, c, data)
		}
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}
