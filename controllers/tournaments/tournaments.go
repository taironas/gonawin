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
	"sort"
	"strings"
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

type MatchJson struct {
	IdNumber int64
	Date     time.Time
	Team1    string
	Team2    string
	Location string
}

type DayJson struct {
	Date    time.Time
	Matches []MatchJson
}

type PhaseJson struct {
	Name string
	Days []DayJson
}

// json tournament calendar handler
// use this handler to get calendar of a tournament.
// the calendar structure is an array of matches of the tournament
// with the location, the teams involved and the date
func CalendarJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
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

		matches := tournamentmdl.Matches(c, tournament.Matches1stStage)

		mapIdTeams := tournamentmdl.MapOfIdTeams(c, *tournament)

		// ToDo: might need to filter information here.

		matchesJson := make([]MatchJson, len(matches))
		for i, m := range matches {
			matchesJson[i].IdNumber = m.IdNumber
			matchesJson[i].Date = m.Date
			matchesJson[i].Team1 = mapIdTeams[m.TeamId1]
			matchesJson[i].Team2 = mapIdTeams[m.TeamId2]
			matchesJson[i].Location = m.Location
		}

		data := struct {
			Matches []MatchJson
		}{
			matchesJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json tournament calendar by day handler
// use this handler to get calendar grouped by day of a tournament.
// the calendar structure is an array of days that contains an array of matches for a tournament
// with the location, the teams involved and the date
func CalendarByDayJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
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

		matches := tournamentmdl.Matches(c, tournament.Matches1stStage)
		matches2ndPhase := tournamentmdl.Matches(c, tournament.Matches2ndStage)

		mapIdTeams := tournamentmdl.MapOfIdTeams(c, *tournament)

		// ToDo: might need to filter information here.
		matchesJson := make([]MatchJson, len(matches))
		for i, m := range matches {
			matchesJson[i].IdNumber = m.IdNumber
			matchesJson[i].Date = m.Date
			matchesJson[i].Team1 = mapIdTeams[m.TeamId1]
			matchesJson[i].Team2 = mapIdTeams[m.TeamId2]
			matchesJson[i].Location = m.Location
		}

		// append 2nd round to first one
		for _, m := range matches2ndPhase {
			var matchJson2ndPhase MatchJson
			matchJson2ndPhase.IdNumber = m.IdNumber
			matchJson2ndPhase.Date = m.Date
			rule := strings.Split(m.Rule, " ")
			matchJson2ndPhase.Team1 = rule[0]
			matchJson2ndPhase.Team2 = rule[1]
			matchJson2ndPhase.Location = m.Location
			// append second round results
			matchesJson = append(matchesJson, matchJson2ndPhase)

		}

		// get array of days from matches
		log.Infof(c, "Tournament Calendar By Day Handler: ready to build days array")

		var days []DayJson
		fillDaysFromMatches(c, &days, matchesJson)
		log.Infof(c, "Tournament Calendar By Day Handler: days array ok")

		data := struct {
			Days []DayJson
		}{
			days,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// from an array of matches build an array of dates where matches occur.
// each element in this array is itself an array of matches that occur that day
func fillDaysFromMatches(c appengine.Context, days *[]DayJson, matches []MatchJson) {

	log.Infof(c, "fill days from matches:  start")

	mapOfDays := make(map[string][]MatchJson)

	log.Infof(c, "fill days from matches:  make map ok")
	const shortForm = "Jan/02/2006"
	for _, m := range matches {
		currentDate := m.Date.Format(shortForm)
		_, ok := mapOfDays[currentDate]
		if ok {
			mapOfDays[currentDate] = append(mapOfDays[currentDate], m)
		} else {
			var arrayMatches []MatchJson
			arrayMatches = append(arrayMatches, m)
			mapOfDays[currentDate] = arrayMatches
		}
	}
	log.Infof(c, "fill days from matches:  map built")

	i := 0
	//arrayDays := *days
	*days = make([]DayJson, len(mapOfDays))
	for key, value := range mapOfDays {
		log.Infof(c, "fill days from matches:  i: %v", i)
		log.Infof(c, "fill days from matches:  key: %v", key)
		log.Infof(c, "fill days from matches:  value: %v", value)
		(*days)[i].Date, _ = time.Parse(shortForm, key)
		(*days)[i].Matches = value
		log.Infof(c, "fill days from matches:  day struct: %v", (*days)[i])
		i++
	}
	log.Infof(c, "fill days from matches:  array of days ready")

	sort.Sort(ByDate(*days))
	log.Infof(c, "fill days from matches:  days are now sorted")

	log.Infof(c, "fill days from matches:  end")
}

func fillDaysFromMatchesInInterval(c appengine.Context, days *[]DayJson, matches []MatchJson, low int64, high int64) {
	var filteredMatches []MatchJson
	for _, v := range matches {
		if v.IdNumber >= low && v.IdNumber <= high {
			filteredMatches = append(filteredMatches, v)
		}
	}
	fillDaysFromMatches(c, days, filteredMatches)
}

// ByDate implements sort.Interface for []Person based on the date field.
type ByDate []DayJson

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// json tournament calendar by phase handler
// use this handler to get calendar grouped by phase of a tournament.
// the calendar structure is an array of phases that contains an array of matches for a tournament
// with the location, the teams involved and the date
func CalendarByPhaseJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "Tournament Calendar by Phase Handler: error extracting permalink err:%v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var tournament *tournamentmdl.Tournament
		tournament, err = tournamentmdl.ById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "Tournament Calendar by Phase Handler: tournament with id:%v was not found %v", tournamentId, err)
			return helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		matches := tournamentmdl.Matches(c, tournament.Matches1stStage)
		matches2ndPhase := tournamentmdl.Matches(c, tournament.Matches2ndStage)

		mapIdTeams := tournamentmdl.MapOfIdTeams(c, *tournament)

		// ToDo: might need to filter information here.
		matchesJson := make([]MatchJson, len(matches))
		for i, m := range matches {
			matchesJson[i].IdNumber = m.IdNumber
			matchesJson[i].Date = m.Date
			matchesJson[i].Team1 = mapIdTeams[m.TeamId1]
			matchesJson[i].Team2 = mapIdTeams[m.TeamId2]
			matchesJson[i].Location = m.Location
		}

		// append 2nd round to first one
		for _, m := range matches2ndPhase {
			var matchJson2ndPhase MatchJson
			matchJson2ndPhase.IdNumber = m.IdNumber
			matchJson2ndPhase.Date = m.Date
			rule := strings.Split(m.Rule, " ")
			matchJson2ndPhase.Team1 = rule[0]
			matchJson2ndPhase.Team2 = rule[1]
			matchJson2ndPhase.Location = m.Location
			// append second round results
			matchesJson = append(matchesJson, matchJson2ndPhase)

		}

		// get array of days from matches
		log.Infof(c, "Tournament Calendar by Phase Handler: ready to build phase array")
		phaseNames := []string{"First Stage", "Round of 16", "Quarter-finals", "Semi-finals", "Third Place", "Finals"}
		limits := make(map[string][]int64)
		limits["First Stage"] = []int64{1, 48}
		limits["Round of 16"] = []int64{49, 56}
		limits["Quarter-finals"] = []int64{57, 60}
		limits["Semi-finals"] = []int64{61, 62}
		limits["Third Place"] = []int64{63, 63}
		limits["Finals"] = []int64{64, 64}

		// first stage: matches 1 to 48
		// Round of 16: matches 49 to 56
		// Quarte-finals: matches 57 to 60
		// Semi-finals: matches 61 to 62
		// Thrid Place: match 63
		// Finals: match 64
		phases := make([]PhaseJson, len(phaseNames))
		for i, _ := range phases {
			phases[i].Name = phaseNames[i]
			log.Infof(c, "Tournament Calendar by Phase Handler: name %s", phases[i].Name)
			low := limits[phases[i].Name][0]
			high := limits[phases[i].Name][1]
			// build days array
			var days []DayJson
			fillDaysFromMatchesInInterval(c, &days, matchesJson, low, high)
			phases[i].Days = days
		}

		data := struct {
			Phases []PhaseJson
		}{
			phases,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}
