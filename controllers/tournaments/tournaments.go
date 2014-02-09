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
	"time"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	tournamentrelshlp "github.com/santiaago/purple-wing/helpers/tournamentrels"

	searchmdl "github.com/santiaago/purple-wing/models/search"
	teammdl "github.com/santiaago/purple-wing/models/team"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	tournamentinvmdl "github.com/santiaago/purple-wing/models/tournamentInvertedIndex"
	tournamentrelmdl "github.com/santiaago/purple-wing/models/tournamentrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

type TournamentData struct {
	Name string
}

// json index tournaments handler
func IndexJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournaments := tournamentmdl.FindAll(c)
		if len(tournaments) == 0 {
			return templateshlp.RenderEmptyJsonArray(w, c)
		}

		fieldsToKeep := []string{"Id", "Name"}
		tournamentsJson := make([]tournamentmdl.TournamentJson, len(tournaments))
		helpers.TransformFromArrayOfPointers(&tournaments, &tournamentsJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tournamentsJson)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json new tournament handler
func NewJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "Tournament New Handler: Error when decoding request body: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotCreate)}
		}

		var data TournamentData
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Errorf(c, "Tournament New Handler: Error when decoding request body: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotCreate)}
		}

		if len(data.Name) <= 0 {
			log.Errorf(c, "Tournamnet New Handler: 'Name' field cannot be empty")
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeNameCannotBeEmpty)}
		} else if t := tournamentmdl.Find(c, "KeyName", helpers.TrimLower(data.Name)); t != nil {
			log.Errorf(c, "Tournament New Handler: That tournament name already exists.")
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentAlreadyExists)}
		} else {
			tournament, err := tournamentmdl.Create(c, data.Name, "description foo", time.Now(), time.Now(), u.Id)
			if err != nil {
				log.Errorf(c, "Tournament New Handler: error when trying to create a tournament: %v", err)
				return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotCreate)}
			}
			// return the newly created tournament
			fieldsToKeep := []string{"Id", "Name"}
			var tJson tournamentmdl.TournamentJson
			helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

			return templateshlp.RenderJson(w, c, tJson)
		}
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// experimental: sar
// json new world cup tournament handler
func NewWorldCupJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		tournament, err := tournamentmdl.CreateWorldCup(c, u.Id)
		if err != nil {
			log.Errorf(c, "Tournament New World Cup Handler: error when trying to create a tournament: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotCreate)}
		}
		// return the newly created tournament
		// fieldsToKeep := []string{"Id", "Name"}
		// var tJson tournamentmdl.TournamentJson
		// helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tournament) //Json)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// Json show tournament handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {

		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournament Show Handler: error when extracting permalink id: %v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var tournament *tournamentmdl.Tournament
		tournament, err = tournamentmdl.ById(c, intID)
		if err != nil {
			log.Errorf(c, "Tournament Show Handler: tournament with id:%v was not found %v", intID, err)
			return helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		participants := tournamentrelshlp.Participants(c, intID)
		teams := tournamentrelshlp.Teams(c, intID)

		// tournament
		fieldsToKeep := []string{"Id", "Name"}
		var tournamentJson tournamentmdl.TournamentJson
		helpers.InitPointerStructure(tournament, &tournamentJson, fieldsToKeep)
		// participant
		participantFieldsToKeep := []string{"Id", "Username"}
		participantsJson := make([]usermdl.UserJson, len(participants))
		helpers.TransformFromArrayOfPointers(&participants, &participantsJson, participantFieldsToKeep)
		// teams
		teamsJson := make([]teammdl.TeamJson, len(teams))
		helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)
		// data
		data := struct {
			Tournament   tournamentmdl.TournamentJson
			Joined       bool
			Participants []usermdl.UserJson
			Teams        []teammdl.TeamJson
		}{
			tournamentJson,
			tournamentmdl.Joined(c, intID, u.Id),
			participantsJson,
			teamsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// Json tournament destroy handler
func DestroyJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {

		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournament Destroy Handler: error when extracting permalink id: %v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFoundCannotDelete)}
		}

		if !tournamentmdl.IsTournamentAdmin(c, intID, u.Id) {
			log.Errorf(c, "Tournament Destroy Handler: user is not admin")
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentDeleteForbiden)}
		}

		// delete all tournament-user relationships
		for _, participant := range tournamentrelshlp.Participants(c, intID) {
			if err := tournamentrelmdl.Destroy(c, intID, participant.Id); err != nil {
				log.Errorf(c, " error when trying to destroy tournament relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, team := range tournamentrelshlp.Teams(c, intID) {
			if err := tournamentteamrelmdl.Destroy(c, intID, team.Id); err != nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the tournament
		tournamentmdl.Destroy(c, intID)

		// return destroyed status
		return templateshlp.RenderJson(w, c, "tournament has been destroyed")
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

//  Json Update tournament handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Tournament Update Handler: error when extracting permalink id: %v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFoundCannotUpdate)}
		}

		if !tournamentmdl.IsTournamentAdmin(c, intID, u.Id) {
			log.Errorf(c, "Tournament Update Handler: user is not admin")
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentUpdateForbiden)}
		}

		var tournament *tournamentmdl.Tournament
		tournament, err = tournamentmdl.ById(c, intID)
		if err != nil {
			log.Errorf(c, "Tournament Update handler: tournament not found. id: %v, err: %v", intID, err)
			return helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFoundCannotUpdate)}
		}

		// only work on name other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf(c, "Tournament Update handler: Error when reading request body err: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
		}

		var updatedData TournamentData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			log.Errorf(c, "Tournament Update handler: Error when decoding request body err: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
		}

		if helpers.IsStringValid(updatedData.Name) && updatedData.Name != tournament.Name {
			// be sure that team with that name does not exist in datastore
			if t := tournamentmdl.Find(c, "KeyName", helpers.TrimLower(updatedData.Name)); t != nil {
				log.Errorf(c, "Tournament New Handler: That tournament name already exists.")
				return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentAlreadyExists)}
			}

			tournament.Name = updatedData.Name
			tournamentmdl.Update(c, intID, tournament)
		} else {
			log.Errorf(c, "Cannot update because updated data are not valid")
			log.Errorf(c, "Update name = %s", updatedData.Name)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotUpdate)}
		}

		// return the updated tournament
		fieldsToKeep := []string{"Id", "Name"}
		var tJson tournamentmdl.TournamentJson
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}

}

