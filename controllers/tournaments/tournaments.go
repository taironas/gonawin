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
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/auth"
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

type Form struct {
	Name  string
	Error string
}

type TournamentData struct {
	Name string
}

type indexData struct {
	Tournaments           []*tournamentmdl.Tournament
	TournamentInputSearch string
}

// index tournaments handler
func Index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var data indexData
	if r.Method == "GET" {
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
		"Tournaments": func() bool { return true },
	}

	templateshlp.RenderWithData(w, r, c, t, data, funcs, "renderTournamentIndex")
}

// new tournament handler
func New(w http.ResponseWriter, r *http.Request) {
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
			tournament, err := tournamentmdl.Create(c, form.Name, "description foo", time.Now(), time.Now(), auth.CurrentUser(r, c).Id)
			if err != nil {
				log.Errorf(c, " error when trying to create a tournament: %v", err)
			}
			// redirect to the newly created tournament page
			http.Redirect(w, r, "/m/tournaments/"+fmt.Sprintf("%d", tournament.Id), http.StatusFound)
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

// show tournament handler
func Show(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil {
		log.Errorf(c, " Unable to find ID in request %v", r)
		http.Redirect(w, r, "/m/tournaments/", http.StatusFound)
		return
	}

	if r.Method == "GET" {
		funcs := template.FuncMap{
			"Joined":            func() bool { return tournamentmdl.Joined(c, intID, auth.CurrentUser(r, c).Id) },
			"TeamJoined":        func(teamId int64) bool { return tournamentmdl.TeamJoined(c, intID, teamId) },
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

		if err != nil {
			helpers.Error404(w)
			return
		}

		participants := tournamentrelshlp.Participants(c, intID)
		teams := tournamentrelshlp.Teams(c, intID)
		candidateTeams := teammdl.Find(c, "AdminId", auth.CurrentUser(r, c).Id)

		tournamentData := struct {
			Tournament     *tournamentmdl.Tournament
			Participants   []*usermdl.User
			Teams          []*teammdl.Team
			CandidateTeams []*teammdl.Team
		}{
			tournament,
			participants,
			teams,
			candidateTeams,
		}
		templateshlp.RenderWithData(w, r, c, t, tournamentData, funcs, "renderTournamentShow")
	}
}

//  Edit tournament handler
func Edit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	intID, err := handlers.PermalinkID(r, c, 3)
	if err != nil {
		http.Redirect(w, r, "/m/tournaments/", http.StatusFound)
		return
	}

	if !tournamentmdl.IsTournamentAdmin(c, intID, auth.CurrentUser(r, c).Id) {
		http.Redirect(w, r, "/m", http.StatusFound)
		return
	}

	var tournament *tournamentmdl.Tournament
	tournament, err = tournamentmdl.ById(c, intID)
	if err != nil {
		log.Errorf(c, " Tournament Edit handler: tournament not found. id: %v", intID)
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

		if helpers.IsStringValid(editName) && editName != tournament.Name {
			tournament.Name = editName
			tournamentmdl.Update(c, intID, tournament)
		} else {
			log.Errorf(c, " cannot update %v", helpers.IsStringValid(editName))
		}
		url := fmt.Sprintf("/m/tournaments/%d", intID)
		http.Redirect(w, r, url, http.StatusFound)
		return
	} else {
		helpers.Error404(w)
		return
	}
}

// Tournament destroy handler
func Destroy(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	tournamentID, err := handlers.PermalinkID(r, c, 4)
	if err != nil {
		http.Redirect(w, r, "/m/tournaments/", http.StatusFound)
		return
	}

	if r.Method == "POST" {
		// delete all tournament-user relationships
		for _, participant := range tournamentrelshlp.Participants(c, tournamentID) {
			if err := tournamentrelmdl.Destroy(c, tournamentID, participant.Id); err != nil {
				log.Errorf(c, " error when trying to destroy tournament relationship: %v", err)
			}
		}
		// delete all tournament-team relationships
		for _, team := range tournamentrelshlp.Teams(c, tournamentID) {
			if err := tournamentteamrelmdl.Destroy(c, tournamentID, team.Id); err != nil {
				log.Errorf(c, " error when trying to destroy team relationship: %v", err)
			}
		}
		// delete the tournament
		tournamentmdl.Destroy(c, tournamentID)

		http.Redirect(w, r, "/m/tournaments", http.StatusFound)
		return
	}
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
		counterTournament := 0
		for _, tournament := range tournaments {
			var tJson tournamentmdl.TournamentJson
			helpers.CopyToPointerStructure(tournament, &tJson)
			helpers.KeepFields(&tJson, fieldsToKeep)
			tournamentsJson[counterTournament] = tJson
			counterTournament++
		}

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
			helpers.CopyToPointerStructure(tournament, &tJson)
			helpers.KeepFields(&tJson, fieldsToKeep)

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
		helpers.CopyToPointerStructure(tournament, &tournamentJson)
		helpers.KeepFields(&tournamentJson, fieldsToKeep)
		// participant
		participantsJson := make([]usermdl.UserJson, len(participants))
		counterParticipant := 0
		participantFieldsToKeep := []string{"Id", "Username"}
		for _, participant := range participants {
			var pJson usermdl.UserJson
			helpers.CopyToPointerStructure(participant, &pJson)
			helpers.KeepFields(&pJson, participantFieldsToKeep)
			participantsJson[counterParticipant] = pJson
			counterParticipant++
		}
		// teams
		teamsJson := make([]teammdl.TeamJson, len(teams))
		counterTeam := 0
		for _, team := range teams {
			var teamJson teammdl.TeamJson
			helpers.CopyToPointerStructure(team, &teamJson)
			helpers.KeepFields(&teamJson, fieldsToKeep)
			teamsJson[counterTeam] = teamJson
			counterTeam++
		}
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
		helpers.CopyToPointerStructure(tournament, &tJson)
		helpers.KeepFields(&tJson, fieldsToKeep)

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
			return templateshlp.RenderEmptyJsonArray(w, c)
		}

		fieldsToKeep := []string{"Id", "Name"}
		tournamentsJson := make([]tournamentmdl.TournamentJson, len(tournaments))
		counterTournament := 0
		for _, tournament := range tournaments {
			var tJson tournamentmdl.TournamentJson
			helpers.CopyToPointerStructure(tournament, &tJson)
			helpers.KeepFields(&tJson, fieldsToKeep)
			tournamentsJson[counterTournament] = tJson
			counterTournament++
		}

		return templateshlp.RenderJson(w, c, tournamentsJson)
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
		counterCandidate := 0
		for _, team := range teams {
			var tJson teammdl.TeamJson
			helpers.CopyToPointerStructure(team, &tJson)
			helpers.KeepFields(&tJson, fieldsToKeep)
			var canditate canditateType
			canditate.Team = tJson
			canditate.Joined = tournamentmdl.TeamJoined(c, tournamentId, team.Id)
			candidatesData[counterCandidate] = canditate
			counterCandidate++
		}
		return templateshlp.RenderJson(w, c, candidatesData)
	} else {
		return helpers.BadRequest{errors.New("Not supported.")}
	}
}
