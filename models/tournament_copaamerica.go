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

package models

import (
	"fmt"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers/log"
)

type CopaAmericaTournament struct {
}

// Map of groups, key: group name, value: string array of teams.
func (cat CopaAmericaTournament) MapOfGroups() map[string][]string {
	var mapCAGroups map[string][]string
	mapCAGroups = make(map[string][]string)
	mapCAGroups["A"] = []string{"Chile", "Mexico", "Ecuador", "Bolivia"}
	mapCAGroups["B"] = []string{"Argentina", "Uruguay", "Paraguay", "Jamaica"}
	mapCAGroups["C"] = []string{"Brazil", "Colombia", "Peru", "Venezuela"}

	return mapCAGroups
}

// Map of country codes, key: team name, value: ISO code
// example: Brazil: BR
func (cat CopaAmericaTournament) MapOfTeamCodes() map[string]string {

	var codes map[string]string
	codes = make(map[string]string)

	codes["Chile"] = "cl"
	codes["Mexico"] = "mx"
	codes["Ecuador"] = "ec"
	codes["Bolivia"] = "bo"
	codes["Argentina"] = "ar"
	codes["Uruguay"] = "uy"
	codes["Paraguay"] = "py"
	codes["Jamaica"] = "jm"
	codes["Brazil"] = "br"
	codes["Colombia"] = "co"
	codes["Peru"] = "pe"
	codes["Venezuela"] = "ve"

	return codes
}

// MapOfGroupMatches return a map with the maches of all the groups
// key: group name, value: array of array of strings with match information:
// MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation
//
// Example:
//
// 	Group A:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (cat CopaAmericaTournament) MapOfGroupMatches() map[string][][]string {

	mapGroupMatches := make(map[string][][]string)

	const (
		cMatchId       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	mA1 := []string{"1", "Jun/11/2015", "Chile", "Ecuador", "Estadio Nacional, Santiago"}
	mA2 := []string{"2", "Jun/12/2015", "Mexico", "Bolivia", "Estadio Sausalito, Viña del Mar"}
	mA3 := []string{"7", "Jun/15/2015", "Ecuador", "Bolivia", "Estadio Elías Figueroa, Valparaíso"}
	mA4 := []string{"8", "Jun/15/2015", "Chile", "Mexico", "Estadio Nacional, Santiago"}
	mA5 := []string{"13", "Jun/19/2015", "Mexico", "Ecuador", "Estadio El Teniente, Rancagua"}
	mA6 := []string{"14", "Jun/19/2015", "Chile", "Bolivia", "Estadio Nacional, Santiago"}

	mB1 := []string{"3", "Jun/13/2015", "Uruguay", "Jamaica", "Estadio Regional de Antofagasta, Antofagasta"}
	mB2 := []string{"4", "Jun/13/2015", "Argentina", "Paraguay", "Estadio La Portada, La Serena"}
	mB3 := []string{"9", "Jun/16/2015", "Paraguay", "Jamaica", "Estadio Regional de Antofagasta, Antofagasta"}
	mB4 := []string{"10", "Jun/16/2015", "Argentina", "Uruguay", "Estadio La Portada, La Serena"}
	mB5 := []string{"15", "Jun/20/2015", "Uruguay", "Paraguay", "Estadio La Portada, La Serena"}
	mB6 := []string{"16", "Jun/20/2015", "Argentina", "Jamaica", "Estadio Sausalito, Viña del Mar"}

	mC1 := []string{"5", "Jun/14/2015", "Colombia", "Venezuela", "Estadio El Teniente, Rancagua"}
	mC2 := []string{"6", "Jun/14/2015", "Brazil", "Peru", "Estadio Municipal Germán Becker, Temuco"}
	mC3 := []string{"11", "Jun/17/2015", "Brazil", "Colombia", "Estadio Monumental David Arellano, Santiago"}
	mC4 := []string{"12", "Jun/18/2015", "Peru", "Venezuela", "Estadio Elías Figueroa, Valparaíso"}
	mC5 := []string{"17", "Jun/21/2015", "Colombia", "Peru", "Estadio Municipal Germán Becker, Temuco"}
	mC6 := []string{"18", "Jun/21/2015", "Brazil", "Venezuela", "Estadio Monumental David Arellano, Santiago"}

	var matchesA [][]string
	var matchesB [][]string
	var matchesC [][]string

	matchesA = append(matchesA, mA1, mA2, mA3, mA4, mA5, mA6)
	matchesB = append(matchesB, mB1, mB2, mB3, mB4, mB5, mB6)
	matchesC = append(matchesC, mC1, mC2, mC3, mC4, mC5, mC6)

	mapGroupMatches["A"] = matchesA
	mapGroupMatches["B"] = matchesB
	mapGroupMatches["C"] = matchesC

	return mapGroupMatches
}

