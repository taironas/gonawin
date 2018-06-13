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

// const (
// 	cFirstStage    = "First Stage"
// 	cRoundOf16     = "Round of 16"
// 	cQuarterFinals = "Quarter-finals"
// 	cSemiFinals    = "Semi-finals"
// 	cThirdPlace    = "Third Place"
// 	cFinals        = "Finals"
// )

// WorldCupTournament2018 represents a World Cup tournament.
//
type WorldCupTournament2018 struct {
}

// MapOfGroups is a map containing the groups of a tournament.
// key: group name, value: string array of teams.
//
func (wct WorldCupTournament2018) MapOfGroups() map[string][]string {
	var mapWCGroups map[string][]string
	mapWCGroups = make(map[string][]string)
	mapWCGroups["A"] = []string{"Russia", "Saudi Arabia", "Egypt", "Uruguay"}
	mapWCGroups["B"] = []string{"Portugal", "Spain", "Morocco", "Iran"}
	mapWCGroups["C"] = []string{"France", "Australia", "Peru", "Denmark"}
	mapWCGroups["D"] = []string{"Argentina", "Iceland", "Croatia", "Nigeria"}
	mapWCGroups["E"] = []string{"Brazil", "Switzerland", "Costa Rica", "Serbia"}
	mapWCGroups["F"] = []string{"Germany", "Mexico", "Sweden", "South Korea"}
	mapWCGroups["G"] = []string{"Belgium", "Panama", "Tunisia", "England"}
	mapWCGroups["H"] = []string{"Poland", "Senegal", "Colombia", "Japan"}

	return mapWCGroups
}

// MapOfTeamCodes is map containing the team codes.
// key: team name, value: code
// example: Paris Saint-Germain: PSG
//
// example: Brazil: BR
func (wct WorldCupTournament2018) MapOfTeamCodes() map[string]string {

	var codes map[string]string
	codes = make(map[string]string)

	codes["Russia"] = "ru"
	codes["Saudi Arabia"] = "sa"
	codes["Egypt"] = "eg"
	codes["Uruguay"] = "uy"
	codes["Portugal"] = "pt"
	codes["Spain"] = "es"
	codes["Morocco"] = "ma"
	codes["Iran"] = "ir"
	codes["France"] = "fr"
	codes["Australia"] = "au"
	codes["Peru"] = "pe"
	codes["Denmark"] = "dk"
	codes["Argentina"] = "ar"
	codes["Iceland"] = "is"
	codes["Croatia"] = "hr"
	codes["Nigeria"] = "ng"
	codes["Brazil"] = "br"
	codes["Switzerland"] = "ch"
	codes["Costa Rica"] = "cr"
	codes["Serbia"] = "rs"
	codes["Germany"] = "de"
	codes["Mexico"] = "mx"
	codes["Sweden"] = "se"
	codes["South Korea"] = "kr"
	codes["Belgium"] = "be"
	codes["Panama"] = "pa"
	codes["Tunisia"] = "tn"
	codes["England"] = "gb-eng"
	codes["Poland"] = "pl"
	codes["Senegal"] = "sn"
	codes["Colombia"] = "co"
	codes["Japan"] = "jp"

	return codes
}

