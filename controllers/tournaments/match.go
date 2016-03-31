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

package tournaments

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// MatchJSON is a variable to hold of match information.
//
type MatchJSON struct {
	ID         int64 `json:"Id"`
	IDNumber   int64 `json:"IdNumber"`
	Date       time.Time
	Team1      string
	Team2      string
	Iso1       string
	Iso2       string
	Location   string
	Result1    int64
	Result2    int64
	HasPredict bool
	Predict    string
	Finished   bool
	Ready      bool
	CanPredict bool
}

// Matches is the handler allowing to get the matches of a tournament.
// use the filter parameter to specify the matches you want:
// if filter is equal to 'first' you wil get matches of the first phase of the tournament.
// if filter is equal to 'second' you will get the matches of the second phase of the tournament.
//
func Matches(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Matches Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var t *mdl.Tournament

	if t, err = extract.Tournament(); err != nil {
		return err
	}

	filter := r.FormValue("filter")
	// if wrong data we set groupby to "first"
	if filter != "first" && filter != "second" {
		filter = "first"
	}

	var matchesJSON []MatchJSON

	if filter == "first" {
		matchesJSON = buildFirstPhaseMatches(c, t, u)
	} else if filter == "second" {
		matchesJSON = buildSecondPhaseMatches(c, t, u)
	}
	data := struct {
		Matches []MatchJSON
	}{
		matchesJSON,
	}

	return templateshlp.RenderJson(w, c, data)
}