// MapOf2ndRoundMatches returns the map of 2nd round matches, of the tournament.
// key: round number, value: array of array of strings with match information:
// MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation
//
// Example:
//
//	round 16:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (cat CopaAmericaTournament) MapOf2ndRoundMatches() map[string][][]string {

	// 17 Quarter-finals
	m2nd1 := []string{"19", "Jun/24/2015", "1A", "3BC", "Estadio Nacional, Santiago"}
	m2nd2 := []string{"20", "Jun/25/2015", "2A", "2C", "Estadio Municipal Germán Becker, Temuco"}
	m2nd3 := []string{"21", "Jun/26/2015", "1B", "3AC", "Estadio Sausalito, Viña del Mar"}
	m2nd4 := []string{"22", "Jun/27/2015", "1C", "2B", "Estadio Municipal de Concepción, Concepción"}
	// 18 Semi-finals
	m2nd5 := []string{"23", "Jun/29/2015", "W19", "W20", "Estadio Nacional, Santiago"}
	m2nd6 := []string{"24", "Jun/30/2015", "W21", "W22", "Estadio Municipal de Concepción, Concepción"}
	//19 Round 19  -  Match for third place
	m2nd7 := []string{"25", "Jul/03/2015", "L23", "L24", "Estadio Municipal de Concepción, Concepción"}
	//"20" Final
	m2nd8 := []string{"26", "Jul/04/2015", "W23", "W24", "Estadio Nacional, Santiago"}

	var round17 [][]string
	var round18 [][]string
	var round19 [][]string
	var round20 [][]string

	round17 = append(round17, m2nd1, m2nd2, m2nd3, m2nd4)
	round18 = append(round18, m2nd5, m2nd6)
	round19 = append(round19, m2nd7)
	round20 = append(round20, m2nd8)

	mapMatches2ndRound := make(map[string][][]string)
	mapMatches2ndRound["17"] = round17
	mapMatches2ndRound["18"] = round18
	mapMatches2ndRound["19"] = round19
	mapMatches2ndRound["20"] = round20

	return mapMatches2ndRound
}

// ArrayOfPhases returns an array of the phases names of world cup tournament:
// FirstStage, RoundOf16, QuarterFinals, SemiFinals, ThirdPlace, Finals
//
func (cat CopaAmericaTournament) ArrayOfPhases() []string {
	return []string{cFirstStage, cQuarterFinals, cSemiFinals, cThirdPlace, cFinals}
}

// MapOfPhaseIntervals returns a map with key the corresponding phase in the copa america tournament
// at value a tuple that represent the match number interval in which the phase take place:
// first stage: matches 1 to 48
// Round of 16: matches 49 to 56
// Quarte-finals: matches 57 to 60
// Semi-finals: matches 61 to 62
// Thrid Place: match 63
// Finals: match 64
func (cat CopaAmericaTournament) MapOfPhaseIntervals() map[string][]int64 {

	limits := make(map[string][]int64)
	limits[cFirstStage] = []int64{1, 18}
	limits[cQuarterFinals] = []int64{19, 22}
	limits[cSemiFinals] = []int64{23, 24}
	limits[cThirdPlace] = []int64{25, 25}
	limits[cFinals] = []int64{26, 26}
	return limits
}