// MapOfGroupMatches is a map containing the matches accessible by group.
//
// Example:
//
// 	Group A:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (wct WorldCupTournament2018) MapOfGroupMatches() map[string][][]string {

	mapGroupMatches := make(map[string][][]string)

	const (
		cMatchID       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	mA1 := []string{"1", "Jun/14/2018", "Russia", "Saudi Arabia", "Moscow"}
	mA2 := []string{"2", "Jun/15/2018", "Egypt", "Uruguay", "Yekaterinburg"}
	mA3 := []string{"17", "Jun/19/2018", "Russia", "Egypt", "Saint Petersburg"}
	mA4 := []string{"18", "Jun/20/2018", "Uruguay", "Saudi Arabia", "Rostov-on-Don"}
	mA5 := []string{"33", "Jun/25/2018", "Uruguay", "Russia", "Samara"}
	mA6 := []string{"34", "Jun/25/2018", "Saudi Arabia", "Egypt", "Volgograd"}

	mB1 := []string{"4", "Jun/15/2018", "Morocco", "Iran", "Saint Petersburg"}
	mB2 := []string{"3", "Jun/15/2018", "Portugal", "Spain", "Sochi"}
	mB3 := []string{"19", "Jun/20/2018", "Portugal", "Morocco", "Moscow"}
	mB4 := []string{"20", "Jun/20/2018", "Iran", "Spain", "Kazan"}
	mB5 := []string{"35", "Jun/25/2018", "Iran", "Portugal", "Saransk"}
	mB6 := []string{"36", "Jun/25/2018", "Spain", "Morocco", "Kaliningrad"}

	mC1 := []string{"5", "Jun/16/2018", "France", "Australia", "Kazan"}
	mC2 := []string{"6", "Jun/16/2018", "Peru", "Denmark", "Saransk"}
	mC3 := []string{"22", "Jun/21/2018", "Denmark", "Australia", "Samara"}
	mC4 := []string{"21", "Jun/21/2018", "France", "Peru", "Yekaterinburg"}
	mC5 := []string{"37", "Jun/26/2018", "Denmark", "France", "Moscow"}
	mC6 := []string{"38", "Jun/26/2018", "Australia", "Peru", "Sochi"}

	mD1 := []string{"7", "Jun/16/2018", "Argentina", "Iceland", "Moscow"}
	mD2 := []string{"8", "Jun/16/2018", "Croatia", "Nigeria", "Kaliningrad"}
	mD3 := []string{"23", "Jun/21/2018", "Argentina", "Croatia", "Nizhny Novgorod"}
	mD4 := []string{"24", "Jun/22/2018", "Nigeria", "Iceland", "Volgograd"}
	mD5 := []string{"39", "Jun/26/2018", "Nigeria", "Argentina", "Saint Petersburg"}
	mD6 := []string{"40", "Jun/26/2018", "Iceland", "Croatia", "Rostov-on-Don"}

	mE1 := []string{"9", "Jun/17/2018", "Brazil", "Switzerland", "Rostov-on-Don"}
	mE2 := []string{"10", "Jun/17/2018", "Costa Rica", "Serbia", "Samara"}
	mE3 := []string{"25", "Jun/22/2018", "Brazil", "Costa Rica", "Saint Petersburg"}
	mE4 := []string{"26", "Jun/22/2018", "Serbia", "Switzerland", "Kaliningrad"}
	mE5 := []string{"41", "Jun/27/2018", "Serbia", "Brazil", "Moscow"}
	mE6 := []string{"42", "Jun/27/2018", "Switzerland", "Costa Rica", "Nizhny Novgorod"}

	mF1 := []string{"11", "Jun/17/2018", "Germany", "Mexico", "Moscow"}
	mF2 := []string{"12", "Jun/18/2018", "Sweden", "South Korea", "Nizhny Novgorod"}
	mF3 := []string{"28", "Jun/23/2018", "South Korea", "Mexico", "Rostov-on-Don"}
	mF4 := []string{"27", "Jun/23/2018", "Germany", "Sweden", "Sochi"}
	mF5 := []string{"43", "Jun/27/2018", "South Korea", "Germany", "Kazan"}
	mF6 := []string{"44", "Jun/27/2018", "Mexico", "Sweden", "Yekaterinburg"}

	mG1 := []string{"13", "Jun/18/2018", "Belgium", "Panama", "Sochi"}
	mG2 := []string{"14", "Jun/18/2018", "Tunisia", "England", "Volgograd"}
	mG3 := []string{"29", "Jun/23/2018", "Belgium", "Tunisia", "Moscow"}
	mG4 := []string{"30", "Jun/24/2018", "England", "Panama", "Nizhny Novgorod"}
	mG5 := []string{"45", "Jun/28/2018", "England", "Belgium", "Kaliningrad"}
	mG6 := []string{"46", "Jun/28/2018", "Panama", "Tunisia", "Saransk"}

	mH1 := []string{"16", "Jun/19/2018", "Colombia", "Japan", "Saransk"}
	mH2 := []string{"15", "Jun/19/2018", "Poland", "Senegal", "Moscow"}
	mH3 := []string{"32", "Jun/24/2018", "Japan", "Senegal", "Yekaterinburg"}
	mH4 := []string{"31", "Jun/24/2018", "Poland", "Colombia", "Kazan"}
	mH5 := []string{"47", "Jun/28/2018", "Japan", "Poland", "Volgograd"}
	mH6 := []string{"48", "Jun/28/2018", "Senegal", "Colombia", "Samara"}

	var matchesA [][]string
	var matchesB [][]string
	var matchesC [][]string
	var matchesD [][]string
	var matchesE [][]string
	var matchesF [][]string
	var matchesG [][]string
	var matchesH [][]string

	matchesA = append(matchesA, mA1, mA2, mA3, mA4, mA5, mA6)
	matchesB = append(matchesB, mB1, mB2, mB3, mB4, mB5, mB6)
	matchesC = append(matchesC, mC1, mC2, mC3, mC4, mC5, mC6)
	matchesD = append(matchesD, mD1, mD2, mD3, mD4, mD5, mD6)
	matchesE = append(matchesE, mE1, mE2, mE3, mE4, mE5, mE6)
	matchesF = append(matchesF, mF1, mF2, mF3, mF4, mF5, mF6)
	matchesG = append(matchesG, mG1, mG2, mG3, mG4, mG5, mG6)
	matchesH = append(matchesH, mH1, mH2, mH3, mH4, mH5, mH6)

	mapGroupMatches["A"] = matchesA
	mapGroupMatches["B"] = matchesB
	mapGroupMatches["C"] = matchesC
	mapGroupMatches["D"] = matchesD
	mapGroupMatches["E"] = matchesE
	mapGroupMatches["F"] = matchesF
	mapGroupMatches["G"] = matchesG
	mapGroupMatches["H"] = matchesH

	return mapGroupMatches
}