// UpdateMatchResult is the handler allowing to update match of tournament with results information.
// from parameter 'result' with format 'result1 result2' the match information is updated accordingly.
//
func UpdateMatchResult(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Update Match Result Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var tournament *mdl.Tournament
	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	var match *mdl.Tmatch
	if match, err = extract.Match(tournament); err != nil {
		return err
	}

	result := r.FormValue("result")
	// is result well formated?
	results := strings.Split(result, " ")
	r1 := 0
	r2 := 0
	if len(results) != 2 {
		log.Errorf(c, "%s unable to get results, lenght not right: %v", desc, results)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}
	if r1, err = strconv.Atoi(results[0]); err != nil {
		log.Errorf(c, "%s unable to get results, error: %v not number 1", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}
	if r2, err = strconv.Atoi(results[1]); err != nil {
		log.Errorf(c, "%s unable to get results, error: %v not number 2", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	if err = mdl.SetResult(c, match, int64(r1), int64(r2), tournament); err != nil {
		log.Errorf(c, "%s unable to set result for match with id:%v error: %v", desc, match.IdNumber, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchCannotUpdate)}
	}

	// return the updated match
	var mjson MatchJSON
	mjson.IDNumber = match.IdNumber
	mjson.Date = match.Date
	rule := strings.Split(match.Rule, " ")

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(tournament); tb == nil {
		log.Errorf(c, "%s TournamentBuilder not found", desc)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	mapIDTeams := tb.MapOfIDTeams(c, tournament)

	if len(rule) > 1 {
		mjson.Team1 = rule[0]
		mjson.Team2 = rule[1]
	} else {
		mjson.Team1 = mapIDTeams[match.TeamId1]
		mjson.Team2 = mapIDTeams[match.TeamId2]
	}
	mjson.Location = match.Location

	mjson.Result1 = match.Result1
	mjson.Result2 = match.Result2

	// publish new activity
	object := mdl.ActivityEntity{Id: match.TeamId1, Type: "tteam", DisplayName: mapIDTeams[match.TeamId1]}
	target := mdl.ActivityEntity{Id: match.TeamId2, Type: "tteam", DisplayName: mapIDTeams[match.TeamId2]}
	verb := ""
	if match.Result1 > match.Result2 {
		verb = fmt.Sprintf("won %d-%d against", match.Result1, match.Result2)
	} else if match.Result1 < match.Result2 {
		verb = fmt.Sprintf("lost %d-%d against", match.Result1, match.Result2)
	} else {
		verb = fmt.Sprintf("tied %d-%d against", match.Result1, match.Result2)
	}
	tournament.Publish(c, "match", verb, object, target)

	return templateshlp.RenderJson(w, c, mjson)
}

// BlockMatchPrediction is the handler allowing to block the prediction for match of tournament.
//
func BlockMatchPrediction(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament block match prediction Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var tournament *mdl.Tournament
	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	var match *mdl.Tmatch
	if match, err = extract.Match(tournament); err != nil {
		return err
	}

	match.CanPredict = false
	if err := mdl.UpdateMatch(c, match); err != nil {
		log.Errorf(c, "%s unable to update match with id :%v", desc, match.IdNumber)
	}

	// return the updated match
	var mjson MatchJSON
	mjson.IDNumber = match.IdNumber
	mjson.Date = match.Date
	rule := strings.Split(match.Rule, " ")

	var tb mdl.TournamentBuilder

	if tb = mdl.GetTournamentBuilder(tournament); tb == nil {
		log.Errorf(c, "%s TournamentBuilder not found", desc)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	mapIDTeams := tb.MapOfIDTeams(c, tournament)

	if len(rule) > 1 {
		mjson.Team1 = rule[0]
		mjson.Team2 = rule[1]
	} else {
		mjson.Team1 = mapIDTeams[match.TeamId1]
		mjson.Team2 = mapIDTeams[match.TeamId2]
	}
	mjson.Location = match.Location

	mjson.Result1 = match.Result1
	mjson.Result2 = match.Result2

	return templateshlp.RenderJson(w, c, mjson)
}

// From a tournament entity return an array of MatchJSON data structure.
// second phase matches will have the specific rules in there team names
func buildMatchesFromTournament(c appengine.Context, t *mdl.Tournament, u *mdl.User) []MatchJSON {
	matchesJSON := buildFirstPhaseMatches(c, t, u)
	matches2ndPhase := buildSecondPhaseMatches(c, t, u)
	matchesJSON = append(matchesJSON, matches2ndPhase...)

	return matchesJSON
}

// From a tournament entity return an array of first phase MatchJSON data structure.
func buildFirstPhaseMatches(c appengine.Context, t *mdl.Tournament, u *mdl.User) []MatchJSON {
	desc := "buildMatchesFromTournament"

	matches := mdl.Matches(c, t.Matches1stStage)
	var predicts mdl.Predicts
	var err error
	if predicts, err = mdl.PredictsByIds(c, u.PredictIds); err != nil {
		log.Errorf(c, "%s predictions not found, %v", desc, err)
		return []MatchJSON{}
	}

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(t); tb == nil {
		log.Errorf(c, "%s tournament builder not found", desc)
		return []MatchJSON{}
	}

	mapIDTeams := tb.MapOfIDTeams(c, t)
	mapTeamCodes := tb.MapOfTeamCodes()

	matchesJSON := make([]MatchJSON, len(matches))
	for i, m := range matches {
		matchesJSON[i].ID = m.Id
		matchesJSON[i].IDNumber = m.IdNumber
		matchesJSON[i].Date = m.Date
		matchesJSON[i].Team1 = mapIDTeams[m.TeamId1]
		matchesJSON[i].Team2 = mapIDTeams[m.TeamId2]
		matchesJSON[i].Iso1 = mapTeamCodes[matchesJSON[i].Team1]
		matchesJSON[i].Iso2 = mapTeamCodes[matchesJSON[i].Team2]

		matchesJSON[i].Location = m.Location
		matchesJSON[i].Result1 = m.Result1
		matchesJSON[i].Result2 = m.Result2
		matchesJSON[i].Finished = m.Finished
		matchesJSON[i].Ready = m.Ready
		matchesJSON[i].CanPredict = m.CanPredict
		if hasMatch, j := predicts.ContainsMatchID(m.Id); hasMatch == true {
			matchesJSON[i].HasPredict = true
			matchesJSON[i].Predict = fmt.Sprintf("%v - %v", predicts[j].Result1, predicts[j].Result2)
		} else {
			matchesJSON[i].HasPredict = false
		}
	}

	return matchesJSON
}

// From a tournament entity return an array of second phase MatchJSON data structure.
// second phase matches will have the specific rules in there team names
func buildSecondPhaseMatches(c appengine.Context, t *mdl.Tournament, u *mdl.User) []MatchJSON {

	matches2ndPhase := mdl.Matches(c, t.Matches2ndStage)

	var predicts mdl.Predicts
	var err error
	if predicts, err = mdl.PredictsByIds(c, u.PredictIds); err != nil {
		return []MatchJSON{}
	}

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(t); tb == nil {
		return []MatchJSON{}
	}

	mapIDTeams := tb.MapOfIDTeams(c, t)
	mapTeamCodes := tb.MapOfTeamCodes()

	matchesJSON := make([]MatchJSON, len(matches2ndPhase))

	// append 2nd round to first one
	for i, m := range matches2ndPhase {
		matchesJSON[i].ID = m.Id
		matchesJSON[i].IDNumber = m.IdNumber
		matchesJSON[i].Date = m.Date
		rule := strings.Split(m.Rule, " ")
		if len(rule) == 2 {
			matchesJSON[i].Team1 = rule[0]
			matchesJSON[i].Team2 = rule[1]
			if _, ok := mapTeamCodes[rule[0]]; ok {
				matchesJSON[i].Iso1 = mapTeamCodes[rule[0]]
			}
			if _, ok := mapTeamCodes[rule[1]]; ok {
				matchesJSON[i].Iso2 = mapTeamCodes[rule[1]]
			}
		} else {
			if m.TeamId1 > 0 {
				matchesJSON[i].Team1 = mapIDTeams[m.TeamId1]
			} else {
				matchesJSON[i].Team1 = rule[0]
			}

			if m.TeamId2 > 0 {
				matchesJSON[i].Team2 = mapIDTeams[m.TeamId2]
			} else {
				matchesJSON[i].Team2 = rule[len(rule)-1]
			}

			matchesJSON[i].Iso1 = mapTeamCodes[mapIDTeams[m.TeamId1]]
			matchesJSON[i].Iso2 = mapTeamCodes[mapIDTeams[m.TeamId2]]

		}

		matchesJSON[i].Location = m.Location
		matchesJSON[i].Result1 = m.Result1
		matchesJSON[i].Result2 = m.Result2
		matchesJSON[i].Finished = m.Finished
		matchesJSON[i].Ready = m.Ready
		matchesJSON[i].CanPredict = m.CanPredict

		if hasMatch, j := predicts.ContainsMatchID(m.Id); hasMatch == true {
			matchesJSON[i].HasPredict = true
			matchesJSON[i].Predict = fmt.Sprintf("%v - %v", predicts[j].Result1, predicts[j].Result2)
		} else {
			matchesJSON[i].HasPredict = false
		}
	}
	return matchesJSON
}