// Create Copa America tournament entity 2015.
func CreateCopaAmerica(c appengine.Context, adminId int64) (*Tournament, error) {
	desc := "Copa America"
	// create new tournament
	log.Infof(c, "%s: start", desc)

	cat := CopaAmericaTournament{}

	mapCATGroups := cat.MapOfGroups()
	mapCountryCodes := cat.MapOfTeamCodes()
	mapTeamId := make(map[string]int64)

	// mapGroupMatches is a map where the key is a string which represent the group
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// map[string][][]string
	// example: {"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}
	mapGroupMatches := cat.MapOfGroupMatches()

	const (
		cMatchId       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	// for date parsing
	const shortForm = "Jan/02/2006"

	// matches1stStageIds is an array of  int64
	// where we allocate IDs of the Tmatches entities
	// we will store them in the tournament entity for easy retreival later on.
	matchesA := mapGroupMatches["A"]
	matches1stStageIds := make([]int64, len(mapGroupMatches)*len(matchesA))

	log.Infof(c, "%s: maps ready", desc)

	// build groups, and teams
	groups := make([]Tgroup, len(mapCATGroups))
	for groupName, teams := range mapCATGroups {
		log.Infof(c, "---------------------------------------")
		log.Infof(c, "%s: working with group: %v", desc, groupName)
		log.Infof(c, "%s: teams: %v", desc, teams)

		var group Tgroup
		group.Name = groupName
		groupIndex := int64(groupName[0]) - 65
		group.Teams = make([]Tteam, len(teams))
		group.Points = make([]int64, len(teams))
		group.GoalsF = make([]int64, len(teams))
		group.GoalsA = make([]int64, len(teams))

		for i, teamName := range teams {
			log.Infof(c, "%s: team: %v", desc, teamName)

			teamID, _, err1 := datastore.AllocateIDs(c, "Tteam", nil, 1)
			if err1 != nil {
				return nil, err1
			}
			log.Infof(c, "%s: team: %v allocateIDs ok", desc, teamName)

			teamkey := datastore.NewKey(c, "Tteam", "", teamID, nil)
			log.Infof(c, "%s: team: %v NewKey ok", desc, teamName)

			team := &Tteam{teamID, teamName, mapCountryCodes[teamName]}
			log.Infof(c, "%s: team: %v instance of team ok", desc, teamName)

			_, err := datastore.Put(c, teamkey, team)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "%s: team: %v put in datastore ok", desc, teamName)
			group.Teams[i] = *team
			mapTeamId[teamName] = teamID
		}

		// build group matches:
		log.Infof(c, "%s: building group matches", desc)

		groupMatches := mapGroupMatches[groupName]
		group.Matches = make([]Tmatch, len(groupMatches))

		for matchIndex, matchData := range groupMatches {
			log.Infof(c, "%s: match data: %v", desc, matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "%s: match: %v allocateIDs ok", desc, matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "%s: match: new key ok", desc)

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalId, _ := strconv.Atoi(matchData[cMatchId])
			emptyrule := ""
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalId),
				matchTime,
				mapTeamId[matchData[cMatchTeam1]],
				mapTeamId[matchData[cMatchTeam2]],
				matchData[cMatchLocation],
				emptyrule,
				emptyresult,
				emptyresult,
				false,
				true,
				true,
			}

			log.Infof(c, "%s: match: build match ok", desc)

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "%s: match: %v put in datastore ok", desc, matchData)
			group.Matches[matchIndex] = *match

			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches1stStageIds[int64(matchInternalId)-1] = matchID
		}

		groupID, _, err1 := datastore.AllocateIDs(c, "Tgroup", nil, 1)
		if err1 != nil {
			return nil, err1
		}

		log.Infof(c, "%s: Group: %v allocate ID ok", desc, groupName)

		groupkey := datastore.NewKey(c, "Tgroup", "", groupID, nil)
		log.Infof(c, "%s: Group: %v New Key ok", desc, groupName)

		group.Id = groupID
		groups[groupIndex] = group
		_, err := datastore.Put(c, groupkey, &group)
		if err != nil {
			return nil, err
		}
		log.Infof(c, "%s: Group: %v put in datastore ok", desc, groupName)
	}

	// build array of group ids
	groupIds := make([]int64, len(groups))
	for i := range groupIds {
		groupIds[i] = groups[i].Id
	}

	log.Infof(c, "%s: build of groups ids complete: %v", desc, groupIds)

	// matches 2nd stage
	matches2ndStageIds := make([]int64, 0)
	userIds := make([]int64, 0)
	teamIds := make([]int64, 0)

	// mapMatches2ndRound  is a map where the key is a string which represent the rounds
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// mapMatches2ndRound is a map[string][][]string
	// example: {"64", "Jul/13/2014", "W61", "W62", "Rio de Janeiro"}
	mapMatches2ndRound := cat.MapOf2ndRoundMatches()

	// build matches 2nd phase
	for roundNumber, roundMatches := range mapMatches2ndRound {
		log.Infof(c, "%s: building 2nd round matches: round number %v", desc, roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "%s: second phase match data: %v", desc, matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "%s: match: %v allocateIDs ok", desc, matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "%s: match: new key ok", desc)

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalId, _ := strconv.Atoi(matchData[cMatchId])

			rule := fmt.Sprintf("%s %s", matchData[cMatchTeam1], matchData[cMatchTeam2])
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalId),
				matchTime,
				0, // second round matches start with ids at 0
				0, // second round matches start with ids at 0
				matchData[cMatchLocation],
				rule,
				emptyresult,
				emptyresult,
				false,
				false,
				true,
			}
			log.Infof(c, "%s: match 2nd round: build match ok", desc)

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "%s: 2nd round match: %v put in datastore ok", desc, matchData)
			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches2ndStageIds = append(matches2ndStageIds, matchID)
		}
	}

	tstart, _ := time.Parse(shortForm, "Jun/11/2015")
	tend, _ := time.Parse(shortForm, "Jul/04/2015")
	adminIds := make([]int64, 1)
	adminIds[0] = adminId
	name := "2015 Copa America"
	description := "Chile"
	var tournament *Tournament
	var err error
	if tournament, err = CreateTournament(c, name, description, tstart, tend, adminId); err != nil {
		log.Infof(c, "%s: something went wrong when creating tournament.", desc)
	} else {
		tournament.GroupIds = groupIds
		tournament.Matches1stStage = matches1stStageIds
		tournament.Matches2ndStage = matches2ndStageIds
		tournament.UserIds = userIds
		tournament.TeamIds = teamIds
		tournament.IsFirstStageComplete = false
		if err1 := tournament.Update(c); err1 != nil {
			log.Infof(c, "%s: unable to udpate tournament.", desc)
		}
	}

	log.Infof(c, "%s: instance of tournament ready", desc)

	return tournament, nil
}

// MapOfIdTeams returns a map of team IDs as keys and team names as values.
//
func (cat CopaAmericaTournament) MapOfIdTeams(c appengine.Context, tournament *Tournament) map[int64]string {

	mapIdTeams := make(map[int64]string)

	groups := Groups(c, tournament.GroupIds)
	for _, g := range groups {
		for _, t := range g.Teams {
			mapIdTeams[t.Id] = t.Name
		}
	}
	return mapIdTeams
}