// MapOf2ndRoundMatches returns the Map of 2nd round matches, of the world cup tournament.
// key: round number, value: array of array of strings with match information ( MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation)
//
// Example:
//
//	round 16:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (wct WorldCupTournament2018) MapOf2ndRoundMatches() map[string][][]string {

	// Round of 16
	m2nd1 := []string{"49", "Jun/30/2018", "1A", "2B", "Sochi"}
	m2nd2 := []string{"50", "Jun/30/2018", "1C", "2D", "Kazan"}
	m2nd3 := []string{"51", "Jul/01/2018", "1B", "2A", "Moscow"}
	m2nd4 := []string{"52", "Jul/01/2018", "1D", "2C", "Nizhny Novgorod"}
	m2nd5 := []string{"53", "Jul/02/2018", "1E", "2F", "Samara"}
	m2nd6 := []string{"54", "Jul/02/2018", "1G", "2H", "Rostov-on-Don"}
	m2nd7 := []string{"55", "Jul/03/2018", "1F", "2E", "Saint Petersburg"}
	m2nd8 := []string{"56", "Jul/03/2018", "1H", "2G", "Moscow"}
	// 17 Quarter-finals
	m2nd9 := []string{"57", "Jul/06/2018", "W49", "W50", "Nizhny Novgorod"}
	m2nd10 := []string{"58", "Jul/06/2018", "W53", "W54", "Kazan"}
	m2nd11 := []string{"59", "Jul/07/2018", "W51", "W52", "Samara"}
	m2nd12 := []string{"60", "Jul/07/2018", "W55", "W56", "Sochi"}
	// 18 Semi-finals
	m2nd13 := []string{"61", "Jul/10/2018", "W57", "W58", "Saint Petersburg"}
	m2nd14 := []string{"62", "Jul/11/2018", "W59", "W60", "Moscow"}
	//19 Round 19  -  Match for third place
	m2nd15 := []string{"63", "Jul/14/2018", "L61", "L62", "Saint Petersburg"}
	//"20" Final
	m2nd16 := []string{"64", "Jul/15/2018", "W61", "W62", "Moscow"}

	var round16 [][]string
	var round17 [][]string
	var round18 [][]string
	var round19 [][]string
	var round20 [][]string

	round16 = append(round16, m2nd1, m2nd2, m2nd3, m2nd4, m2nd5, m2nd6, m2nd7, m2nd8)
	round17 = append(round17, m2nd9, m2nd10, m2nd11, m2nd12)
	round18 = append(round18, m2nd13, m2nd14)
	round19 = append(round19, m2nd15)
	round20 = append(round20, m2nd16)

	mapMatches2ndRound := make(map[string][][]string)
	mapMatches2ndRound["16"] = round16
	mapMatches2ndRound["17"] = round17
	mapMatches2ndRound["18"] = round18
	mapMatches2ndRound["19"] = round19
	mapMatches2ndRound["20"] = round20

	return mapMatches2ndRound
}