// json search tournaments handler
func SearchJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0) {

		words := helpers.SetOfStrings(keywords)
		ids, err := tournamentinvmdl.GetIndexes(c, words)
		if err != nil {
			log.Errorf(c, "Tournament Search Handler: tournaments.Index, error occurred when getting indexes of words: %v", err)
			return helpers.InternalServerError{errors.New(helpers.ErrorCodeTournamentCannotSearch)}
		}
		result := searchmdl.TournamentScore(c, keywords, ids)
		log.Infof(c, "result from TournamentScore: %v", result)
		tournaments := tournamentmdl.ByIds(c, result)
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
		tournamentsJson := make([]tournamentmdl.TournamentJson, len(tournaments))
		helpers.TransformFromArrayOfPointers(&tournaments, &tournamentsJson, fieldsToKeep)
		// we should not directly return an array. so we add an extra layer.
		data := struct {
			Tournaments []tournamentmdl.TournamentJson `json:",omitempty"`
		}{
			tournamentsJson,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json team candidates for a specific tournament:
func CandidateTeamsJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Candidate Teams Handler: error extracting permalink err:%v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		if _, err1 := tournamentmdl.ById(c, tournamentId); err1 != nil {
			log.Errorf(c, "Candidate Teams Handler: tournament not found err:%v", err)
			return helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		// query teams
		teams := teammdl.Find(c, "AdminId", u.Id)
		type canditateType struct {
			Team   teammdl.TeamJson
			Joined bool
		}
		fieldsToKeep := []string{"Id", "Name"}
		candidatesData := make([]canditateType, len(teams))

		for counterCandidate, team := range teams {
			var tJson teammdl.TeamJson
			helpers.InitPointerStructure(team, &tJson, fieldsToKeep)
			var canditate canditateType
			canditate.Team = tJson
			canditate.Joined = tournamentmdl.TeamJoined(c, tournamentId, team.Id)
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
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json tournament participants handler
// use this handler to get participants of a tournament.
func ParticipantsJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "Tournament Participants Handler: error extracting permalink err:%v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		participants := tournamentrelshlp.Participants(c, tournamentId)
		// participant
		participantFieldsToKeep := []string{"Id", "Username"}
		participantsJson := make([]usermdl.UserJson, len(participants))
		helpers.TransformFromArrayOfPointers(&participants, &participantsJson, participantFieldsToKeep)
		// data
		data := struct {
			Participants []usermdl.UserJson
		}{
			participantsJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json tournament groups handler
// use this handler to get groups of a tournament.
func GroupsJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "Tournament Groups Handler: error extracting permalink err:%v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *tournamentmdl.Tournament
		tournament, err = tournamentmdl.ById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "Tournament Group Handler: tournament with id:%v was not found %v", tournamentId, err)
			return helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		groups := tournamentmdl.Groups(c, tournament.GroupIds)
		// ToDo: might need to filter information here.
		data := struct {
			Groups []*tournamentmdl.Tgroup
		}{
			groups,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}
