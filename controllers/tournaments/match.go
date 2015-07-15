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

// A MatchJson is a variable to hold of match information.
type MatchJson struct {
	Id         int64
	IdNumber   int64
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

// Json tournament Matches handler
// use this handler to get the matches of a tournament.
// use the filter parameter to specify the matches you want:
// if filter is equal to 'first' you wil get matches of the first phase of the tournament.
// if filter is equal to 'second' you will get the matches of the second phase of the tournament.
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

	log.Infof(c, "%s ready to build days array", desc)
	var matchesJson []MatchJson

	if filter == "first" {
		matchesJson = buildFirstPhaseMatches(c, t, u)
	} else if filter == "second" {
		matchesJson = buildSecondPhaseMatches(c, t, u)
	}
	data := struct {
		Matches []MatchJson
	}{
		matchesJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// Update Match handler.
// Update match of tournament with results information.
// from parameter 'result' with format 'result1 result2' the match information is updated accordingly.
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
	var mjson MatchJson
	mjson.IdNumber = match.IdNumber
	mjson.Date = match.Date
	rule := strings.Split(match.Rule, " ")

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(tournament); tb == nil {
		log.Errorf(c, "%s TournamentBuilder not found", desc)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	mapIdTeams := tb.MapOfIdTeams(c, tournament)

	if len(rule) > 1 {
		mjson.Team1 = rule[0]
		mjson.Team2 = rule[1]
	} else {
		mjson.Team1 = mapIdTeams[match.TeamId1]
		mjson.Team2 = mapIdTeams[match.TeamId2]
	}
	mjson.Location = match.Location

	mjson.Result1 = match.Result1
	mjson.Result2 = match.Result2

	// publish new activity
	object := mdl.ActivityEntity{Id: match.TeamId1, Type: "tteam", DisplayName: mapIdTeams[match.TeamId1]}
	target := mdl.ActivityEntity{Id: match.TeamId2, Type: "tteam", DisplayName: mapIdTeams[match.TeamId2]}
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

// Block Match prediction handler.
// Block prediction for match of tournament.
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
	var mjson MatchJson
	mjson.IdNumber = match.IdNumber
	mjson.Date = match.Date
	rule := strings.Split(match.Rule, " ")

	var tb mdl.TournamentBuilder

	if tb = mdl.GetTournamentBuilder(tournament); tb == nil {
		log.Errorf(c, "%s TournamentBuilder not found", desc)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	mapIdTeams := tb.MapOfIdTeams(c, tournament)

	if len(rule) > 1 {
		mjson.Team1 = rule[0]
		mjson.Team2 = rule[1]
	} else {
		mjson.Team1 = mapIdTeams[match.TeamId1]
		mjson.Team2 = mapIdTeams[match.TeamId2]
	}
	mjson.Location = match.Location

	mjson.Result1 = match.Result1
	mjson.Result2 = match.Result2

	return templateshlp.RenderJson(w, c, mjson)
}

// From a tournament entity return an array of MatchJson data structure.
// second phase matches will have the specific rules in there team names
func buildMatchesFromTournament(c appengine.Context, t *mdl.Tournament, u *mdl.User) []MatchJson {
	desc := "buildMatchesFromTournament"
	log.Infof(c, "%s start", desc)
	matchesJson := buildFirstPhaseMatches(c, t, u)
	log.Infof(c, "%s done with buildFirstPhaseMatches", desc)
	matches2ndPhase := buildSecondPhaseMatches(c, t, u)
	log.Infof(c, "%s done with buildSecondPhaseMatches", desc)
	matchesJson = append(matchesJson, matches2ndPhase...)

	return matchesJson
}

// From a tournament entity return an array of first phase MatchJson data structure.
func buildFirstPhaseMatches(c appengine.Context, t *mdl.Tournament, u *mdl.User) []MatchJson {
	desc := "buildMatchesFromTournament"
	log.Infof(c, "%s start", desc)
	log.Infof(c, "%s %v", desc, t.Matches1stStage)

	matches := mdl.Matches(c, t.Matches1stStage)
	log.Infof(c, "%s number of matches retrieved %v", desc, len(matches))
	var predicts mdl.Predicts
	var err error
	if predicts, err = mdl.PredictsByIds(c, u.PredictIds); err != nil {
		log.Infof(c, "%s predictions not found, %v", desc, err)
		return []MatchJson{}
	}

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(t); tb == nil {
		log.Infof(c, "%s tournament builder not found", desc)
		return []MatchJson{}
	}

	mapIdTeams := tb.MapOfIdTeams(c, t)
	mapTeamCodes := tb.MapOfTeamCodes()

	matchesJson := make([]MatchJson, len(matches))
	for i, m := range matches {
		matchesJson[i].Id = m.Id
		matchesJson[i].IdNumber = m.IdNumber
		matchesJson[i].Date = m.Date
		matchesJson[i].Team1 = mapIdTeams[m.TeamId1]
		matchesJson[i].Team2 = mapIdTeams[m.TeamId2]
		matchesJson[i].Iso1 = mapTeamCodes[matchesJson[i].Team1]
		matchesJson[i].Iso2 = mapTeamCodes[matchesJson[i].Team2]

		matchesJson[i].Location = m.Location
		matchesJson[i].Result1 = m.Result1
		matchesJson[i].Result2 = m.Result2
		matchesJson[i].Finished = m.Finished
		matchesJson[i].Ready = m.Ready
		matchesJson[i].CanPredict = m.CanPredict
		if hasMatch, j := predicts.ContainsMatchId(m.Id); hasMatch == true {
			matchesJson[i].HasPredict = true
			matchesJson[i].Predict = fmt.Sprintf("%v - %v", predicts[j].Result1, predicts[j].Result2)
		} else {
			matchesJson[i].HasPredict = false
		}
	}

	return matchesJson
}

// From a tournament entity return an array of second phase MatchJson data structure.
// second phase matches will have the specific rules in there team names
func buildSecondPhaseMatches(c appengine.Context, t *mdl.Tournament, u *mdl.User) []MatchJson {

	matches2ndPhase := mdl.Matches(c, t.Matches2ndStage)

	var predicts mdl.Predicts
	var err error
	if predicts, err = mdl.PredictsByIds(c, u.PredictIds); err != nil {
		return []MatchJson{}
	}

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(t); tb == nil {
		return []MatchJson{}
	}

	mapIdTeams := tb.MapOfIdTeams(c, t)
	mapTeamCodes := tb.MapOfTeamCodes()

	matchesJson := make([]MatchJson, len(matches2ndPhase))

	// append 2nd round to first one
	for i, m := range matches2ndPhase {
		matchesJson[i].Id = m.Id
		matchesJson[i].IdNumber = m.IdNumber
		matchesJson[i].Date = m.Date
		rule := strings.Split(m.Rule, " ")
		if len(rule) == 2 {
			matchesJson[i].Team1 = rule[0]
			matchesJson[i].Team2 = rule[1]
			if _, ok := mapTeamCodes[rule[0]]; ok {
				matchesJson[i].Iso1 = mapTeamCodes[rule[0]]
			}
			if _, ok := mapTeamCodes[rule[1]]; ok {
				matchesJson[i].Iso2 = mapTeamCodes[rule[1]]
			}
		} else {
			if m.TeamId1 > 0 {
				matchesJson[i].Team1 = mapIdTeams[m.TeamId1]
			} else {
				matchesJson[i].Team1 = rule[0]
			}

			if m.TeamId2 > 0 {
				matchesJson[i].Team2 = mapIdTeams[m.TeamId2]
			} else {
				matchesJson[i].Team2 = rule[len(rule)-1]
			}

			matchesJson[i].Iso1 = mapTeamCodes[mapIdTeams[m.TeamId1]]
			matchesJson[i].Iso2 = mapTeamCodes[mapIdTeams[m.TeamId2]]

		}

		matchesJson[i].Location = m.Location
		matchesJson[i].Result1 = m.Result1
		matchesJson[i].Result2 = m.Result2
		matchesJson[i].Finished = m.Finished
		matchesJson[i].Ready = m.Ready
		matchesJson[i].CanPredict = m.CanPredict

		if hasMatch, j := predicts.ContainsMatchId(m.Id); hasMatch == true {
			matchesJson[i].HasPredict = true
			matchesJson[i].Predict = fmt.Sprintf("%v - %v", predicts[j].Result1, predicts[j].Result2)
		} else {
			matchesJson[i].HasPredict = false
		}
	}
	return matchesJson
}