// ArrayOfPhases returns an array of the phases names of champions league tournament: (QuarterFinals, SemiFinals, Finals).
//
func (wct WorldCupTournament2018) ArrayOfPhases() []string {
	return []string{cFirstStage, cRoundOf16, cQuarterFinals, cSemiFinals, cThirdPlace, cFinals}
}

// MapOfPhaseIntervals builds a map with key the corresponding phase in the world cup tournament
// at value a tuple that represent the match number interval in which the phase take place:
// first stage: matches 1 to 48
// Round of 16: matches 49 to 56
// Quarte-finals: matches 57 to 60
// Semi-finals: matches 61 to 62
// Third Place: match 63
// Finals: match 64
//
func (wct WorldCupTournament2018) MapOfPhaseIntervals() map[string][]int64 {

	limits := make(map[string][]int64)
	limits[cFirstStage] = []int64{1, 48}
	limits[cRoundOf16] = []int64{49, 56}
	limits[cQuarterFinals] = []int64{57, 60}
	limits[cSemiFinals] = []int64{61, 62}
	limits[cThirdPlace] = []int64{63, 63}
	limits[cFinals] = []int64{64, 64}
	return limits
}

// CreateWorldCup2018 creates world cup tournament entity 2018.
//
func CreateWorldCup2018(c appengine.Context, adminID int64) (*Tournament, error) {
	// create new tournament
	log.Infof(c, "World Cup: start")

	wct := WorldCupTournament2018{}

	mapWCGroups := wct.MapOfGroups()
	mapCountryCodes := wct.MapOfTeamCodes()
	mapTeamID := make(map[string]int64)

	// mapGroupMatches is a map where the key is a string which represent the group
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// map[string][][]string
	// example: {"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}
	mapGroupMatches := wct.MapOfGroupMatches()

	const (
		cMatchID       = 0
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
	matches1stStageIds := make([]int64, 8*len(matchesA))

	log.Infof(c, "World Cup: maps ready")

	// build groups, and teams
	groups := make([]Tgroup, len(mapWCGroups))
	for groupName, teams := range mapWCGroups {
		log.Infof(c, "---------------------------------------")
		log.Infof(c, "World Cup: working with group: %v", groupName)
		log.Infof(c, "World Cup: teams: %v", teams)

		var group Tgroup
		group.Name = groupName
		groupIndex := int64(groupName[0]) - 65
		group.Teams = make([]Tteam, len(teams))
		group.Points = make([]int64, len(teams))
		group.GoalsF = make([]int64, len(teams))
		group.GoalsA = make([]int64, len(teams))
		for i, teamName := range teams {
			log.Infof(c, "World Cup: team: %v", teamName)

			teamID, _, err1 := datastore.AllocateIDs(c, "Tteam", nil, 1)
			if err1 != nil {
				return nil, err1
			}
			log.Infof(c, "World Cup: team: %v allocateIDs ok", teamName)

			teamkey := datastore.NewKey(c, "Tteam", "", teamID, nil)
			log.Infof(c, "World Cup: team: %v NewKey ok", teamName)

			team := &Tteam{teamID, teamName, mapCountryCodes[teamName]}
			log.Infof(c, "World Cup: team: %v instance of team ok", teamName)

			_, err := datastore.Put(c, teamkey, team)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "World Cup: team: %v put in datastore ok", teamName)
			group.Teams[i] = *team
			mapTeamID[teamName] = teamID
		}

		// build group matches:
		log.Infof(c, "World Cup: building group matches")

		groupMatches := mapGroupMatches[groupName]
		group.Matches = make([]Tmatch, len(groupMatches))

		for matchIndex, matchData := range groupMatches {
			log.Infof(c, "World Cup: match data: %v", matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "World Cup: match: %v allocateIDs ok", matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "World Cup: match: new key ok")

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalID, _ := strconv.Atoi(matchData[cMatchID])
			emptyrule := ""
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalID),
				matchTime,
				mapTeamID[matchData[cMatchTeam1]],
				mapTeamID[matchData[cMatchTeam2]],
				matchData[cMatchLocation],
				emptyrule,
				emptyresult,
				emptyresult,
				false,
				true,
				true,
			}
			log.Infof(c, "World Cup: match: build match ok")

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "World Cup: match: %v put in datastore ok", matchData)
			group.Matches[matchIndex] = *match

			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches1stStageIds[int64(matchInternalID)-1] = matchID
		}

		groupID, _, err1 := datastore.AllocateIDs(c, "Tgroup", nil, 1)
		if err1 != nil {
			return nil, err1
		}

		log.Infof(c, "World Cup: Group: %v allocate Id ok", groupName)

		groupkey := datastore.NewKey(c, "Tgroup", "", groupID, nil)
		log.Infof(c, "World Cup: Group: %v New Key ok", groupName)

		group.Id = groupID
		groups[groupIndex] = group
		_, err := datastore.Put(c, groupkey, &group)
		if err != nil {
			return nil, err
		}
		log.Infof(c, "World Cup: Group: %v put in datastore ok", groupName)
	}

	// build array of group ids
	groupIds := make([]int64, 8)
	for i := range groupIds {
		groupIds[i] = groups[i].Id
	}

	log.Infof(c, "World Cup: build of groups ids complete: %v", groupIds)

	// matches 2nd stage
	var matches2ndStageIds []int64
	var userIds []int64
	var teamIds []int64
	// mapMatches2ndRound  is a map where the key is a string which represent the rounds
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// mapMatches2ndRound is a map[string][][]string
	// example: {"64", "Jul/13/2014", "W61", "W62", "Rio de Janeiro"}
	mapMatches2ndRound := wct.MapOf2ndRoundMatches()

	// build matches 2nd phase
	for roundNumber, roundMatches := range mapMatches2ndRound {
		log.Infof(c, "World Cup: building 2nd round matches: round number %v", roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "World Cup: second phase match data: %v", matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "World Cup: match: %v allocateIDs ok", matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "World Cup: match: new key ok")

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalID, _ := strconv.Atoi(matchData[cMatchID])

			rule := fmt.Sprintf("%s %s", matchData[cMatchTeam1], matchData[cMatchTeam2])
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalID),
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
			log.Infof(c, "World Cup: match 2nd round: build match ok")

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "World Cup: 2nd round match: %v put in datastore ok", matchData)
			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches2ndStageIds = append(matches2ndStageIds, matchID)
		}
	}

	tstart, _ := time.Parse(shortForm, "Jun/14/2018")
	tend, _ := time.Parse(shortForm, "Jul/15/2018")
	adminIds := make([]int64, 1)
	adminIds[0] = adminID
	name := "2018 FIFA World Cup"
	description := "Russia"
	var tournament *Tournament
	var err error
	if tournament, err = CreateTournament(c, name, description, tstart, tend, adminID); err != nil {
		log.Infof(c, "World Cup: something went wrong when creating tournament.")
	} else {
		tournament.GroupIds = groupIds
		tournament.Matches1stStage = matches1stStageIds
		tournament.Matches2ndStage = matches2ndStageIds
		tournament.UserIds = userIds
		tournament.TeamIds = teamIds
		tournament.IsFirstStageComplete = false
		if err1 := tournament.Update(c); err1 != nil {
			log.Infof(c, "World Cup: unable to udpate tournament.")
		}
	}

	log.Infof(c, "World Cup: instance of tournament ready")

	return tournament, nil
}

// MapOfIDTeams builds a map of teams from tournament entity.
//
func (wct WorldCupTournament2018) MapOfIDTeams(c appengine.Context, tournament *Tournament) map[int64]string {

	mapIDTeams := make(map[int64]string)

	groups := Groups(c, tournament.GroupIds)
	for _, g := range groups {
		for _, t := range g.Teams {
			mapIDTeams[t.Id] = t.Name
		}
	}
	return mapIDTeams
}
