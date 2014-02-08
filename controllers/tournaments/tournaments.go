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

	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json new tournament handler
func NewJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when reading request body")}
		}

		var data TournamentData
		err = json.Unmarshal(body, &data)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when decoding request body")}
		}

		if len(data.Name) <= 0 {
			return helpers.InternalServerError{errors.New("'Name' field cannot be empty")}
		} else if t := tournamentmdl.Find(c, "KeyName", helpers.TrimLower(data.Name)); t != nil {
			return helpers.InternalServerError{errors.New("That tournament name already exists")}
		} else {
			tournament, err := tournamentmdl.Create(c, data.Name, "description foo", time.Now(), time.Now(), u.Id)
			if err != nil {
				log.Errorf(c, " error when trying to create a tournament: %v", err)
				return helpers.InternalServerError{errors.New("error when trying to create a tournament")}
			}
			// return the newly created tournament
			fieldsToKeep := []string{"Id", "Name"}
			var tJson tournamentmdl.TournamentJson
			helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

			return templateshlp.RenderJson(w, c, tJson)
		}
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// Json show tournament handler
func ShowJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	intID, err := handlers.PermalinkID(r, c, 4)
	if err != nil {
		return helpers.NotFound{err}
	}

	if r.Method == "GET" {
		var tournament *tournamentmdl.Tournament
		tournament, err = tournamentmdl.ById(c, intID)
		if err != nil {
			return helpers.NotFound{err}
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
	} else {
		return helpers.BadRequest{errors.New("Not supported.")}
	}
}

// Json tournament destroy handler
func DestroyJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	intID, err := handlers.PermalinkID(r, c, 4)
	if err != nil {
		return helpers.NotFound{err}
	}

	if r.Method == "POST" {
		if !tournamentmdl.IsTournamentAdmin(c, intID, u.Id) {
			return helpers.Forbidden{errors.New("tournament can only be deleted by the tournament administrator")}
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
	} else {
		return helpers.BadRequest{errors.New("Not supported.")}
	}
}

//  Json Update tournament handler
func UpdateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		intID, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			return helpers.NotFound{err}
		}

		if !tournamentmdl.IsTournamentAdmin(c, intID, u.Id) {
			return helpers.Forbidden{errors.New("tournament can only be updated by the tournament administrator")}
		}

		var tournament *tournamentmdl.Tournament
		tournament, err = tournamentmdl.ById(c, intID)
		if err != nil {
			log.Errorf(c, " Tournament Edit handler: tournament not found. id: %v", intID)
			return helpers.NotFound{err}
		}

		// only work on name other values should not be editable
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when reading request body")}
		}

		var updatedData TournamentData
		err = json.Unmarshal(body, &updatedData)
		if err != nil {
			return helpers.InternalServerError{errors.New("Error when decoding request body")}
		}

		if helpers.IsStringValid(updatedData.Name) && updatedData.Name != tournament.Name {
			tournament.Name = updatedData.Name
			tournamentmdl.Update(c, intID, tournament)
		} else {
			log.Errorf(c, "Cannot update because updated data are not valid")
			log.Errorf(c, "Update name = %s", updatedData.Name)
		}
		// return the updated tournament
		fieldsToKeep := []string{"Id", "Name"}
		var tJson tournamentmdl.TournamentJson
		helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

		return templateshlp.RenderJson(w, c, tJson)
	} else {
		return helpers.BadRequest{errors.New("Not supported.")}
	}
}

// json search tournaments handler
func SearchJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)
	log.Infof(c, "json search handler.")
	keywords := r.FormValue("q")
	if r.Method == "GET" && (len(keywords) > 0) {
		words := helpers.SetOfStrings(keywords)
		ids, err := tournamentinvmdl.GetIndexes(c, words)
		if err != nil {
			log.Errorf(c, " tournaments.Index, error occurred when getting indexes of words: %v", err)
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
	} else {
		return helpers.BadRequest{errors.New("not supported.")}
	}
}

// json team candidates for a specific tournament:
func CandidateTeamsJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)
	log.Infof(c, "json team candidates handler.")

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			return helpers.NotFound{err}
		}
		if _, err1 := tournamentmdl.ById(c, tournamentId); err1 != nil {
			return helpers.NotFound{err}
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
	} else {
		return helpers.BadRequest{errors.New("Not supported.")}
	}
}

// json tournament participants handler
// use this handler to get participants of a tournament.
func ParticipantsJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)
	log.Infof(c, "json tournament participants handler.")

	tournamentId, err := handlers.PermalinkID(r, c, 3)
	if err != nil {
		return helpers.NotFound{err}
	}

	if r.Method == "GET" {
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
	} else {
		return helpers.BadRequest{errors.New("Not supported.")}
	}
}
